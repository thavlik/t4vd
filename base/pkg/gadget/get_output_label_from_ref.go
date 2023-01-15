package gadget

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/filter/pkg/api"
	"go.uber.org/zap"
)

var ErrLabelNotFound = errors.New("label not found")

func GetOutputLabelFromRef(
	ctx context.Context,
	query string,
	ref *DataRef,
	log *zap.Logger,
) (*api.Label, error) {
	gadgetName, channel, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(
		"%s/output/%s/y?%s",
		ResolveBaseURL(gadgetName),
		channel,
		query,
	)
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get input labels")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrLabelNotFound
		}
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf(
			"%s: %s: %s",
			url,
			resp.Status,
			string(body),
		)
	}
	var label api.Label
	if err := json.NewDecoder(resp.Body).Decode(&label); err != nil {
		return nil, errors.Wrap(err, "failed to decode label")
	}
	return &label, nil
}
