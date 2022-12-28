package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"go.uber.org/zap"
)

func (s *s3VidCache) List(
	ctx context.Context,
	marker string,
) (videoIDs []string, isTruncated bool, nextMarker string, err error) {
	sess := s3.New(base.AWSSession())
	var m *string
	if marker != "" {
		m = aws.String(marker)
	}
	resp, err := sess.ListObjectsWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(s.bucketName),
		Marker: m,
	})
	if err != nil {
		return nil, false, "", errors.Wrap(err, "s3")
	}
	for _, item := range resp.Contents {
		i := strings.Index(*item.Key, ".")
		if i == -1 {
			s.log.Warn("malformed s3 object key",
				zap.String("bucket", s.bucketName),
				zap.String("key", *item.Key))
			continue
		}
		videoIDs = append(videoIDs, (*item.Key)[:i])
	}
	return videoIDs, aws.BoolValue(resp.IsTruncated), aws.StringValue(resp.NextMarker), nil
}
