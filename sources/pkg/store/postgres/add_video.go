package postgres

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) AddVideo(
	projectID string,
	video *api.Video,
	blacklist bool,
	submitterID string,
) error {
	if _, err := s.db.Exec(
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
				thumbnail,
				blacklist,
				project,
				submitter
			)
			VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				$11, $12, $13, $14, $15,
				$16, $17
			)
			ON CONFLICT (id) DO UPDATE
			SET (channel, uploader, viewcount, blacklist) = (EXCLUDED.channel, EXCLUDED.uploader, EXCLUDED.viewcount, EXCLUDED.blacklist)`,
			videosTable,
		),
		store.ScopedResourceID(projectID, video.ID),
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
		video.Blacklist,
		projectID,
		submitterID,
	); err != nil {
		return errors.Wrap(err, "postgres")
	}
	return nil
}
