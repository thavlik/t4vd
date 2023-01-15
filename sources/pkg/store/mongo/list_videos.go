package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *mongoStore) ListVideos(
	ctx context.Context,
	projectID string,
) ([]*api.Video, error) {
	cursor, err := s.videos.Find(
		ctx,
		map[string]interface{}{
			"project": projectID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	videos := make([]*api.Video, len(docs))
	for i, doc := range docs {
		videos[i] = convertVideoDoc(doc)
	}
	return videos, nil
}

func (s *mongoStore) ListVideoIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	return getFieldFromDocs(
		ctx,
		projectID,
		s.videos,
		"v",
		blacklist,
	)
}

func convertVideoDoc(m map[string]interface{}) *api.Video {
	return &api.Video{
		ID:          m["v"].(string),
		Blacklist:   m["blacklist"].(bool),
		SubmitterID: m["submitter"].(string),
		Submitted:   m["submitted"].(int64),
	}
}
