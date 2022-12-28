package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) CreateProject(
	project *api.Project,
) error {
	var v struct {
		CreatorID string `bson:"creator"`
	}
	if err := s.projects.FindOne(context.Background(),
		map[string]interface{}{
			"_id": project.ID,
		},
	).Decode(&v); err == nil {
		if v.CreatorID != project.CreatorID {
			return errors.New("only the project creator can change the name")
		}
	} else if err != mongo.ErrNoDocuments {
		return errors.Wrap(err, "mongo")
	}
	if _, err := s.projects.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": project.ID,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"name":    project.Name,
				"creator": project.CreatorID,
				"group":   project.GroupID,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
