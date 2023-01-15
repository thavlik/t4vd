package server

import "github.com/pkg/errors"

func stringTagsOnly(p map[string]interface{}) error {
	for k, v := range p {
		if _, ok := v.(string); !ok {
			return errors.Errorf("invalid payload: %s is not a string", k)
		}
	}
	return nil
}
