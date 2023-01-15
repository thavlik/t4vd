package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/base/pkg/gadget/metadata"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

// Server is a server that samples frames from an input gadget
// and allows the user to label those frames. The labels are
// stored in a database, and they can be sampled like any other
// gadget output.
type Server struct {
	labelStore   labelstore.LabelStore
	gadgetID     string
	maxBatchSize int
	inputRef     *gadget.DataRef
	log          *zap.Logger
}

// NewServer creates a new filter server.
func NewServer(
	labelStore labelstore.LabelStore,
	gadgetID string,
	maxBatchSize int,
	defaultRef *gadget.DataRef,
	log *zap.Logger,
) *Server {
	return &Server{
		labelStore,
		gadgetID,
		maxBatchSize,
		defaultRef,
		log,
	}
}

// Listen starts the filter server.
func (s *Server) Listen(port int) error {
	meta := &metadata.Metadata{
		GadgetID:     s.gadgetID,
		Name:         "filter",
		MaxBatchSize: s.maxBatchSize,
		Inputs: []*metadata.Channel{{
			Name: "default",
		}},
		Outputs: []*metadata.Channel{{
			Name: "default",
		}},
	}

	router := mux.NewRouter()
	router.HandleFunc("/metadata", gadget.HandleGetMetadata(meta, s.log))

	input := s.inputRef

	// setup the proxy methods for the default input channel
	gadget.SetupInputChannel(router, s.gadgetID, "default", input, s.log)

	// setup methods for the default output channel
	// these are for inserting/retrieving the labels associated with
	// this gadget as well as retrieving the transformed data
	router.HandleFunc("/output/default/x", gadget.HandleGetOutputDataFromRef( // retrieve transformed data by id (identity in this case)
		s.gadgetID,
		input,
		s.log,
	))
	router.HandleFunc("/output/default/y", gadget.HandleOutputLabel( // retrieve a specific label by id (labels stored by this gadget)
		s.labelStore,
		input,
		validateTags,
		s.log,
	))

	router.HandleFunc("/sample/output/default/y", gadget.HandleSampleOutputLabelsFromStore( // sample labels stored by this gadget
		s.labelStore,
		s.maxBatchSize,
		s.log,
	))

	// get/set the state of the input channel
	router.HandleFunc("/state/{channel}", gadget.HandleInputState(
		map[string]*gadget.DataRef{
			"default": input,
		},
		s.log,
	))

	router.HandleFunc("/healthz", base.HealthHandler)
	router.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}
