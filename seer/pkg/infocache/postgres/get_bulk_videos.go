package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
)

func (c *postgresInfoCache) GetBulkVideos(
	ctx context.Context,
	videoIDs []string,
) ([]*api.VideoDetails, error) {
	rows, err := c.db.QueryContext(
		ctx,
		fmt.Sprintf(
			`SELECT
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
			FROM %s WHERE id = ANY($1)`,
			cachedVideosTable,
		),
		pq.Array(videoIDs),
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var output []*api.VideoDetails
	for rows.Next() {
		video := &api.VideoDetails{}
		if err := rows.Scan(
			&video.ID,
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
		); err != nil {
			return nil, errors.Wrap(err, "postgres")
		}
		output = append(output, video)
	}
	return output, nil
}
