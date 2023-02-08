package ytdl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

type Options struct {
	VideoFormat       string
	AudioFormat       string
	AudioSampleRate   int
	AudioChannelCount int
	SkipAudio         bool
	SkipVideo         bool
}

func Download(
	ctx context.Context,
	input string,
	w io.WriteCloser,
	options *Options,
	onProgress chan<- *base.DownloadProgress,
	log *zap.Logger,
) error {
	base.ProgressDownload(ctx, onProgress)
	defer base.ProgressDownload(ctx, onProgress)
	var command string
	if options.SkipAudio {
		if options.SkipVideo {
			return errors.New("both audio and video are skipped")
		}
		// video only
		command = fmt.Sprintf(
			`youtube-dl -f "bestvideo[ext=%s]/%s" -o - -- "%s"`,
			options.VideoFormat,
			options.VideoFormat,
			input,
		)
	} else if options.SkipVideo {
		// audio only
		sr := options.AudioSampleRate
		if sr == 0 {
			sr = 44100 // default
		}
		command = fmt.Sprintf(
			`youtube-dl -o - -f bestaudio --extract-audio --audio-format %s --audio-quality 0 --postprocessor-args "-osr %d -ac %d" -- %s`,
			options.AudioFormat,
			sr,
			options.AudioChannelCount,
			input,
		)
	} else {
		// audio and video
		command = fmt.Sprintf(
			`youtube-dl --merge-output-format %s -f "bestvideo[ext=%s]+bestaudio[ext=%s]/bestvideo+bestaudio" -o - -- "%s"`,
			options.VideoFormat,
			options.VideoFormat,
			options.AudioFormat,
			input,
		)
	}
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = LogWriter(w, log, 1000*1000, ctx, onProgress)
	pr, pw := io.Pipe()
	cmd.Stderr = pw
	fmt.Printf("> %s\n", command)
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "start")
	}
	stderrDone := make(chan error)
	scanner := bufio.NewScanner(pr)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	go func() {
		stderrDone <- func() (err error) {
			for scanner.Scan() {
				text := scanner.Text()
				fmt.Println(text)
				if strings.TrimSpace(text) == "ERROR: unable to download video data: HTTP Error 403: Forbidden" {
					// Retry?
					return fmt.Errorf("youtube returned error 403")
				}
			}
			if err := scanner.Err(); err == io.ErrClosedPipe {
				// The pipe was closed because the command exited
				return nil
			} else if err != nil {
				// ERROR: unable to download video data: HTTP Error 403: Forbidden
				fmt.Println("TODO: detect age restriction")
				log.Error("stderr scanner error", zap.Error(err))
				return errors.Wrap(err, "scanner")
			}
			return nil
		}()
	}()
	log.Debug("waiting on youtube-dl termination")
	exited := make(chan error)
	go func() { exited <- cmd.Wait() }()
	var err error
	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		_ = pr.Close()
		_ = pw.Close()
		_ = w.Close()
		return ctx.Err()
	case err = <-exited:
		_ = pr.Close()
		_ = pw.Close()
		_ = w.Close()
		if err != nil {
			return fmt.Errorf("failed to run '%s': %v", command, err)
		}
	}
	if err := <-stderrDone; err != nil {
		return errors.Wrap(err, "stderr parser")
	}
	return nil
}
