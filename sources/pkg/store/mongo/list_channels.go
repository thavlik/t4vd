package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *mongoStore) ListChannels(
	ctx context.Context,
	projectID string,
) ([]*api.Channel, error) {
	cursor, err := s.channels.Find(
		ctx,
		map[string]interface{}{
			"project": projectID,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	channels := make([]*api.Channel, len(docs))
	for i, doc := range docs {
		channels[i] = convertChannelDoc(doc)
	}
	return channels, nil
}

func (s *mongoStore) ListChannelIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	return getFieldFromDocs(
		ctx,
		projectID,
		s.channels,
		"c",
		blacklist,
	)
}

func convertChannelDoc(m map[string]interface{}) *api.Channel {
	return &api.Channel{
		ID:        m["c"].(string),
		Blacklist: m["blacklist"].(bool),
	}
}
