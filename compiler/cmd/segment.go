package main

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/bjjvb/base/pkg/base"
)

var segmentArgs struct {
	seer         base.ServiceOptions
	inputBucket  string
	outputBucket string
	videoID      string
	format       string
	segmentTime  time.Duration
	skipFrames   int
}

var segmentCmd = &cobra.Command{
	Use: "segment",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if segmentArgs.inputBucket == "" {
			return errors.New("missing --input-bucket")
		}
		if segmentArgs.outputBucket == "" {
			return errors.New("missing --output-bucket")
		}
		if segmentArgs.videoID == "" {
			return errors.New("missing --video-id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	base.AddServiceFlags(segmentCmd, "seer", &segmentArgs.seer, 0)
	segmentCmd.PersistentFlags().StringVar(&segmentArgs.inputBucket, "input-bucket", "", "full length video cache bucket name")
	segmentCmd.PersistentFlags().StringVar(&segmentArgs.outputBucket, "output-bucket", "", "bucket to receive the segments")
	segmentCmd.PersistentFlags().StringVar(&segmentArgs.videoID, "video-id", "", "video id")
	segmentCmd.PersistentFlags().StringVar(&segmentArgs.format, "video-format", "webm", "video format (must match what is in bucket)")
	segmentCmd.PersistentFlags().DurationVar(&segmentArgs.segmentTime, "segment-time", time.Minute, "max duration of each segment")
	segmentCmd.PersistentFlags().IntVar(&segmentArgs.skipFrames, "skip-frames", 2, "number of frames to skip during recoding")
	ConfigureCommand(segmentCmd)
}
