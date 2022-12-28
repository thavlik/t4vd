package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ds *mongoDataStore) GetDataset(
	ctx context.Context,
	projectID string,
	datasetID string,
) (*api.Dataset, error) {
	if datasetID == "" {
		// prefer latest complete dataset
		return ds.getLatestDataset(ctx, projectID, true)
	}
	return ds.getSpecificDataset(ctx, projectID, datasetID)
}

func (ds *mongoDataStore) getDatasetVideoIDs(
	ctx context.Context,
	datasetID string,
) ([]string, error) {
	query, err := ds.outputVideos.Find(
		ctx,
		map[string]interface{}{
			"ds": datasetID,
		},
		options.Find().SetProjection(map[string]interface{}{
			"v": 1,
		}))
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := query.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	videoIDs := make([]string, len(docs))
	for i, doc := range docs {
		videoIDs[i] = doc["v"].(string)
	}
	return videoIDs, nil
}

func (ds *mongoDataStore) getLatestDataset(
	ctx context.Context,
	projectID string,
	complete bool,
) (*api.Dataset, error) {
	// Use the latest complete dataset
	result := ds.outputDatasets.FindOne(
		ctx,
		map[string]interface{}{
			"p": projectID,
			"c": complete,
		},
		options.FindOne().
			SetSort(map[string]interface{}{
				"t": -1,
			}).
			SetProjection(map[string]interface{}{
				"_id": 1,
				"t":   1,
			}))
	if err := result.Err(); err == mongo.ErrNoDocuments {
		if complete {
			// Try and get the latest incomplete dataset
			return ds.getLatestDataset(ctx, projectID, false)
		}
		return nil, datastore.ErrDatasetNotFound
	}
	doc := make(map[string]interface{})
	if err := result.Decode(&doc); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	datasetID := doc["_id"].(string)
	return ds.getSpecificDataset(ctx, projectID, datasetID)
}

func (ds *mongoDataStore) getSpecificDataset(
	ctx context.Context,
	projectID string,
	datasetID string,
) (*api.Dataset, error) {
	result := ds.outputDatasets.FindOne(
		ctx,
		map[string]interface{}{
			"_id": datasetID,
			"p":   projectID,
		})
	if err := result.Err(); err == mongo.ErrNoDocuments {
		return nil, datastore.ErrDatasetNotFound
	}
	doc := make(map[string]interface{})
	if err := result.Decode(&doc); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	timestamp := doc["t"].(int64)
	complete := doc["c"].(bool)
	videoIDs, err := ds.getDatasetVideoIDs(ctx, datasetID)
	if err != nil {
		return nil, errors.Wrap(err, "getDatasetVideoIDs")
	}
	videos, err := datastore.ResolveVideos(
		ctx,
		ds.seer,
		ds,
		videoIDs,
		nil,
		ds.log,
	)
	if err != nil {
		return nil, errors.Wrap(err, "ResolveVideos")
	}
	return &api.Dataset{
		ID:        datasetID,
		Timestamp: timestamp,
		Videos:    videos,
		Complete:  complete,
	}, nil
}
