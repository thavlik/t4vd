package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

func handleChannels(
	s api.Sources,
	projectID string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() (err error) {
			switch r.Method {
			case http.MethodGet:
				retCode, err = handleGetChannels(w, r, s, projectID, log)
			case http.MethodPut:
				retCode, err = handlePutChannel(w, r, s, projectID, log)
			case http.MethodDelete:
				retCode, err = handleDeleteChannel(w, r, s, projectID, log)
			default:
				retCode, err = http.StatusMethodNotAllowed, base.InvalidMethod(r.Method)
			}
			return err
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}

func handleGetChannels(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (int, error) {
	channels, err := s.ListChannels(
		r.Context(),
		api.ListChannelsRequest{
			ProjectID: projectID,
		},
	)
	if err != nil {
		return 0, errors.Wrap(err, "sources.ListChannels")
	}
	w.Header().Set("Content-Type", "application/json")
	return 0, json.NewEncoder(w).Encode(channels)
}

func handlePutChannel(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (int, error) {
	var req api.AddChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json.Decode")
	}
	req.ProjectID = projectID
	channel, err := s.AddChannel(
		r.Context(),
		req,
	)
	if err != nil {
		return 0, errors.Wrap(err, "store.AddChannel")
	}
	w.Header().Set("Content-Type", "application/json")
	return 0, json.NewEncoder(w).Encode(channel)
}

func handleDeleteChannel(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (int, error) {
	var req api.RemoveChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json.Decode")
	}
	req.ProjectID = projectID
	if _, err := s.RemoveChannel(
		r.Context(),
		req,
	); err != nil {
		return 0, errors.Wrap(err, "store.RemoveChannel")
	}
	return 0, nil
}
