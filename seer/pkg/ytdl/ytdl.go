package ytdl

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
) error {
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
			if exiterr, ok := err.(*exec.ExitError); ok {
				if exiterr.ExitCode() == 101 {
					// standard exit code for MaxDownloadsReached
					done <- nil
					return
				}
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

func LogWriter(
	w io.Writer,
	log *zap.Logger,
	printBytes int,
	ctx context.Context,
	onProgress chan<- *base.DownloadProgress,
) io.Writer {
	return &logWriter{
		ctx:        ctx,
		w:          w,
		log:        log,
		start:      time.Now(),
		printBytes: printBytes,
		onProgress: onProgress,
	}
}

type logWriter struct {
	ctx        context.Context
	w          io.Writer
	log        *zap.Logger
	total      int64
	start      time.Time
	pcur       int
	printBytes int
	onProgress chan<- *base.DownloadProgress
}

func mb(b int64) string {
	v := float64(b) / (1000.0 * 1000.0)
	return fmt.Sprintf("%.2f MiB", v)
}

func (t *logWriter) Write(p []byte) (int, error) {
	n, err := t.w.Write(p)
	if err != nil {
		t.log.Error("error writing",
			zap.Error(err),
			zap.String("total", mb(t.total)))
		return 0, err
	}
	first := t.total == 0
	t.total += int64(n)
	t.pcur += n
	if first || t.pcur > t.printBytes {
		t.pcur = 0
		t.log.Debug("still downloading video",
			zap.String("total", mb(t.total)),
			base.Elapsed(t.start))
		if t.onProgress != nil {
			elapsed := time.Since(t.start)
			select {
			case <-t.ctx.Done():
				return 0, t.ctx.Err()
			case t.onProgress <- &base.DownloadProgress{
				Total:   int64(t.total),
				Rate:    float64(t.total) / float64(elapsed) / float64(time.Second), // TODO: calculate a more temporally local rate
				Elapsed: elapsed,
			}:
			}
		}
	}
	return n, nil
}

func Download(
	ctx context.Context,
	input string,
	w io.WriteCloser,
	videoFormat string,
	includeAudio bool,
	onProgress chan<- *base.DownloadProgress,
	log *zap.Logger,
) error {
	base.ProgressDownload(ctx, onProgress)
	defer base.ProgressDownload(ctx, onProgress)
	command := fmt.Sprintf(
		`youtube-dl -f "bestvideo[ext=webm]/webm" -o - -- "%s"`,
		input,
	)
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
