package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *postgresStore) ListVideos(
	ctx context.Context,
	projectID string,
) ([]*api.Video, error) {
	rows, err := s.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT v, blacklist, submitter, submitted
		FROM %s WHERE project = $1`,
		videosTable,
	), projectID)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var videos []*api.Video
	for rows.Next() {
		video := &api.Video{}
		if err := rows.Scan(
			&video.ID,
			&video.Blacklist,
			&video.SubmitterID,
			&video.Submitted,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		videos = append(videos, video)
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
			"SELECT v FROM %s WHERE blacklist = $1 AND project = $2",
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
		videoIDs = append(videoIDs, id)
	}
	return videoIDs, nil
}
