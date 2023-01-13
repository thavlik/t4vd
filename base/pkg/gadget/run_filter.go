package gadget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

func GetGadget(name string) (*Gadget, error) {
	// TODO: resolve the gadget name into an http endpoint
	// TODO: query the gadget's http server for the gadget info
	// build a Gadget instance from the info
	return nil, errors.New("unimplemented")
}

func Entry(
	port int,
	projectID string,
	labelStore labelstore.LabelStore,
	inputGadgetName string,
	inputChannel string,
	defaultBatchSize int,
	maxBatchSize int,
	log *zap.Logger,
) error {
	inputGadget, err := GetGadget(inputGadgetName)
	if err != nil {
		return errors.Wrap(err, "failed to get input gadget")
	}
	channel, err := inputGadget.GetOutput(
		context.Background(),
		inputChannel,
	)
	if err != nil {
		return errors.Wrap(err, "failed to get input channel")
	}
	return NewFilter(
		projectID,
		labelStore,
		channel,
		defaultBatchSize,
		maxBatchSize,
		log,
	).Listen(port)
}

// Filter is a server that samples frames from an input gadget
// and allows the user to label those frames. The labels are
// stored in a database, and they can be sampled like any other
// gadget output.
type Filter struct {
	projectID        string
	labelStore       labelstore.LabelStore
	channel          OutputChannel
	defaultBatchSize int
	maxBatchSize     int
	log              *zap.Logger
}

// NewFilter creates a new filter server.
func NewFilter(
	projectID string,
	labelStore labelstore.LabelStore,
	channel OutputChannel,
	defaultBatchSize int,
	maxBatchSize int,
	log *zap.Logger,
) *Filter {
	return &Filter{
		projectID,
		labelStore,
		channel,
		defaultBatchSize,
		maxBatchSize,
		log,
	}
}

// Listen starts the filter server.
func (s *Filter) Listen(port int) error {
	mux := mux.NewRouter()
	mux.HandleFunc("/input/default", handleGetFrameInput(
		s.channel,
		s.defaultBatchSize,
		s.maxBatchSize,
		s.log,
	))
	mux.HandleFunc("/output/{channel}", s.handleOutput()) // GET or PUT output data, per channel
	mux.HandleFunc("/label/{id}", s.handleGetLabel())     // retrieve a specific label by id
	mux.HandleFunc("/readyz", base.ReadyHandler)          // readiness probe
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}

// extractBatchSize extracts the batch size from the query string.
// If the batch size is not specified, the default batch size is used.
// If the batch size is too large, an error is returned.
func extractBatchSize(
	query url.Values,
	defaultBatchSize int,
	maxBatchSize int,
) (int, error) {
	if v := query.Get("s"); v != "" {
		sz, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse batch size")
		}
		batchSize := int(sz)
		if batchSize > maxBatchSize {
			return 0, errors.New("batch size too large")
		}
		return batchSize, nil
	}
	return defaultBatchSize, nil
}

// handleGetFrameInput handler for sampling frames from the
// given input channel.
// If batchSize is 0, the default batch size is used.
// This same handler should be used for all gadgets
// that sample frames as input.
// Use this handler to retrieve input frames to be labeled,
// then PUT those frames to the output channel.
func handleGetFrameInput(
	inputChannel OutputChannel,
	defaultBatchSize int,
	maxBatchSize int,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return errors.New("method not allowed")
			}
			batchSize, err := extractBatchSize(
				r.URL.Query(),
				defaultBatchSize,
				maxBatchSize,
			)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return errors.Wrap(err, "failed to extract batch size")
			}
			start := time.Now()
			frames, err := SampleFrames(
				r.Context(),
				inputChannel,
				batchSize,
			)
			if err != nil {
				return errors.Wrap(err, "failed to sample input frames")
			}
			log.Debug("sampled input frames",
				zap.Int("batchSize", batchSize),
				base.Elapsed(start))
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(frames)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(r.RequestURI, zap.Error(err))
		}
	}
}

