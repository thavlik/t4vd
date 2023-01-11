package postgres

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/sources/pkg/api"
)

func (s *postgresStore) ListChannels(
	ctx context.Context,
	projectID string,
) ([]*api.Channel, error) {
	rows, err := s.db.QueryContext(
		ctx,
		fmt.Sprintf(`
			SELECT c, blacklist
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
		channel := &api.Channel{}
		if err := rows.Scan(
			&channel.ID,
			&channel.Blacklist,
		); err != nil {
			return nil, errors.Wrap(err, "scan")
		}
		channels = append(channels, channel)
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
			"SELECT c FROM %s WHERE blacklist = $1 AND project = $2",
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
		channelIDs = append(channelIDs, id)
	}
	return channelIDs, nil
}
