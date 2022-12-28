package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/slideshow/pkg/slideshow"
)

var defaultReadAhead = 1024 * 10

var getFrameArgs struct {
	bucket    string
	videoID   string
	time      time.Duration
	readAhead int
	out       string
}

var getFrameCmd = &cobra.Command{
	Use: "get-frame",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if getFrameArgs.videoID == "" {
			return errors.New("missing --video-id")
		}
		ext := filepath.Ext(getFrameArgs.out)
		if ext != "jpg" && ext != "jpeg" && ext != "png" {
			return fmt.Errorf("unsupported output format '%s'", ext)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		var r io.ReadSeeker
		s3r, err := slideshow.NewS3ReadSeeker(
			getFrameArgs.bucket,
			getFrameArgs.videoID,
		)
		if err != nil {
			return errors.Wrap(err, "NewS3ReadSeeker")
		}
		if getFrameArgs.readAhead == 0 {
			r = s3r
		} else {
			r = slideshow.NewReadAhead(s3r, getFrameArgs.readAhead)
		}
		frame, err := slideshow.GetSingleFrameFromReader(r, getFrameArgs.time)
		if err != nil {
			return errors.Wrap(err, "GetSingleFrameFromFile")
		}
		if err := func() error {
			f, err := os.OpenFile(
				getFrameArgs.out,
				os.O_CREATE|os.O_WRONLY,
				0644,
			)
			if err != nil {
				return errors.Wrap(err, "OpenFile")
			}
			defer f.Close()
			ext := filepath.Ext(getFrameArgs.out)
			if ext == "jpg" || ext == "jpeg" {
				if err := jpeg.Encode(f, frame, &jpeg.Options{
					Quality: 100,
				}); err != nil {
					return errors.Wrap(err, "jpeg")
				}
			} else if ext == "png" {
				if err := png.Encode(f, frame); err != nil {
					return errors.Wrap(err, "png")
				}
			} else {
				panic(base.Unreachable)
			}
			return nil
		}(); err != nil {
			return err
		}
		fmt.Println(time.Since(start).String())
		return nil
	},
}

func init() {
	ConfigureCommand(getFrameCmd)
	getFrameCmd.PersistentFlags().IntVarP(&getFrameArgs.readAhead, "read-ahead", "r", defaultReadAhead, "minimum number of bytes to fetch from s3 at any time")
	getFrameCmd.PersistentFlags().StringVarP(&getFrameArgs.videoID, "video-id", "v", "", "youtube video ID")
	getFrameCmd.PersistentFlags().StringVar(&getFrameArgs.bucket, "bucket", "ytvids", "s3 bucket containing full length webm videos")
	getFrameCmd.PersistentFlags().DurationVarP(&getFrameArgs.time, "time", "t", 0, "seek time")
	getFrameCmd.PersistentFlags().StringVarP(&getFrameArgs.out, "out", "o", "out.jpg", "out path")
}
