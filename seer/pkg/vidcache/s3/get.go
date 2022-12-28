package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/seer/pkg/vidcache"
)

func (s *s3VidCache) Get(ctx context.Context, videoID string, w io.Writer) error {
	key := videoKey(videoID)
	sess := s3.New(base.AWSSession())
	result, err := sess.GetObjectWithContext(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return vidcache.ErrVideoNotCached
		}
		return errors.Wrap(err, "s3")
	}
	if _, err := io.Copy(w, result.Body); err != nil {
		return errors.Wrap(err, "copy")
	}
	return nil
}
