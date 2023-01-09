package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) GetProjectIDsForPlaylist(
	ctx context.Context,
	playlistID string,
) (projectIDs []string, err error) {
	cursor, err := s.channels.Find(
		ctx,
		map[string]interface{}{
			"p":         playlistID,
			"blacklist": false,
		},
		options.Find().SetProjection(map[string]interface{}{
			"project": 1,
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	return decodeProjectIDsWithContext(ctx, cursor)
}
