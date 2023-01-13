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

func (s *Server) handleAddVideo() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			var req sources.AddVideoRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(ctx, userID, req.ProjectID); err != nil {
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			req.SubmitterID = userID
			resp, err := s.sources.AddVideo(ctx, req)
			if err != nil {
				return errors.Wrap(err, "sources.AddVideo")
			}
			output, err := resolveVideo(ctx, s.seerOpts, resp)
			if err != nil {
				return errors.Wrap(err, "resolveVideo")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(output); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleRemoveVideo() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req sources.RemoveVideoRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(r.Context(), userID, req.ProjectID); err != nil {
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.sources.RemoveVideo(context.Background(), req)
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

func (s *Server) handleListVideos() http.HandlerFunc {
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
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.sources.ListVideos(
				r.Context(),
				sources.ListVideosRequest{
					ProjectID: projectID,
				})
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			output, err := resolveBulkVideos(
				r.Context(),
				s.seerOpts,
				resp.Videos,
			)
			if err != nil {
				return errors.Wrap(err, "resolveVideos")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(output); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

type video struct {
	ID        string             `json:"id"`
	Blacklist bool               `json:"blacklist"`
	Info      *seer.VideoDetails `json:"info"`
}

func resolveVideo(
	ctx context.Context,
	opts base.ServiceOptions,
	v *sources.Video,
) (*video, error) {
	output := &video{
		ID:        v.ID,
		Blacklist: v.Blacklist,
	}
	resolved, err := seer.NewSeerClientFromOptions(opts).
		GetVideoDetails(
			ctx,
			seer.GetVideoDetailsRequest{
				Input: v.ID,
			})
	if err != nil {
		if strings.Contains(err.Error(), infocache.ErrCacheUnavailable.Error()) {
			return output, nil
		}
		return nil, errors.Wrap(err, "seer.GetVideoDetails")
	}
	output.Info = &resolved.Details
	return output, nil
}

func resolveBulkVideos(
	ctx context.Context,
	opts base.ServiceOptions,
	videos []*sources.Video,
) ([]*video, error) {
	output := make([]*video, len(videos))
	videoIDs := make([]string, len(videos))
	for i, v := range videos {
		videoIDs[i] = v.ID
		output[i] = &video{
			ID:        v.ID,
			Blacklist: v.Blacklist,
		}
	}
	resolved, err := seer.NewSeerClientFromOptions(opts).
		GetBulkVideosDetails(
			ctx,
			seer.GetBulkVideosDetailsRequest{
				VideoIDs: videoIDs,
			},
		)
	if err != nil {
		return nil, errors.Wrap(err, "seer.GetBulkPlaylistsDetails")
	}
	for _, info := range resolved.Videos {
		for i, v := range videos {
			if v.ID == info.ID {
				output[i].Info = info
				break
			}
		}
	}
	return output, nil
}

func (s *Server) handleGetVideoThumbnail() http.HandlerFunc {
	return s.handler(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) error {
			videoID := r.URL.Query().Get("v")
			if videoID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			w.Header().Set("Content-Type", "image/jpeg")
			if err := seer.GetVideoThumbnail(
				r.Context(),
				s.seerOpts,
				videoID,
				w,
			); err == seer.ErrNotCached {
				w.WriteHeader(http.StatusNotFound)
				return nil
				//w.Header().Set("Content-Type", "image/svg+xml")
				//_, err := w.Write(pendingSvg)
				//return err
			} else if err != nil {
				return errors.Wrap(err, "seer.GetVideoThumbnail")
			}
			return nil
		})
}
