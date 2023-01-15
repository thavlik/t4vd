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
		"_id":     label.ID,
		"creator": label.CreatorID,
		"created": label.Created.UnixNano(),
		"project": label.ProjectID,
	}
	if len(label.Comment) > 0 {
		doc["comment"] = label.Comment
	}
	if len(label.Tags) > 0 {
		doc["tags"] = label.Tags
	}
	if !label.Deleted.IsZero() {
		doc["deleted"] = label.Deleted.UnixNano()
	}
	if label.DeleterID != "" {
		doc["deleter"] = label.DeleterID
	}
	if label.Pad != 0 {
		doc["pad"] = label.Pad
	}
	if label.Seek != 0 {
		doc["seek"] = label.Seek
	}
	if len(label.Payload) > 0 {
		doc["payload"] = label.Payload
	}
	if label.Parent != nil {
		doc["parent"] = label.Parent.ID
	}
	if _, err := l.col.InsertOne(
		context.Background(),
		doc,
	); err != nil {
		return errors.Wrap(err, "mongo")
	}
	return nil
}
