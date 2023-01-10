package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/seer/pkg/api"
	"github.com/thavlik/t4vd/seer/pkg/infocache"
	"go.uber.org/zap"
)

func write(w http.ResponseWriter, data []byte) error {
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "w.Write")
	}
	if _, err := w.Write([]byte("\n")); err != nil {
		return errors.Wrap(err, "w.Write")
	}
	return nil
}

func (s *Server) handleGetChannelVideoIDs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		channelID := r.URL.Query().Get("c")
		if channelID == "" {
			http.Error(w, "missing channel ID", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		resp, err := s.GetChannelVideoIDs(
			r.Context(),
			api.GetChannelVideoIDsRequest{
				ID: channelID,
			})
		if err == nil {
			// all of the video IDs are available
			for _, id := range resp.VideoIDs {
				if err := write(w, []byte(id)); err != nil {
					s.log.Error("write", zap.Error(err))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			return
		} else if err == infocache.ErrCacheUnavailable {
			// wait for the data to be available
			w.Header().Set("Connection", "keep-alive")
			// request a channel that receives all
			// video IDs as they yielded by the ytdl
			// query and then send them to the client
			channelVideoIDsChan, chanErr := s.channelVideoIDsChan(r.Context(), channelID)
			done := r.Context().Done()
			for {
				select {
				case <-done:
					return
				case err, ok := <-chanErr:
					if !ok || err == nil {
						return
					}
					s.log.Error("channelVideoIDsChan", zap.Error(err))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				case id, ok := <-channelVideoIDsChan:
					if !ok {
						return
					}
					if err := write(w, []byte(id)); err != nil {
						s.log.Error("write", zap.Error(err))
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
		} else if err != nil {
			s.log.Error("GetChannelVideoIDs", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func channelTopic(id string) string {
	return "channel:" + id
}

// channelVideoIDsChan creates a channel that is filled
// with all the known video IDs for a channel and then
// receives new video IDs as they are discovered.
func (s *Server) channelVideoIDsChan(
	ctx context.Context,
	channelID string,
) (<-chan string, <-chan error) {
	ch := make(chan string, 32)
	err := make(chan error, 1)
	go func() {
		defer close(ch)
		defer close(err)
		err <- func() error {
			// keep track of which video IDs have been sent
			sent := make(map[string]struct{})
			// first we subscribe to the topic so we don't
			// miss any video IDs that are discovered while
			// we are retrieving the cache
			sub, err := s.pubsub.Subscribe(ctx, channelTopic(channelID))
			if err != nil {
				return err
			}
			videoIDs := sub.Messages(ctx)
			// retrieve the cache. it will tell us if it is
			// complete. if it is not complete, the stream
			// should have the rest
			cachedVideoIDs, complete, err := s.cachedVideoIDs.List(ctx, channelID)
			if err != nil {
				return errors.Wrap(err, "cachedVideoIDs.List")
			}
			for _, videoID := range cachedVideoIDs {
				sent[videoID] = struct{}{}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ch <- videoID:
					continue
				}
			}
			if complete {
				// we've sent all the video IDs for the channel
				return nil
			}
			// get the rest of the video IDs from the topic
			for {
				select {
				case <-ctx.Done():
					return nil
				case videoID, ok := <-videoIDs:
					if string(videoID) == "" {
						// we've reached the end of the topic
						return nil
					} else if !ok {
						return errors.New("video IDs channel closed unexpectedly")
					} else if _, ok := sent[string(videoID)]; ok {
						// we've already sent this video ID
						continue
					}
					sent[string(videoID)] = struct{}{}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ch <- string(videoID):
						continue
					}
				}
			}
		}()
	}()
	return ch, err
}

func (s *Server) GetChannelVideoIDs(ctx context.Context, req api.GetChannelVideoIDsRequest) (*api.GetChannelVideoIDsResponse, error) {
	if req.ID == "" {
		return nil, errors.New("missing id")
	}
	log := s.log.With(zap.String("channelID", req.ID))
	log.Debug("querying channel videos")
	videoIDs, recency, err := s.infoCache.GetChannelVideoIDs(ctx, req.ID)
	if err == infocache.ErrCacheUnavailable {
		log.Debug("cached channel info not available")
		if err := s.scheduleChannelQuery(req.ID); err != nil {
			return nil, err
		}
		return nil, err
	} else if err != nil {
		return nil, errors.Wrap(err, "infocache.GetChannelVideoIDs")
	}
	log.Debug("using cached videos for channel",
		zap.Int("numVideos", len(videoIDs)),
		zap.String("age", time.Since(recency).String()))
	if time.Since(recency) > maxRecency {
		if err := s.scheduleChannelQuery(req.ID); err != nil {
			return nil, err
		}
	}
	return &api.GetChannelVideoIDsResponse{VideoIDs: videoIDs}, nil
}

func (s *Server) scheduleChannelQuery(id string) error {
	s.log.Debug("asynchronously querying channel details", zap.String("id", id))
	body, err := json.Marshal(&entity{
		Type: channel,
		ID:   id,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added channel query to scheduler", zap.String("id", id))
	return nil
}
