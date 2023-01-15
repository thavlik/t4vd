package gadget

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// ExtractBatchSize extracts the batch size from the query string.
// If the batch size is not specified, the default batch size is used.
// If the batch size is too large, an error is returned.
func ExtractBatchSize(
	query url.Values,
	maxBatchSize int,
) (int, error) {
	if v := query.Get("s"); v != "" {
		sz, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse batch size")
		}
		batchSize := int(sz)
		if batchSize > maxBatchSize {
			return 0, errors.New("batch size too large")
		}
		return batchSize, nil
	}
	return 0, errors.New("batch size not specified")
}
