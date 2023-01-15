package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/base/pkg/gadget/metadata"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.uber.org/zap"
)

type Server struct {
	labelStore   labelstore.LabelStore
	inputRef     *gadget.DataRef
	maxBatchSize int
	log          *zap.Logger
}

func NewServer(
	labelStore labelstore.LabelStore,
	initInputRef *gadget.DataRef,
	maxBatchSize int,
	log *zap.Logger,
) *Server {
	return &Server{
		labelStore,
		initInputRef,
		maxBatchSize,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	meta := &metadata.Metadata{
		Name:         "crop",
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

	// setup the input channel proxy methods
	gadget.SetupInputChannel(router, "default", input, s.log)

	// the output data handler that actually applies the crop
	router.HandleFunc("/output/default/x", handleGetCroppedOutput(
		s.labelStore,
		input,
		s.log,
	))

	// standard handler to GET/PUT individual output labels
	router.HandleFunc("/output/default/y", gadget.HandleOutputLabel(
		s.labelStore,
		input,
		validateCrop,
		s.log,
	))

	// output label sampler that crops the data
	router.HandleFunc("/sample/output/default/y", handleSampleCroppedOutput(
		s.maxBatchSize,
		s.log,
	))

	return (&http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}).ListenAndServe()
}
