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

func handleVideos(
	s api.Sources,
	projectID string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() (err error) {
			switch r.Method {
			case http.MethodGet:
				retCode, err = handleGetVideos(w, r, s, projectID, log)
			case http.MethodPut:
				retCode, err = handlePutVideo(w, r, s, projectID, log)
			case http.MethodDelete:
				retCode, err = handleDeleteVideo(w, r, s, projectID, log)
			default:
				retCode, err = http.StatusMethodNotAllowed, base.InvalidMethod(r.Method)
			}
			return err
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}

func handleGetVideos(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (int, error) {
	videos, err := s.ListVideos(
		r.Context(),
		api.ListVideosRequest{
			ProjectID: projectID,
		})
	if err != nil {
		return 0, errors.Wrap(err, "sources.ListVideos")
	}
	w.Header().Set("Content-Type", "application/json")
	return 0, json.NewEncoder(w).Encode(videos)
}

func handlePutVideo(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (int, error) {
	var req api.AddVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json.Decode")
	}
	req.ProjectID = projectID
	video, err := s.AddVideo(
		r.Context(),
		req,
	)
	if err != nil {
		return 0, errors.Wrap(err, "store.AddVideo")
	}
	w.Header().Set("Content-Type", "application/json")
	return 0, json.NewEncoder(w).Encode(video)
}

func handleDeleteVideo(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (int, error) {
	var req api.RemoveVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json.Decode")
	}
	req.ProjectID = projectID
	if _, err := s.RemoveVideo(
		r.Context(),
		req,
	); err != nil {
		return 0, errors.Wrap(err, "store.RemoveVideo")
	}
	return 0, nil
}