// handleGetVideoInput handler for sampling video from the
// given channel on the associated gadget. Use this to sample
// videos to be labeled by the host gadget.
func handleGetVideoInput(
	inputChannel OutputChannel,
	defaultBatchSize int,
	maxBatchSize int,
	padding time.Duration,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return errors.New("method not allowed")
			}
			batchSize, err := extractBatchSize(
				r.URL.Query(),
				defaultBatchSize,
				maxBatchSize,
			)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return errors.Wrap(err, "failed to extract batch size")
			}
			start := time.Now()
			videos, err := SampleVideos(
				r.Context(),
				inputChannel,
				batchSize,
				padding,
			)
			if err != nil {
				return errors.Wrap(err, "failed to sample input frames")
			}
			log.Debug("sampled input videos",
				zap.Int("batchSize", batchSize),
				base.Elapsed(start))
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(videos)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error(r.RequestURI, zap.Error(err))
		}
	}
}

// handleGetOutput handler for sampling frames from the
// given output channel. Use this to retrieve labeled frames
// from the host gadget.
func (s *Filter) handleGetOutput(
	w http.ResponseWriter,
	r *http.Request,
) error {
	batchSize, err := extractBatchSize(
		r.URL.Query(),
		s.defaultBatchSize,
		s.maxBatchSize,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.Wrap(err, "failed to extract batch size")
	}
	start := time.Now()
	var labels []*api.Label
	channel, ok := mux.Vars(r)["channel"]
	if !ok || channel == "default" {
		// any tags are ok
		labels, err = s.labelStore.Sample(
			r.Context(),
			s.projectID,
			batchSize,
		)
		if err != nil {
			return errors.Wrap(err, "labelstore.Sample")
		}
	} else {
		// only labels with a tag matching the channel name are ok
		labels, err = s.labelStore.SampleWithTags(
			r.Context(),
			s.projectID,
			batchSize,
			strings.Split(channel, ","),
			r.URL.Query().Get("all") == "1",
		)
		if err != nil {
			return errors.Wrap(err, "labelstore.SampleWithTags")
		}
	}
	s.log.Debug("sampled output",
		zap.Int("batchSize", batchSize),
		base.Elapsed(start))
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(labels)
}

// SampleLabels samples labels from the output channel.
// This is gadget-specific as it requires querying the
// gadget's underlying storage.
func (s *Filter) handleOutput() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			switch r.Method {
			case http.MethodGet:
				return s.handleGetOutput(w, r)
			case http.MethodPut:
				return s.handlePutLabel(w, r)
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
				return errors.New("method not allowed")
			}
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error(r.RequestURI, zap.Error(err))
		}
	}
}

// handlePutLabel handler for inserting a label into the
// output channel. The frontend invokes this to insert
// labels into the host gadget's underlying storage.
func (s *Filter) handlePutLabel(
	w http.ResponseWriter,
	r *http.Request,
) error {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid content type")
	}
	var label api.Label
	if err := json.NewDecoder(r.Body).Decode(&label); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.Wrap(err, "failed to decode label")
	}
	if err := s.labelStore.Insert(&label); err != nil {
		return errors.Wrap(err, "failed to insert label")
	}
	s.log.Debug("inserted label",
		zap.String("projectID", label.ProjectID))
	return nil
}

// handleGetLabel handler for retrieving a label from the
// output channel. The frontend invokes this to retrieve
// a specific label's metadata from the host gadget's underlying
// storage given its id.
func (s *Filter) handleGetLabel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return errors.New("method not allowed")
			}
			id, ok := mux.Vars(r)["id"]
			if !ok || id == "" {
				w.WriteHeader(http.StatusBadRequest)
				return errors.New("missing id")
			}
			labels, err := s.labelStore.Get(r.Context(), id)
			if err == labelstore.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return err
			} else if err != nil {
				return errors.Wrap(err, "failed to list labels")
			}
			w.Header().Set("Content-Type", "application/json")
			return json.NewEncoder(w).Encode(labels)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error(r.RequestURI, zap.Error(err))
		}
	}
}
