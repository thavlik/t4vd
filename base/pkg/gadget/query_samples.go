package gadget

import (
	"context"
	"encoding/json"
	"net/http"
)

func querySamples(
	ctx context.Context,
	endpoint string,
	channelName string,
	headers map[string]string,
	out interface{},
) error {
	req, err := http.NewRequest(
		http.MethodGet,
		channelEndpoint(endpoint, channelName),
		nil,
	)
	if err != nil {
		return err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req = req.WithContext(ctx)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return err
	}
	return nil
}
