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

func handlePlaylists(
	s api.Sources,
	projectID string,
	log *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retCode := http.StatusInternalServerError
		if err := func() (err error) {
			switch r.Method {
			case http.MethodGet:
				retCode, err = handleGetPlaylists(w, r, s, projectID, log)
			case http.MethodPut:
				retCode, err = handlePutPlaylist(w, r, s, projectID, log)
			case http.MethodDelete:
				retCode, err = handleDeletePlaylist(w, r, s, projectID, log)
			default:
				retCode, err = http.StatusMethodNotAllowed, base.InvalidMethod(r.Method)
			}
			return err
		}(); err != nil {
			gadget.HandlerError(r, w, retCode, err, log)
		}
	}
}

func handleGetPlaylists(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (statusCode int, err error) {
	playlists, err := s.ListPlaylists(
		r.Context(),
		api.ListPlaylistsRequest{
			ProjectID: projectID,
		},
	)
	if err != nil {
		return 0, errors.Wrap(err, "sources.ListPlaylists")
	}
	w.Header().Set("Content-Type", "application/json")
	return 0, json.NewEncoder(w).Encode(playlists)
}

func handlePutPlaylist(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (statusCode int, err error) {
	var req api.AddPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json.Decode")
	}
	req.ProjectID = projectID
	playlist, err := s.AddPlaylist(
		r.Context(),
		req,
	)
	if err != nil {
		return 0, errors.Wrap(err, "sources.AddPlaylist")
	}
	w.Header().Set("Content-Type", "application/json")
	return 0, json.NewEncoder(w).Encode(playlist)
}

func handleDeletePlaylist(
	w http.ResponseWriter,
	r *http.Request,
	s api.Sources,
	projectID string,
	log *zap.Logger,
) (statusCode int, err error) {
	var req api.RemovePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return http.StatusBadRequest, errors.Wrap(err, "json.Decode")
	}
	req.ProjectID = projectID
	if _, err := s.RemovePlaylist(
		r.Context(),
		req,
	); err != nil {
		return 0, errors.Wrap(err, "sources.RemovePlaylist")
	}
	return 0, nil
}
