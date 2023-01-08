package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *mongoStore) GetProjectIDsForVideo(
	ctx context.Context,
	videoID string,
) (projectIDs []string, err error) {
	cursor, err := s.videos.Find(
		ctx,
		map[string]interface{}{
			"v": videoID,
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
