package ytdl

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/thavlik/t4vd/base/pkg/base"
	"go.uber.org/zap"
)

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
