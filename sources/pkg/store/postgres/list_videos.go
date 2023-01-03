package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) ListVideos(
	ctx context.Context,
	projectID string,
) ([]*api.Video, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT
			id,
			title,
			description,
			channel,
			channelid,
			duration,
			viewcount,
			width,
			height,
			fps,
			uploaddate,
			uploader,
			uploaderid,
			thumbnail,
			blacklist
		FROM %s WHERE project = $1`,
		videosTable,
	), projectID)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var videos []*api.Video
	for rows.Next() {
		var id string
		var blacklist bool
		if err := rows.Scan(
			&id,
			&blacklist,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		videos = append(videos, &api.Video{
			ID:        store.ExtractResourceID(id),
			Blacklist: blacklist,
		})
	}
	return videos, nil
}

func (s *postgresStore) ListVideoIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT id FROM %s WHERE blacklist = $1 AND project = $2",
			videosTable,
		),
		blacklist,
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var videoIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		videoIDs = append(videoIDs, store.ExtractResourceID(id))
	}
	return videoIDs, nil
}
