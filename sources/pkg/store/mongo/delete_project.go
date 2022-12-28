package mongo

import (
	"context"

	"github.com/pkg/errors"
)

func (s *mongoStore) DeleteProject(id string) error {
	if _, err := s.projects.DeleteOne(
		context.Background(),
		map[string]interface{}{
			"_id": id,
		},
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
