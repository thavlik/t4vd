package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *Server) handleGetRandomMarker() http.HandlerFunc {
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
			if resp, err := s.sources.IsProjectEmpty(r.Context(), sources.IsProjectEmptyRequest{
				ProjectID: projectID,
			}); err != nil {
				return errors.Wrap(err, "sources")
			} else if resp.IsEmpty {
				// Project is empty, return 404.
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("project is empty"))
				return nil
			}
			resp, err := slideshow.NewSlideShowClientFromOptions(s.slideshow).
				GetRandomMarker(
					r.Context(),
					slideshow.GetRandomMarker{
						ProjectID: projectID,
					})
			if err != nil {
				return errors.Wrap(err, "slideshow")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleGetFrame() http.HandlerFunc {
	// Skip rbac for performance purposes, but it
	// should be enabled if security is an issue.
	return s.handler(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) error {
			videoID := r.URL.Query().Get("v")
			if videoID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			tv, err := strconv.ParseInt(r.URL.Query().Get("t"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			t := time.Duration(tv)
			w.Header().Set("Content-Type", "image/jpeg")
			if err := slideshow.GetFrame(
				r.Context(),
				s.slideshow,
				videoID,
				t,
				w,
			); err != nil {
				return errors.Wrap(err, "slideshow.GetFrame")
			}
			return nil
		})
}
