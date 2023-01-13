package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
)

func (l *mongoLabelStore) Insert(
	label *api.Label,
) error {
	doc := map[string]interface{}{
		"_id":       label.ID,
		"submitter": label.SubmitterID,
		"submitted": label.Timestamp,
		"video":     label.Marker.VideoID,
		"timestamp": int64(label.Marker.Timestamp),
		"tags":      label.Tags,
		"project":   label.ProjectID,
	}
	if label.ParentID != "" {
		doc["parent"] = label.ParentID
	}
	if _, err := l.col.InsertOne(
		context.Background(),
		doc,
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
