package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
	"github.com/thavlik/t4vd/sources/pkg/store"
)

func (s *postgresStore) ListChannels(
	ctx context.Context,
	projectID string,
) ([]*api.Channel, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT
				id,
				blacklist
			FROM %s WHERE project = $1`,
			channelsTable,
		),
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var channels []*api.Channel
	for rows.Next() {
		var id string
		var blacklist bool
		if err := rows.Scan(
			&id,
			&blacklist,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		channels = append(channels, &api.Channel{
			ID:        store.ExtractResourceID(id),
			Blacklist: blacklist,
		})
	}
	return channels, nil
}

func (s *postgresStore) ListChannelIDs(
	ctx context.Context,
	projectID string,
	blacklist bool,
) ([]string, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(
			"SELECT id FROM %s WHERE blacklist = $1 AND project = $2",
			channelsTable,
		),
		blacklist,
		projectID,
	)
	if err != nil {
		return nil, errors.Wrap(err, "postgres")
	}
	defer rows.Close()
	var channelIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		channelIDs = append(channelIDs, store.ExtractResourceID(id))
	}
	return channelIDs, nil
}
