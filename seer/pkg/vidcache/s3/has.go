package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
)

func (s *s3VidCache) Has(
	ctx context.Context,
	videoID string,
) (bool, error) {
	key := videoKey(videoID)
	sess := s3.New(base.AWSSession())
	head, err := sess.HeadObjectWithContext(
		ctx,
		&s3.HeadObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		})
	if err != nil {
		// HeadObject returns non-standard 404 https://github.com/aws/aws-sdk-go/issues/2095
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, errors.Wrap(err, "s3")
	}
	return aws.Int64Value(head.ContentLength) > 0, nil
}
