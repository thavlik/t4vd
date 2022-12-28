package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/compiler/pkg/datastore"
	"go.uber.org/zap"
)

func (ds *postgresDataStore) GetDataset(
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

func (ds *postgresDataStore) getLatestDataset(
	ctx context.Context,
	projectID string,
	complete bool,
) (*api.Dataset, error) {
	row := ds.db.QueryRowContext(
		ctx,
		fmt.Sprintf(
			"SELECT id FROM %s WHERE p = $1 AND c = $2 ORDER BY t DESC",
			outputDatasetsTable,
		),
		projectID,
		complete,
	)
	var id string
	if err := row.Scan(&id); err == sql.ErrNoRows {
		if complete {
			// default to latest incomplete dataset
			return ds.getLatestDataset(ctx, projectID, false)
		}
		return nil, datastore.ErrDatasetNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "sql")
	}
	return ds.getSpecificDataset(ctx, projectID, id)
}

func (ds *postgresDataStore) getSpecificDataset(
	ctx context.Context,
	projectID string,
	datasetID string,
) (*api.Dataset, error) {
	if datasetID == "" {
		return nil, errors.New("missing datasetID")
	}
	row := ds.db.QueryRowContext(
		ctx,
		fmt.Sprintf(
			"SELECT t, c FROM %s WHERE id = $1 AND p = $2",
			outputDatasetsTable,
		),
		datasetID,
		projectID,
	)
	dataset := api.Dataset{ID: datasetID}
	if err := row.Scan(
		&dataset.Timestamp,
		&dataset.Complete,
	); err == sql.ErrNoRows {
		return nil, datastore.ErrDatasetNotFound
	} else if err != nil {
		return nil, errors.Wrap(err, "sql")
	}
	videoIDs, err := ds.getDatasetVideoIDs(ctx, datasetID)
	if err != nil {
		return nil, errors.Wrap(err, "getDatasetVideoIDs")
	}
	ds.log.Debug("resolving videos", zap.Int("num", len(videoIDs)))
	dataset.Videos, err = datastore.ResolveVideos(
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
	return &dataset, nil
}

func (ds *postgresDataStore) getDatasetVideoIDs(
	ctx context.Context,
	datasetID string,
) ([]string, error) {
	if datasetID == "" {
		return nil, errors.New("missing datasetID")
	}
	rows, err := ds.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT id, v FROM %s WHERE ds = $1",
			outputVideosTable,
		),
		datasetID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sql")
	}
	defer rows.Close()
	var videoIDs []string
	for rows.Next() {
		var id int64
		var v string
		if err := rows.Scan(&id, &v); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		videoIDs = append(videoIDs, v)
	}
	return videoIDs, nil
}
