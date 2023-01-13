package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getRecency(
	col *mongo.Collection,
	id string,
) (time.Time, error) {
	doc := make(map[string]interface{})
	if err := col.FindOne(
		context.Background(),
		map[string]interface{}{
			"_id": id,
		},
	).Decode(&doc); err == mongo.ErrNoDocuments {
		return time.Time{}, infocache.ErrCacheUnavailable
	} else if err != nil {
		return time.Time{}, errors.Wrap(err, "mongo")
	}
	v, ok := doc["updated"].(int64)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid type for field 'updated': %T", doc["updated"])
	}
	updated := time.Unix(0, v)
	return updated, nil
}

func checkCacheRecency(
	col *mongo.Collection,
	id string,
) (bool, error) {
	updated, err := getRecency(col, id)
	if err == infocache.ErrCacheUnavailable {
		return false, nil
	} else if err != nil {
		return false, errors.Wrap(err, "getCacheRecency")
	}
	return time.Since(updated) < infocache.CacheRecency, nil
}

func getVideoIDs(
	ctx context.Context,
	joins *mongo.Collection,
	keyName string,
	value string,
) ([]string, error) {
	cursor, err := joins.Find(
		ctx,
		map[string]interface{}{
			keyName: value,
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	if len(docs) == 0 {
		return nil, infocache.ErrCacheUnavailable
	}
	videoIDs := make([]string, len(docs))
	for i, doc := range docs {
		videoIDs[i] = doc["v"].(string)
	}
	return videoIDs, nil
}

func refreshCache(
	col *mongo.Collection,
	id string,
	timestamp time.Time,
) error {
	if _, err := col.UpdateOne(
		context.Background(),
		map[string]interface{}{
			"_id": id,
		},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"updated": timestamp.UnixNano(),
			},
		},
		options.Update().SetUpsert(true),
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}

func setVideoIDs(
	joins *mongo.Collection,
	originKey string,
	originValue string,
	videoIDs []string,
) error {
	operations := make([]mongo.WriteModel, len(videoIDs)+1)
	operations[0] = mongo.NewDeleteManyModel().
		SetFilter(map[string]interface{}{
			originKey: originValue,
		})
	for i, videoID := range videoIDs {
		operations[i+1] = mongo.NewInsertOneModel().
			SetDocument(map[string]interface{}{
				originKey: originValue,
				"v":       videoID,
			})
	}
	_, err := joins.BulkWrite(context.Background(), operations)
	if err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
