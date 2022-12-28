package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	slideshow "github.com/thavlik/bjjvb/slideshow/pkg/api"
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
				return errors.Wrap(err, "parse time query")
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
