package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/thavlik/bjjvb/compiler/pkg/api"
)

func (ds *postgresDataStore) SaveDataset(
	ctx context.Context,
	projectID string,
	videos []*api.Video,
	complete bool,
	timestamp time.Time,
) (*api.Dataset, error) {
	datasetID := uuid.New().String()
	if _, err := ds.db.ExecContext(
		ctx,
		fmt.Sprintf(
			`INSERT INTO %s (id, t, c, p) VALUES ($1, $2, $3, $4)`,
			outputDatasetsTable,
		),
		datasetID,
		timestamp.UnixNano(),
		complete,
		projectID,
	); err != nil {
		return nil, errors.Wrap(err, "sql")
	}
	for _, video := range videos {
		if video.ID == "" {
			return nil, errors.New("sanity check failed: video has empty id")
		}
		if _, err := ds.db.ExecContext(
			ctx,
			fmt.Sprintf(
				`INSERT INTO %s (v, ds) VALUES ($1, $2)`,
				outputVideosTable,
			),
			video.ID,
			datasetID,
		); err != nil {
			return nil, errors.Wrap(err, "sql")
		}
	}
	return &api.Dataset{
		ID:        datasetID,
		Timestamp: timestamp.UnixNano(),
		Videos:    videos,
	}, nil
}
