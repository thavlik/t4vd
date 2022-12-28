package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	seer "github.com/thavlik/bjjvb/seer/pkg/api"
	sources "github.com/thavlik/bjjvb/sources/pkg/api"
)

func (s *Server) handleAddChannel() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req sources.AddChannelRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			req.SubmitterID = userID
			resp, err := s.sources.AddChannel(context.Background(), req)
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleRemoveChannel() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req sources.RemoveChannelRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			resp, err := s.sources.RemoveChannel(context.Background(), req)
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleListChannels() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			projectID := r.URL.Query().Get("p")
			if projectID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			resp, err := s.sources.ListChannels(r.Context(), sources.ListChannelsRequest{
				ProjectID: projectID,
			})
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp.Channels); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleGetChannelAvatar() http.HandlerFunc {
	return s.handler(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) error {
			channelID := r.URL.Query().Get("c")
			if channelID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			w.Header().Set("Content-Type", "image/jpeg")
			if err := seer.GetChannelAvatar(
				r.Context(),
				s.seerOpts,
				channelID,
				w,
			); err == seer.ErrNotCached {
				w.WriteHeader(http.StatusNotFound)
				return nil
			} else if err != nil {
				return errors.Wrap(err, "seer.GetChannelAvatar")
			}
			return nil
		})
}
