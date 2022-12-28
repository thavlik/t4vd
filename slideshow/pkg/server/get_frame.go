package server

import (
	"bytes"
	"context"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/slideshow/pkg/imgcache"
	"github.com/thavlik/bjjvb/slideshow/pkg/slideshow"
	"go.uber.org/zap"
)

func getFrame(
	ctx context.Context,
	imgCache imgcache.ImgCache,
	bucket string,
	videoID string,
	t time.Duration,
	log *zap.Logger,
) (data []byte, err error) {
	log = log.With(
		zap.String("id", videoID),
		zap.Int64("t", int64(t)))
	start := time.Now()
	data, err = imgCache.GetImage(ctx, videoID, t)
	if err == nil {
		log.Debug("retrieved frame from imgcache", base.Elapsed(start))
		return data, nil
	} else if err == imgcache.ErrNotCached {
		// Get frame from bucket
		log.Debug("extracting frame from bucket")
		frame, err := slideshow.GetSingleFrameFromBucket(
			bucket,
			videoID+".webm",
			t,
		)
		if err != nil {
			return nil, errors.Wrap(err, "GetSingleFrameFromBucket")
		}
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, frame, &jpeg.Options{
			Quality: 100,
		}); err != nil {
			return nil, errors.Wrap(err, "jpeg.Encode")
		}
		data = buf.Bytes()
		// Synchrously cache the image. Used by filter microservice
		// to ensure the stack is good to go. This needs to sync
		// so the frame can be discarded if imgcache fails.
		if err := imgCache.SetImage(videoID, t, data); err != nil {
			return nil, errors.Wrap(err, "imgcache.SetImage")
		}
		log.Debug("retrieved frame from bucket",
			base.Elapsed(start))
		return data, nil
	}
	return nil, errors.Wrap(err, "imgcache.GetImage")
}

func (s *Server) handleGetFrame() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() error {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
				return errors.New("bad method")
			}
			w.Header().Set("Content-Type", "image/jpeg")
			videoID := r.URL.Query().Get("v")
			if videoID == "" {
				return errors.New("missing videoID from query")
			}
			tv, err := strconv.ParseInt(r.URL.Query().Get("t"), 10, 64)
			if err != nil {
				return errors.Wrap(err, "parse time query")
			}
			noDownload := r.URL.Query().Get("nodownload") == "1"
			t := time.Duration(tv)
			s.log.Debug("handleGetFrame",
				zap.String("videoID", videoID),
				zap.String("t", t.String()),
				zap.Bool("noDownload", noDownload))
			data, err := getFrame(
				r.Context(),
				s.imgCache,
				s.bucket,
				videoID,
				t,
				s.log,
			)
			if err != nil {
				return errors.Wrap(err, "getFrame")
			}
			if !noDownload {
				if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
					return errors.Wrap(err, "copy")
				}
			}
			return nil
		}(); err != nil {
			s.log.Error("get frame handler error", zap.Error(err))
		}
	}
}
