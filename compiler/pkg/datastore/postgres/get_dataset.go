package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/compiler/pkg/api"
	"github.com/thavlik/t4vd/compiler/pkg/datastore"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
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
	var err error
	dataset.Videos, err = ds.getDatasetVideos(ctx, datasetID)
	if err != nil {
		return nil, errors.Wrap(err, "getDatasetVideoIDs")
	}
	return &dataset, nil
}

func (ds *postgresDataStore) getDatasetVideos(
	ctx context.Context,
	datasetID string,
) ([]*api.Video, error) {
	if datasetID == "" {
		return nil, errors.New("missing datasetID")
	}
	rows, err := ds.db.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT
				v,
				submitter,
				submitted,
				sourcety,
				sourceid,
				details
			FROM %s WHERE ds = $1`,
			outputVideosTable,
		),
		datasetID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "sql")
	}
	defer rows.Close()
	var videos []*api.Video
	for rows.Next() {
		v := &api.Video{Source: &api.VideoSource{}}
		var sourceTy, sourceID sql.NullString
		details := make(map[string]interface{})
		if err := rows.Scan(
			&v.ID,
			&v.Source.SubmitterID,
			&v.Source.Submitted,
			&sourceTy,
			&sourceID,
			&details,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		if sourceTy.Valid {
			v.Source.Type = sourceTy.String
		}
		if sourceID.Valid {
			v.Source.ID = sourceID.String
		}
		v.Details = (*api.VideoDetails)(seer.ConvertVideoDetails(details))
		videos = append(videos, v)
	}
	return videos, nil
}
