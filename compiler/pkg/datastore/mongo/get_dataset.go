package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
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

func (ds *mongoDataStore) getDatasetVideos(
	ctx context.Context,
	datasetID string,
) ([]*api.Video, error) {
	query, err := ds.outputVideos.Find(
		ctx,
		map[string]interface{}{
			"ds": datasetID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	var docs []map[string]interface{}
	if err := query.All(ctx, &docs); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	videos := make([]*api.Video, len(docs))
	for i, doc := range docs {
		source := doc["source"].(map[string]interface{})
		sourceID, _ := source["id"].(string)
		sourceTy, _ := source["type"].(string)
		videos[i] = &api.Video{
			ID:      doc["v"].(string),
			Details: (*api.VideoDetails)(seer.ConvertVideoDetails(doc["details"].(map[string]interface{}))),
			Source: &api.VideoSource{
				ID:          sourceID,
				Type:        sourceTy,
				SubmitterID: source["submitter"].(string),
				Submitted:   source["submitted"].(int64),
			},
		}
	}
	return videos, nil
}

func (ds *mongoDataStore) getLatestDataset(
	ctx context.Context,
	projectID string,
	complete bool,
) (*api.Dataset, error) {
	// Use the latest complete dataset
	doc := make(map[string]interface{})
	if err := ds.outputDatasets.FindOne(
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
			}),
	).Decode(&doc); err == mongo.ErrNoDocuments {
		if complete {
			// Try and get the latest incomplete dataset
			return ds.getLatestDataset(ctx, projectID, false)
		}
		return nil, datastore.ErrDatasetNotFound
	} else if err != nil {
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
	doc := make(map[string]interface{})
	if err := result.Decode(&doc); err == mongo.ErrNoDocuments {
		return nil, datastore.ErrDatasetNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	timestamp, ok := doc["t"].(int64)
	if !ok {
		return nil, errors.New("invalid timestamp")
	}
	complete, _ := doc["c"].(bool)
	videos, err := ds.getDatasetVideos(ctx, datasetID)
	if err != nil {
		return nil, errors.Wrap(err, "getDatasetVideos")
	}
	return &api.Dataset{
		ID:        datasetID,
		Timestamp: timestamp,
		Videos:    videos,
		Complete:  complete,
	}, nil
}
