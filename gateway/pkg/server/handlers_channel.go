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
)

func (s *Server) handleAddChannel() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			var req sources.AddChannelRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(ctx, userID, req.ProjectID); err != nil {

				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			req.SubmitterID = userID
			resp, err := s.sources.AddChannel(ctx, req)
			if err != nil {
				return errors.Wrap(err, "sources.AddChannel")
			}
			output, err := resolveChannel(ctx, s.seerOpts, resp)
			if err != nil {
				return errors.Wrap(err, "resolveChannel")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(output); err != nil {
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
			ctx := r.Context()
			var req sources.RemoveChannelRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(ctx, userID, req.ProjectID); err != nil {

				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.sources.RemoveChannel(ctx, req)
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
			ctx := r.Context()
			projectID := r.URL.Query().Get("p")
			if projectID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			if err := s.ProjectAccess(ctx, userID, projectID); err != nil {

				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.sources.ListChannels(
				ctx,
				sources.ListChannelsRequest{
					ProjectID: projectID,
				})
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			output, err := resolveBulkChannels(
				ctx,
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

func resolveChannel(
	ctx context.Context,
	opts base.ServiceOptions,
	ch *sources.Channel,
) (*channel, error) {
	output := &channel{
		ID:        ch.ID,
		Blacklist: ch.Blacklist,
	}
	resp, err := seer.NewSeerClientFromOptions(opts).
		GetChannelDetails(
			ctx,
			seer.GetChannelDetailsRequest{
				Input: ch.ID,
			},
		)
	if err != nil {
		if strings.Contains(err.Error(), infocache.ErrCacheUnavailable.Error()) {
			return output, nil
		}
		return nil, errors.Wrap(err, "seer.GetBulkChannelsDetails")
	}
	output.Info = &resp.Details
	return output, nil
}

func resolveBulkChannels(
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
