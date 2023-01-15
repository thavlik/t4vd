package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) CreateProject(
	project *api.Project,
) error {
	if _, err := s.projects.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": project.ID,
		},
		map[string]interface{}{
			"$setOnInsert": map[string]interface{}{
				"creator": project.CreatorID,
				"created": project.Created,
			},
			"$set": map[string]interface{}{
				"name":  project.Name,
				"group": project.GroupID,
				"tags":  project.Tags,
				"desc":  project.Description,
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
