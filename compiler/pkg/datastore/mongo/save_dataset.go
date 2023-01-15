package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
)

func (ds *mongoDataStore) SaveDataset(
	ctx context.Context,
	projectID string,
	videos []*api.Video,
	complete bool,
	timestamp time.Time,
) (*api.Dataset, error) {
	id := uuid.New().String()
	docs := make([]interface{}, len(videos))
	for i, video := range videos {
		docs[i] = map[string]interface{}{
			"v":       video.ID,
			"ds":      id,
			"details": seer.FlattenVideoDetails((*seer.VideoDetails)(video.Details)),
			"source": map[string]interface{}{
				"id":        video.Source.ID,
				"type":      video.Source.Type,
				"submitter": video.Source.SubmitterID,
				"submitted": video.Source.Submitted,
			},
		}
	}
	if n, err := ds.outputVideos.InsertMany(ctx, docs); err != nil {
		return nil, errors.Wrap(err, "mongo")
	} else if got, expected := len(n.InsertedIDs), len(videos); got != expected {
		return nil, fmt.Errorf("inserted count (%d) does not equal expected value (%d)", got, expected)
	}
	if _, err := ds.outputDatasets.InsertOne(ctx, map[string]interface{}{
		"_id": id,
		"t":   timestamp.UnixNano(),
		"c":   complete,
		"p":   projectID,
	}); err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	return &api.Dataset{
		ID:        id,
		Videos:    videos,
		Timestamp: timestamp.UnixNano(),
	}, nil
}
