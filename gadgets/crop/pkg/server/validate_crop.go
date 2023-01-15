package server

import (
	"fmt"
)

var knownFields = []string{"x0", "y0", "x1", "y1", "t0", "t1"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func validateCrop(p map[string]interface{}) error {
	for k, v := range p {
		if !contains(knownFields, k) {
			return fmt.Errorf("unknown field %s", k)
		}
		if _, ok := v.(float64); !ok {
			return fmt.Errorf("field %s must be a number but is a %T", v, k)
		}
	}
	return nil
}
