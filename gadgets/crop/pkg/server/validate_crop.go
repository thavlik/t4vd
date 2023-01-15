package server

import (
	"fmt"

	"github.com/thavlik/t4vd/base/pkg/base"
)

var knownFields = []string{"x0", "y0", "x1", "y1", "t0", "t1"}

func validateCrop(p map[string]interface{}) error {
	for k, v := range p {
		if !base.Contains(knownFields, k) {
			return fmt.Errorf("unknown field %s", k)
		}
		if _, ok := v.(float64); !ok {
			return fmt.Errorf("field %s must be a number but is a %T", v, k)
		}
	}
	if base.MapContainsAny(p, "x0", "y0", "x1", "y1") {
		if !base.MapContainsAll(p, "x0", "y0", "x1", "y1") {
			return fmt.Errorf("x0, y0, x1, and y1 must all be present")
		}
	}
	if base.MapContainsAny(p, "t0", "t1") {
		if !base.MapContainsAll(p, "t0", "t1") {
			return fmt.Errorf("t0 and t1 must both be present")
		}
	}
	return nil
}
