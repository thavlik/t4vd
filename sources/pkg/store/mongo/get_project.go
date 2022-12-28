package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *mongoStore) GetProject(
	ctx context.Context,
	projectID string,
) (*api.Project, error) {
	doc := make(map[string]interface{})
	if err := s.projects.FindOne(
		ctx,
		map[string]interface{}{
			"_id": projectID,
		},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, store.ErrProjectNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	return convertProjectDoc(doc), nil
}
