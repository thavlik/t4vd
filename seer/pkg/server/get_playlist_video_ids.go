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

func playlistTopic(id string) string {
	return "playlist:" + id
}

func (s *Server) handleGetPlaylistVideoIDs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playlistID := r.URL.Query().Get("p")
		if playlistID == "" {
			http.Error(w, "missing playlist ID", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		resp, err := s.GetPlaylistVideoIDs(
			r.Context(),
			api.GetPlaylistVideoIDsRequest{
				ID: playlistID,
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
			playlistVideoIDsChan, chanErr := s.playlistVideoIDsChan(r.Context(), playlistID)
			done := r.Context().Done()
			for {
				select {
				case <-done:
					return
				case err, ok := <-chanErr:
					if !ok || err == nil {
						return
					}
					s.log.Error("playlistVideoIDsChan", zap.Error(err))
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				case id, ok := <-playlistVideoIDsChan:
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
			s.log.Error("GetPlaylistVideoIDs", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// playlistVideoIDsChan creates a channel that is filled
// with all the known video IDs for a playlist and then
// receives new video IDs as they are discovered.
func (s *Server) playlistVideoIDsChan(
	ctx context.Context,
	playlistID string,
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
			sub, err := s.pubsub.Subscribe(ctx, playlistTopic(playlistID))
			if err != nil {
				return err
			}
			videoIDs := sub.Messages(ctx)
			// retrieve the cache. it will tell us if it is
			// complete. if it is not complete, the stream
			// should have the rest
			cachedVideoIDs, complete, err := s.cachedVideoIDs.List(ctx, playlistID)
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

func (s *Server) GetPlaylistVideoIDs(ctx context.Context, req api.GetPlaylistVideoIDsRequest) (*api.GetPlaylistVideoIDsResponse, error) {
	if req.ID == "" {
		return nil, errors.New("missing id")
	}
	log := s.log.With(zap.String("req.ID", req.ID))
	log.Debug("querying playlist videos")
	videoIDs, recency, err := s.infoCache.GetPlaylistVideoIDs(ctx, req.ID)
	if err == infocache.ErrCacheUnavailable {
		log.Debug("cached playlist info not available")
		if err := s.schedulePlaylistQuery(req.ID); err != nil {
			return nil, err
		}
		return nil, err
	} else if err != nil {
		return nil, errors.Wrap(err, "infocache.GetPlaylistVideoIDs")
	}
	log.Debug("using cached videos for playlist",
		zap.Int("numVideos", len(videoIDs)),
		zap.String("age", time.Since(recency).String()))
	if time.Since(recency) > maxRecency {
		if err := s.schedulePlaylistQuery(req.ID); err != nil {
			return nil, err
		}
	}
	return &api.GetPlaylistVideoIDsResponse{VideoIDs: videoIDs}, nil
}

func (s *Server) schedulePlaylistQuery(id string) error {
	s.log.Debug("asynchronously querying playlist details", zap.String("id", id))
	body, err := json.Marshal(&entity{
		Type: playlist,
		ID:   id,
	})
	if err != nil {
		return errors.Wrap(err, "json")
	}
	if err := s.querySched.Add(string(body)); err != nil {
		return errors.Wrap(err, "scheduler.Add")
	}
	s.log.Debug("added playlist query to scheduler", zap.String("id", id))
	return nil
}
