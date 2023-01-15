package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	filter "github.com/thavlik/t4vd/filter/pkg/api"
	slideshow "github.com/thavlik/t4vd/slideshow/pkg/api"
)

func (s *Server) handleGetFilterStack() http.HandlerFunc {
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
			var size int64 = 5
			if s := r.URL.Query().Get("s"); s != "" {
				var err error
				size, err = strconv.ParseInt(s, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return nil
				} else if size > 10 {
					size = 10
				}
			}
			resp, err := slideshow.NewSlideShowClientFromOptions(
				s.slideshow,
			).GetRandomStack(
				r.Context(),
				slideshow.GetRandomStack{
					ProjectID: projectID,
					Size:      int(size),
				})
			if err != nil {
				return errors.Wrap(err, "filter")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleFilterClassify() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			ctx := r.Context()
			var label filter.Label
			if err := json.NewDecoder(r.Body).Decode(&label); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			if err := s.ProjectAccess(ctx, userID, label.ProjectID); err != nil {
				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			resp, err := s.filter.Classify(ctx, label)
			if err != nil {
				return errors.Wrap(err, "filter.Classify")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}
