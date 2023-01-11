package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/iam"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
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
			if err := s.ProjectAccess(r.Context(), userID, req.ProjectID); err != nil {
				s.log.Warn("project access denied", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			req.SubmitterID = userID
			resp, err := s.sources.AddChannel(context.Background(), req)
			if err != nil {
				if strings.Contains(err.Error(), infocache.ErrCacheUnavailable.Error()) {
					// TODO: this is a hack, we should have a better way to handle this
					// In the future, the handler should wait to see if the download
					// completes before returning Accepted, which indicates that the
					// request was accepted but not yet completed.
					w.WriteHeader(http.StatusAccepted)
					return nil
				}
				return errors.Wrap(err, "sources.AddChannel")
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
			if err := s.ProjectAccess(r.Context(), userID, req.ProjectID); err != nil {
				s.log.Warn("project access denied", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return nil
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
			if err := s.ProjectAccess(r.Context(), userID, projectID); err != nil {
				s.log.Warn("project access denied", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.sources.ListChannels(
				r.Context(),
				sources.ListChannelsRequest{
					ProjectID: projectID,
				})
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			output, err := resolveChannels(
				r.Context(),
				s.seerOpts,
				resp.Channels,
			)
			if err != nil {
				return errors.Wrap(err, "resolveChannels")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(output); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

type channel struct {
	ID        string               `json:"id"`
	Blacklist bool                 `json:"blacklist"`
	Info      *seer.ChannelDetails `json:"info"`
}

func resolveChannels(
	ctx context.Context,
	opts base.ServiceOptions,
	channels []*sources.Channel,
) ([]*channel, error) {
	output := make([]*channel, len(channels))
	channelIDs := make([]string, len(channels))
	for i, v := range channels {
		channelIDs[i] = v.ID
		output[i] = &channel{
			ID:        v.ID,
			Blacklist: v.Blacklist,
		}
	}
	resolved, err := seer.NewSeerClientFromOptions(opts).
		GetBulkChannelsDetails(
			ctx,
			seer.GetBulkChannelsDetailsRequest{
				ChannelIDs: channelIDs,
			},
		)
	if err != nil {
		return nil, errors.Wrap(err, "seer.GetBulkChannelsDetails")
	}
	for _, info := range resolved.Channels {
		for i, ch := range channels {
			if ch.ID == info.ID {
				output[i].Info = info
				break
			}
		}
	}
	return output, nil
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
				//w.Header().Set("Content-Type", "image/svg+xml")
				//_, err := w.Write(pendingSvg)
				//return err
			} else if err != nil {
				return errors.Wrap(err, "seer.GetChannelAvatar")
			}
			return nil
		})
}
