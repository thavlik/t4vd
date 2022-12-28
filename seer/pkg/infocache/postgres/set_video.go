package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *postgresInfoCache) SetVideo(video *api.VideoDetails) error {
	if _, err := c.db.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (
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
				thumbnail
			)
			VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10, $11,
				$12, $13, $14
			)
			ON CONFLICT (id) DO UPDATE
			SET (
				title,
				description,
				viewcount,
				uploader
			) = (
				EXCLUDED.title,
				EXCLUDED.description,
				EXCLUDED.viewcount,
				EXCLUDED.uploader
			)`,
			cachedVideosTable,
		),
		video.ID,
		video.Title,
		video.Description,
		video.Channel,
		video.ChannelID,
		video.Duration,
		video.ViewCount,
		video.Width,
		video.Height,
		video.FPS,
		video.UploadDate,
		video.Uploader,
		video.UploaderID,
		video.Thumbnail,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
