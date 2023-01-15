package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/thavlik/t4vd/compiler/pkg/api"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
)

func (ds *postgresDataStore) SaveDataset(
	ctx context.Context,
	projectID string,
	videos []*api.Video,
	complete bool,
	timestamp time.Time,
) (*api.Dataset, error) {
	datasetID := uuid.New().String()
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "begin tx")
	}
	for _, video := range videos {
		if video.ID == "" {
			return nil, errors.New("sanity check failed: video has empty id")
		} else if video.Details == nil {
			return nil, errors.New("sanity check failed: video has nil details")
		}
		if _, err := tx.ExecContext(
			ctx,
			fmt.Sprintf(
				`INSERT INTO %s (
					v,
					ds,
					submitter,
					submitted,
					sourcety,
					sourceid,
					details
				) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				outputVideosTable,
			),
			video.ID,
			datasetID,
			video.Source.SubmitterID,
			video.Source.Submitted,
			sql.NullString{
				Valid:  video.Source.Type != "",
				String: video.Source.Type,
			},
			sql.NullString{
				Valid:  video.Source.ID != "",
				String: video.Source.ID,
			},
			seer.FlattenVideoDetails((*seer.VideoDetails)(video.Details)),
		); err != nil {
			return nil, errors.Wrap(err, "tx exec")
		}
	}
	if _, err := tx.ExecContext(
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
		return nil, errors.Wrap(err, "tx exec")
	}
	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "commit tx")
	}
	return &api.Dataset{
		ID:        datasetID,
		Timestamp: timestamp.UnixNano(),
		Videos:    videos,
	}, nil
}
