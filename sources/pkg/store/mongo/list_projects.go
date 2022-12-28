package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
)

func (s *mongoStore) ListProjects(
	ctx context.Context,
) ([]*api.Project, error) {
	cursor, err := s.projects.Find(ctx, map[string]interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	projects := make([]*api.Project, len(docs))
	for i, doc := range docs {
		projects[i] = convertProjectDoc(doc)
	}
	return projects, nil
}

func convertProjectDoc(m map[string]interface{}) *api.Project {
	return &api.Project{
		ID:        m["_id"].(string),
		Name:      m["name"].(string),
		CreatorID: m["creator"].(string),
		GroupID:   m["group"].(string),
	}
}
