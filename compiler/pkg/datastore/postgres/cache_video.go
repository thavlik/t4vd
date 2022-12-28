package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/thavlik/t4vd/compiler/pkg/api"
)

func (ds *postgresDataStore) CacheVideo(
	ctx context.Context,
	video *api.Video,
) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (id, title, description, thumbnail, uploaddate, uploader, uploaderid, channel, channelid, duration, viewcount, width, height, fps)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT (id) DO UPDATE
SET (title, description, channel, uploader, viewcount) = (EXCLUDED.title, EXCLUDED.description, EXCLUDED.channel, EXCLUDED.uploader, EXCLUDED.viewcount)`,
		videoCacheTable,
	)
	if _, err := ds.db.ExecContext(
		ctx,
		query,
		video.ID,
		video.Title,
		video.Description,
		video.Thumbnail,
		video.UploadDate,
		video.Uploader,
		video.UploaderID,
		video.Channel,
		video.ChannelID,
		video.Duration,
		video.ViewCount,
		video.Width,
		video.Height,
		video.FPS,
	); err != nil {
		return errors.Wrap(err, "sql")
	}
	return nil
}

func (ds *postgresDataStore) CacheBulkVideos(
	ctx context.Context,
	videos []*api.Video,
) error {
	for _, video := range videos {
		if err := ds.CacheVideo(ctx, video); err != nil {
			return err
		}
	}
	return nil
}
