package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/sources/pkg/api"
	"github.com/thavlik/bjjvb/sources/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ids, err := getDocIDs(ctx, projectID, s.channels, blacklist)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func getDocIDs(
	ctx context.Context,
	projectID string,
	col *mongo.Collection,
	blacklist bool,
) ([]string, error) {
	cursor, err := col.Find(
		ctx,
		map[string]interface{}{
			"blacklist": blacklist,
			"project":   projectID,
		},
		options.Find().SetProjection(map[string]interface{}{
			"_id": 1,
		}))
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	IDs := make([]string, len(docs))
	for i, doc := range docs {
		IDs[i] = store.ExtractResourceID(doc["_id"].(string))
	}
	return IDs, nil
}

func convertChannelDoc(m map[string]interface{}) *api.Channel {
	return &api.Channel{
		ID:        store.ExtractResourceID(m["_id"].(string)),
		Name:      m["name"].(string),
		Avatar:    m["avatar"].(string),
		Subs:      m["subs"].(string),
		Blacklist: m["blacklist"].(bool),
	}
}
