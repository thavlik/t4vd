package server

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
)

func keepOnly(vids []*compiler.Video, videoIDs []string) []*compiler.Video {
	var videos []*compiler.Video
	for _, video := range vids {
		if base.Contains(videoIDs, video.ID) {
			videos = append(videos, video)
		}
	}
	return videos
}

func listCachedVideoIDs(bucket string) ([]string, error) {
	sess := s3.New(base.AWSSession())
	var videoIDs []string
	b := aws.String(bucket)
	var marker *string
	for {
		resp, err := sess.ListObjects(&s3.ListObjectsInput{
			Bucket: b,
			Marker: marker,
		})
		if err != nil {
			return nil, errors.Wrap(err, "ListObjects")
		}
		for _, item := range resp.Contents {
			key := aws.StringValue(item.Key)
			i := strings.Index(key, ".webm")
			if i != -1 && aws.Int64Value(item.Size) > 0 {
				videoIDs = append(videoIDs, key[:i])
			}
		}
		if !aws.BoolValue(resp.IsTruncated) {
			break
		}
		marker = resp.NextMarker
	}
	return videoIDs, nil
}
