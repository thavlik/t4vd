// Code generated by oto; DO NOT EDIT.

package api

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/pacedotdev/oto/otohttp"
)

var (
	houndReportChannelDetailsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_channel_details_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportChannelDetails",
	})
	houndReportChannelDetailsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_channel_details_success_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportChannelDetails that does not return with an error",
	})

	houndReportChannelVideoTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_channel_video_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportChannelVideo",
	})
	houndReportChannelVideoSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_channel_video_success_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportChannelVideo that does not return with an error",
	})

	houndReportPlaylistDetailsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_playlist_details_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportPlaylistDetails",
	})
	houndReportPlaylistDetailsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_playlist_details_success_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportPlaylistDetails that does not return with an error",
	})

	houndReportPlaylistVideoTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_playlist_video_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportPlaylistVideo",
	})
	houndReportPlaylistVideoSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_playlist_video_success_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportPlaylistVideo that does not return with an error",
	})

	houndReportVideoDetailsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_video_details_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportVideoDetails",
	})
	houndReportVideoDetailsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_video_details_success_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportVideoDetails that does not return with an error",
	})

	houndReportVideoDownloadProgressTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_video_download_progress_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportVideoDownloadProgress",
	})
	houndReportVideoDownloadProgressSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "hound_report_video_download_progress_success_total",
		Help: "Auto-generated metric incremented on every call to Hound.ReportVideoDownloadProgress that does not return with an error",
	})
)

type Hound interface {
	ReportChannelDetails(context.Context, ChannelDetails) (*Void, error)
	ReportChannelVideo(context.Context, ChannelVideo) (*Void, error)
	ReportPlaylistDetails(context.Context, PlaylistDetails) (*Void, error)
	ReportPlaylistVideo(context.Context, PlaylistVideo) (*Void, error)
	ReportVideoDetails(context.Context, VideoDetails) (*Void, error)
	ReportVideoDownloadProgress(context.Context, VideoDownloadProgress) (*Void, error)
}

type houndServer struct {
	server *otohttp.Server
	hound  Hound
}

func RegisterHound(server *otohttp.Server, hound Hound) {
	handler := &houndServer{
		server: server,
		hound:  hound,
	}
	server.Register("Hound", "ReportChannelDetails", handler.handleReportChannelDetails)
	server.Register("Hound", "ReportChannelVideo", handler.handleReportChannelVideo)
	server.Register("Hound", "ReportPlaylistDetails", handler.handleReportPlaylistDetails)
	server.Register("Hound", "ReportPlaylistVideo", handler.handleReportPlaylistVideo)
	server.Register("Hound", "ReportVideoDetails", handler.handleReportVideoDetails)
	server.Register("Hound", "ReportVideoDownloadProgress", handler.handleReportVideoDownloadProgress)
}

func (s *houndServer) handleReportChannelDetails(w http.ResponseWriter, r *http.Request) {
	houndReportChannelDetailsTotal.Inc()
	var request ChannelDetails
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.hound.ReportChannelDetails(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	houndReportChannelDetailsSuccessTotal.Inc()
}

func (s *houndServer) handleReportChannelVideo(w http.ResponseWriter, r *http.Request) {
	houndReportChannelVideoTotal.Inc()
	var request ChannelVideo
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.hound.ReportChannelVideo(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	houndReportChannelVideoSuccessTotal.Inc()
}

func (s *houndServer) handleReportPlaylistDetails(w http.ResponseWriter, r *http.Request) {
	houndReportPlaylistDetailsTotal.Inc()
	var request PlaylistDetails
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.hound.ReportPlaylistDetails(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	houndReportPlaylistDetailsSuccessTotal.Inc()
}

func (s *houndServer) handleReportPlaylistVideo(w http.ResponseWriter, r *http.Request) {
	houndReportPlaylistVideoTotal.Inc()
	var request PlaylistVideo
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.hound.ReportPlaylistVideo(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	houndReportPlaylistVideoSuccessTotal.Inc()
}

func (s *houndServer) handleReportVideoDetails(w http.ResponseWriter, r *http.Request) {
	houndReportVideoDetailsTotal.Inc()
	var request VideoDetails
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.hound.ReportVideoDetails(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	houndReportVideoDetailsSuccessTotal.Inc()
}

func (s *houndServer) handleReportVideoDownloadProgress(w http.ResponseWriter, r *http.Request) {
	houndReportVideoDownloadProgressTotal.Inc()
	var request VideoDownloadProgress
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.hound.ReportVideoDownloadProgress(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	houndReportVideoDownloadProgressSuccessTotal.Inc()
}

type ChannelDetails struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Subs   string `json:"subs"`
}

type VideoDetails struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
	UploadDate  string `json:"uploadDate"`
	Uploader    string `json:"uploader"`
	UploaderID  string `json:"uploaderID"`
	Channel     string `json:"channel"`
	ChannelID   string `json:"channelID"`
	Duration    int64  `json:"duration"`
	ViewCount   int64  `json:"viewCount"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FPS         int    `json:"fPS"`
}

type ChannelVideo struct {
	ChannelID string       `json:"channelID"`
	NumVideos int          `json:"numVideos"`
	Video     VideoDetails `json:"video"`
}

type Void struct {
	Error string `json:"error,omitempty"`
}

type PlaylistDetails struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Channel   string `json:"channel"`
	ChannelID string `json:"channelID"`
	NumVideos int    `json:"numVideos"`
}

type PlaylistVideo struct {
	PlaylistID string       `json:"playlistID"`
	NumVideos  int          `json:"numVideos"`
	Video      VideoDetails `json:"video"`
}

type VideoDownloadProgress struct {
	ID      string  `json:"id"`
	Total   int64   `json:"total"`
	Rate    float64 `json:"rate"`
	Elapsed int64   `json:"elapsed"`
}