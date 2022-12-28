package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/thavlik/bjjvb/compiler/pkg/api"
	"github.com/thavlik/bjjvb/compiler/pkg/datastore"
)

func (ds *postgresDataStore) GetCachedVideo(
	ctx context.Context,
	id string,
) (*api.Video, error) {
	if id == "" {
		return nil, errors.New("missing id")
	}
	row := ds.db.QueryRowContext(
		ctx,
		fmt.Sprintf("SELECT title, description, thumbnail, uploaddate, uploader, uploaderid, channel, channelid, duration, viewcount, width, height, fps FROM %s WHERE id = $1", videoCacheTable),
		id,
	)
	video := api.Video{ID: id}
	if err := row.Scan(
		&video.Title,
		&video.Description,
		&video.Thumbnail,
		&video.UploadDate,
		&video.Uploader,
		&video.UploaderID,
		&video.Channel,
		&video.ChannelID,
		&video.Duration,
		&video.ViewCount,
		&video.Width,
		&video.Height,
		&video.FPS,
	); err == sql.ErrNoRows {
		return nil, datastore.ErrVideoNotCached
	} else if err != nil {
		return nil, errors.Wrap(err, "sql")
	}
	return &video, nil
}
