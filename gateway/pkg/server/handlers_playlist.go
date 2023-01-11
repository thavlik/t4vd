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

func (s *Server) handleAddPlaylist() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			var req sources.AddPlaylistRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(ctx, userID, req.ProjectID); err != nil {
				s.log.Warn("project access denied", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			req.SubmitterID = userID
			resp, err := s.sources.AddPlaylist(ctx, req)
			if err != nil {
				return errors.Wrap(err, "sources.AddPlaylist")
			}
			output, err := resolvePlaylist(ctx, s.seerOpts, resp)
			if err != nil {
				return errors.Wrap(err, "resolvePlaylist")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(output); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleRemovePlaylist() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req sources.RemovePlaylistRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(r.Context(), userID, req.ProjectID); err != nil {
				s.log.Warn("project access denied", zap.Error(err))
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.sources.RemovePlaylist(context.Background(), req)
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

func (s *Server) handleListPlaylists() http.HandlerFunc {
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
			resp, err := s.sources.ListPlaylists(
				r.Context(),
				sources.ListPlaylistsRequest{
					ProjectID: projectID,
				})
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			output, err := resolveBulkPlaylists(
				r.Context(),
				s.seerOpts,
				resp.Playlists,
			)
			if err != nil {
				return errors.Wrap(err, "resolvePlaylists")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(output); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func resolvePlaylist(
	ctx context.Context,
	opts base.ServiceOptions,
	pl *sources.Playlist,
) (*playlist, error) {
	output := &playlist{
		ID:        pl.ID,
		Blacklist: pl.Blacklist,
	}
	resolved, err := seer.NewSeerClientFromOptions(opts).
		GetPlaylistDetails(
			ctx,
			seer.GetPlaylistDetailsRequest{
				Input: pl.ID,
			},
		)
	if err != nil {
		if strings.Contains(err.Error(), infocache.ErrCacheUnavailable.Error()) {
			return output, nil
		}
		return nil, errors.Wrap(err, "seer.GetPlaylistDetails")
	}
	output.Info = &resolved.Details
	return output, nil
}

func resolveBulkPlaylists(
	ctx context.Context,
	opts base.ServiceOptions,
	playlists []*sources.Playlist,
) ([]*playlist, error) {
	output := make([]*playlist, len(playlists))
	playlistIDs := make([]string, len(playlists))
	for i, v := range playlists {
		playlistIDs[i] = v.ID
		output[i] = &playlist{
			ID:        v.ID,
			Blacklist: v.Blacklist,
		}
	}
	resolved, err := seer.NewSeerClientFromOptions(opts).
		GetBulkPlaylistsDetails(
			ctx,
			seer.GetBulkPlaylistsDetailsRequest{
				PlaylistIDs: playlistIDs,
			},
		)
	if err != nil {
		return nil, errors.Wrap(err, "seer.GetBulkPlaylistsDetails")
	}
	for _, info := range resolved.Playlists {
		for i, pl := range playlists {
			if pl.ID == info.ID {
				output[i].Info = info
				break
			}
		}
	}
	return output, nil
}

type playlist struct {
	ID        string                `json:"id"`
	Blacklist bool                  `json:"blacklist"`
	Info      *seer.PlaylistDetails `json:"info"`
}

func (s *Server) handleGetPlaylistThumbnail() http.HandlerFunc {
	return s.handler(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) error {
			playlistID := r.URL.Query().Get("list")
			if playlistID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			w.Header().Set("Content-Type", "image/jpeg")
			if err := seer.GetPlaylistThumbnail(
				r.Context(),
				s.seerOpts,
				playlistID,
				w,
			); err == seer.ErrNotCached {
				w.WriteHeader(http.StatusNotFound)
				return nil
				//w.Header().Set("Content-Type", "image/svg+xml")
				//_, err := w.Write(pendingSvg)
				//return err
			} else if err != nil {
				return errors.Wrap(err, "seer.GetPlaylistThumbnail")
			}
			return nil
		})
}
