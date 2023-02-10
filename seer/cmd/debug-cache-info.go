package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/ytdl"
	"go.uber.org/zap"
)

var debugCacheInfoArgs struct {
	inputBucket  string
	outputBucket string
}

var debugCacheInfoCmd = &cobra.Command{
	Use:  "debug-cache-info",
	Long: "This command takes an input bucket, extracts all the video IDs from the key names, fetches the corresponding metadata json from youtube, and stores it in the output bucket.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		base.CheckEnv("INPUT_BUCKET", &debugCacheInfoArgs.inputBucket)
		if debugCacheInfoArgs.inputBucket == "" {
			return errors.New("--input-bucket is required")
		}
		base.CheckEnv("OUTPUT_BUCKET", &debugCacheInfoArgs.outputBucket)
		if debugCacheInfoArgs.outputBucket == "" {
			return errors.New("--output-bucket is required")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		s3Client := s3.New(base.AWSSession())
		var nextMarker *string
		for {
			resp, err := s3Client.ListObjectsWithContext(
				cmd.Context(),
				&s3.ListObjectsInput{
					Bucket: aws.String(debugCacheInfoArgs.inputBucket),
					Marker: nextMarker,
				},
			)
			if err != nil {
				return err
			}
			for _, obj := range resp.Contents {
				key := aws.StringValue(obj.Key)
				ext := filepath.Ext(key)
				id := key[:len(key)-len(ext)]
				if len(id) != 11 {
					// not a video
					continue
				}
				outKey := id + ".json"
				head, err := s3Client.HeadObjectWithContext(
					cmd.Context(),
					&s3.HeadObjectInput{
						Bucket: aws.String(debugCacheInfoArgs.outputBucket),
						Key:    aws.String(outKey),
					},
				)
				if err == nil && aws.Int64Value(head.ContentLength) > 0 {
					// info is already cached
					base.DefaultLog.Info("info already cached", zap.String("id", id))
					continue
				} else if err != nil && !strings.Contains(err.Error(), "NotFound") {
					return errors.Wrap(err, "head object")
				}
				// info is not cached
				videos := make(chan *api.VideoDetails, 1)
				ytdlDone := make(chan error, 1)
				go func() {
					ytdlDone <- ytdl.Query(
						cmd.Context(),
						id,
						videos,
						1,
						base.DefaultLog,
					)
				}()
				select {
				case <-cmd.Context().Done():
					return cmd.Context().Err()
				case video, ok := <-videos:
					if !ok {
						break
					}
					body, err := json.Marshal(video)
					if err != nil {
						return err
					}
					if _, err := s3Client.PutObjectWithContext(
						cmd.Context(),
						&s3.PutObjectInput{
							Bucket: aws.String(debugCacheInfoArgs.outputBucket),
							Key:    aws.String(outKey),
							Body:   aws.ReadSeekCloser(bytes.NewReader(body)),
						},
					); err != nil {
						return errors.Wrap(err, "put object")
					}
					base.DefaultLog.Info(
						"cached video info",
						zap.String("id", id),
						zap.String("title", video.Title),
					)
				}
				if err := <-ytdlDone; err != nil {
					base.DefaultLog.Error(
						"ytdl query failed",
						zap.String("id", id),
						zap.Error(err),
					)
					continue
				}
			}
			if !aws.BoolValue(resp.IsTruncated) {
				break
			}
			nextMarker = resp.NextMarker
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(debugCacheInfoCmd)
	debugCacheInfoCmd.Flags().StringVarP(&debugCacheInfoArgs.inputBucket, "input-bucket", "i", "", "input bucket containing the videos/audio files")
	debugCacheInfoCmd.Flags().StringVarP(&debugCacheInfoArgs.outputBucket, "output-bucket", "o", "", "output bucket containing the json info files")
}
