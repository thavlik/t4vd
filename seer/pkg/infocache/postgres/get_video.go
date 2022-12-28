package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/seer/pkg/api"
	"github.com/thavlik/bjjvb/seer/pkg/infocache"
)

func (c *postgresInfoCache) GetVideo(
	ctx context.Context,
	videoID string,
) (*api.VideoDetails, error) {
	row := c.db.QueryRowContext(
		ctx,
		fmt.Sprintf(
			`SELECT
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
			FROM %s WHERE id = $1`,
			cachedVideosTable,
		),
		videoID,
	)
	video := &api.VideoDetails{ID: videoID}
	if err := row.Scan(
		&video.Title,
		&video.Description,
		&video.Channel,
		&video.ChannelID,
		&video.Duration,
		&video.ViewCount,
		&video.Width,
		&video.Height,
		&video.FPS,
		&video.UploadDate,
		&video.Uploader,
		&video.UploaderID,
		&video.Thumbnail,
	); err == sql.ErrNoRows {
		return nil, infocache.ErrCacheUnavailable
	} else if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	return video, nil
}
