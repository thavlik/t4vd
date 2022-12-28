package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *mongoStore) ListProjectsCreatedBy(
	ctx context.Context,
	userID string,
) ([]*api.Project, error) {
	query := map[string]interface{}{
		"creator": userID,
	}
	cursor, err := s.projects.Find(ctx, query)
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
