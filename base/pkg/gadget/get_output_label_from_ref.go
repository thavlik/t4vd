package gadget

import (
	"bytes"
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
	label *api.Label,
	ref *DataRef,
	log *zap.Logger,
) (*api.Label, error) {
	gadgetName, channel, err := ref.Get(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: ensure gadgetName points to gadgetID
	body, err := json.Marshal(label)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode label")
	}
	url := fmt.Sprintf(
		"%s/output/%s/y",
		ResolveBaseURL(gadgetName),
		channel,
	)
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body),
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
	var output api.Label
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, errors.Wrap(err, "failed to decode label")
	}
	return &output, nil
}
