package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
)

func (s *Server) handleGetDataset() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			projectID := r.URL.Query().Get("p")
			if projectID == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			resp, err := s.compiler.GetDataset(r.Context(), compiler.GetDatasetRequest{
				ProjectID: projectID,
			})
			if err != nil {
				if strings.Contains(err.Error(), datastore.ErrDatasetNotFound.Error()) {
					w.WriteHeader(http.StatusNotFound)
					return nil
				}
				return errors.Wrap(err, "compiler")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}
