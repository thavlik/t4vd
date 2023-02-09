package ytdl

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"go.uber.org/zap"
)

var ErrAgeRestricted = errors.New("video is age restricted")

//var ytAgeRestrictedMsg = "ERROR: Sign in to confirm your age"

func Query(
	ctx context.Context,
	input string,
	videos chan<- *api.VideoDetails,
	limit int,
	log *zap.Logger,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("panic: %v", r)
		}
	}()
	command := "youtube-dl -i -j"
	if limit != 0 {
		command += fmt.Sprintf(" --max-downloads %d", limit)
	}
	command += fmt.Sprintf(` -- %s`, input)
	cmd := exec.Command("bash", "-c", command)
	r, w := io.Pipe()
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	fmt.Printf("> %s\n", command)
	start := time.Now()
	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "start")
	}
	scanner := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	stdoutDone := make(chan error, 1)
	go func() {
		defer close(stdoutDone)
		stdoutDone <- func() (err error) {
			defer close(videos)
			for scanner.Scan() {
				output := make(map[string]interface{})
				if err := json.Unmarshal(scanner.Bytes(), &output); err != nil {
					// Ignore invalid lines vomited by youtube-dl
					continue
				}
				video := api.ConvertVideoDetails(output)
				log.Debug("got video", zap.String("videoID", video.ID))
				select {
				case <-ctx.Done():
					return ctx.Err()
				case videos <- video:
				}
			}
			if err := scanner.Err(); err != nil {
				log.Error("scanner error", zap.Error(err))
				return errors.Wrap(err, "scanner")
			}
			return nil
		}()
	}()
	log.Debug("waiting on youtube-dl termination")
	done := make(chan error, 1)
	go func() {
		err := cmd.Wait()
		if err != nil && limit != 0 {
			if exiterr, ok := err.(*exec.ExitError); ok && exiterr.ExitCode() == 101 {
				// standard exit code for MaxDownloadsReached
				done <- nil
				return
			}
		}
		done <- err
	}()
	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		_ = w.Close()
		return ctx.Err()
	case err := <-done:
		_ = w.Close()
		if err != nil {
			return fmt.Errorf("failed to run '%s': %v", command, err)
		}
	}
	if err := <-stdoutDone; err != nil {
		return errors.Wrap(err, "stdout")
	}
	log.Debug("youtube-dl completed", base.Elapsed(start))
	return nil
}
