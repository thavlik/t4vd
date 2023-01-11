package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
	"go.uber.org/zap"
)

func (s *Server) handleCreateProject() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req sources.Project
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if req.Name == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			req.ID = ""            // disallow explicit
			req.GroupID = ""       // disallow explicit
			req.CreatorID = userID // enforced by rbac
			resp, err := s.sources.CreateProject(context.Background(), req)
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

func (s *Server) handleDeleteProject() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req sources.DeleteProject
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			// TODO: fix RBAC
			resp, err := s.sources.DeleteProject(context.Background(), req)
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

func (s *Server) handleGetProject() http.HandlerFunc {
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
			project, err := s.sources.GetProject(
				r.Context(),
				sources.GetProject{
					ID: projectID,
				})
			if err != nil {
				if strings.Contains(err.Error(), store.ErrProjectNotFound.Error()) {
					w.WriteHeader(http.StatusNotFound)
					return nil
				}
				return errors.Wrap(err, "sources")
			}
			users, err := s.iam.ListGroupMembers(
				r.Context(),
				project.GroupID,
			)
			if err != nil {
				return errors.Wrap(err, "iam.ListGroupMembers")
			}
			collabs := make([]*collaborator, len(users))
			for i, user := range users {
				collabs[i] = &collaborator{
					ID:       user.ID,
					Username: user.Username,
				}
			}
			s.log.Debug("retrieved project",
				zap.String("id", project.ID),
				zap.String("name", project.Name),
				zap.Int("numCollabs", len(collabs)))
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(&getProjectResponse{
				ID:            project.ID,
				Name:          project.Name,
				CreatorID:     project.CreatorID,
				GroupID:       project.GroupID,
				Collaborators: collabs,
			}); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

type collaborator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type getProjectResponse struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	CreatorID     string          `json:"creatorID"`
	GroupID       string          `json:"groupID"`
	Collaborators []*collaborator `json:"collaborators"`
}

func (s *Server) handleListProjects() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodGet,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			resp, err := s.sources.ListProjects(r.Context(), sources.ListProjectsRequest{
				VisibleToUserID: userID,
			})
			if err != nil {
				return errors.Wrap(err, "sources")
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(resp.Projects); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}

func (s *Server) handleProjectAddCollaborator() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req struct {
				UserID    string `json:"userID"`
				ProjectID string `json:"projectID"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(r.Context(), userID, req.ProjectID); err != nil {

				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			project, err := s.sources.GetProject(
				context.Background(),
				sources.GetProject{ID: req.ProjectID})
			if err != nil {
				return errors.Wrap(err, "sources.GetProject")
			}
			// TODO: make sure user has additional permissions to add collaborator
			// admin permissions?
			if err := s.iam.AddUserToGroup(
				userID,
				project.GroupID,
			); err != nil {
				return errors.Wrap(err, "iam.AddUserToGroup")
			}
			return nil
		})
}

func (s *Server) handleProjectRemoveCollaborator() http.HandlerFunc {
	return s.rbacHandler(
		http.MethodPost,
		iam.NullPermissions,
		func(userID string, w http.ResponseWriter, r *http.Request) error {
			var req struct {
				UserID    string `json:"userID"`
				ProjectID string `json:"projectID"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return errors.Wrap(err, "decoder")
			}
			if err := s.ProjectAccess(r.Context(), userID, req.ProjectID); err != nil {

				w.WriteHeader(http.StatusForbidden)
				return nil
			}
			project, err := s.sources.GetProject(
				context.Background(),
				sources.GetProject{ID: req.ProjectID})
			if err != nil {
				return errors.Wrap(err, "sources.GetProject")
			}
			// TODO: make sure user has permission to remove collaborator
			// admin permissions?
			if err := s.iam.RemoveUserFromGroup(
				userID,
				project.GroupID,
			); err != nil {
				return errors.Wrap(err, "iam.RemoveUserFromGroup")
			}
			return nil
		})
}

func (s *Server) handleProjectExists() http.HandlerFunc {
	return s.handler(
		http.MethodGet,
		func(w http.ResponseWriter, r *http.Request) (err error) {
			name := r.URL.Query().Get("n")
			if name == "" {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			var res struct {
				Exists bool `json:"exists"`
			}
			res.Exists, err = func() (bool, error) {
				if _, err := s.sources.GetProjectByName(
					r.Context(),
					sources.GetProjectByName{
						Name: name,
					},
				); err != nil {
					if strings.Contains(err.Error(), store.ErrProjectNotFound.Error()) {
						return false, nil
					}
					return false, errors.Wrap(err, "sources")
				}
				return true, nil
			}()
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(&res); err != nil {
				return errors.Wrap(err, "encoder")
			}
			return nil
		})
}
