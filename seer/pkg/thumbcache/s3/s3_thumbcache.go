package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/thumbcache"
	"go.uber.org/zap"
)

type s3ThumbCache struct {
	bucketName string
	log        *zap.Logger
}

func NewS3ThumbCache(
	bucketName string,
	log *zap.Logger,
) thumbcache.ThumbCache {
	return &s3ThumbCache{bucketName, log}
}

func thumbKey(id string) string {
	return fmt.Sprintf("%s.jpg", id)
}

func freeMultipartUploads(
	bucket, key string,
	log *zap.Logger,
) error {
	s3Client := s3.New(base.AWSSession())
	uploads, err := s3Client.ListMultipartUploads(&s3.ListMultipartUploadsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return errors.Wrap(err, "ListMultipartUploads")
	}
	for _, upload := range uploads.Uploads {
		if *upload.Key == key {
			if _, err := s3Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				Bucket:   aws.String(bucket),
				Key:      aws.String(key),
				UploadId: upload.UploadId,
			}); err != nil {
				log.Error("failed to abort multipart upload",
					zap.String("uploadId", *upload.UploadId),
					zap.Error(err))
			}
		}
	}
	return nil
}
