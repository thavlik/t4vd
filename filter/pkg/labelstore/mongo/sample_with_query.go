package mongo

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"github.com/thavlik/t4vd/filter/pkg/labelstore"
	"go.mongodb.org/mongo-driver/bson"
)

func (l *mongoLabelStore) SampleWithQuery(
	ctx context.Context,
	input *labelstore.SampleWithQueryInput,
) ([]*api.Label, error) {
	match := map[string]interface{}{
		"project": input.ProjectID,
	}
	var query interface{} = match
	if input.All {
		if len(input.Tags) > 0 {
			match["tags"] = map[string]interface{}{
				"$all": input.Tags,
			}
		}
		if len(input.Payload) > 0 {
			match["payload"] = input.Payload
		}
	} else if len(input.Tags) > 0 || len(input.Payload) > 0 {
		var parts []interface{}
		if len(input.Tags) > 0 {
			parts = append(parts, map[string]interface{}{
				"tags": map[string]interface{}{
					"$in": input.Tags,
				},
			})
		}
		if len(input.Payload) > 0 {
			for k, v := range input.Payload {
				parts = append(parts, map[string]interface{}{
					k: v,
				})
			}
		}
		query = []interface{}{
			map[string]interface{}{
				"$and": []interface{}{
					match,
					map[string]interface{}{
						"$or": parts,
					},
				},
			},
		}
	}
	if len(input.ExcludeTags) > 0 {
		query = map[string]interface{}{
			"$and": []interface{}{
				query,
				map[string]interface{}{
					"tags": map[string]interface{}{
						"$nin": input.ExcludeTags,
					},
				}},
		}
	}
	cur, err := l.col.Aggregate(ctx, []bson.M{
		{"$match": query},
		{"$sample": bson.M{"size": input.BatchSize}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "mongo")
	}
	defer cur.Close(ctx)
	var labels []*api.Label
	for cur.Next(ctx) {
		doc := make(map[string]interface{})
		if err := cur.Decode(&doc); err != nil {
			return nil, errors.Wrap(err, "mongo")
		}
		label := api.NewLabelFromMap(doc)
		labels = append(labels, label)
	}
	return labels, nil
}
