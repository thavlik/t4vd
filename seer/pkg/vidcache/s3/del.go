package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"go.uber.org/zap"
)

func (s *s3VidCache) Del(videoID string) error {
	key := videoKey(videoID)
	log := s.log.With(
		zap.String("bucket", s.bucketName),
		zap.String("key", key))
	if err := freeMultipartUploads(
		s.bucketName,
		key,
		log,
	); err != nil {
		log.Error("failed to free multipart upload", zap.Error(err))
	}
	if _, err := s3.New(base.AWSSession()).
		DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		}); err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
			return nil
		}
		return errors.Wrap(err, "s3")
	}
	return nil
}
