package mongo

import (
	"context"

	"github.com/hashicorp/go-multierror"
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
	var multi error
	if _, err := s.videos.DeleteMany(
		context.Background(),
		map[string]interface{}{
			"project": id,
		},
	); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "delete videos"))
	}
	if _, err := s.playlists.DeleteMany(
		context.Background(),
		map[string]interface{}{
			"project": id,
		},
	); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "delete playlists"))
	}
	if _, err := s.channels.DeleteMany(
		context.Background(),
		map[string]interface{}{
			"project": id,
		},
	); err != nil {
		multi = multierror.Append(multi, errors.Wrap(err, "delete channels"))
	}
	return multi
}
