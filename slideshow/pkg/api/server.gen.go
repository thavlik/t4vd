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
	slideShowGetRandomMarkerTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "slide_show_get_random_marker_total",
		Help: "Auto-generated metric incremented on every call to SlideShow.GetRandomMarker",
	})
	slideShowGetRandomMarkerSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "slide_show_get_random_marker_success_total",
		Help: "Auto-generated metric incremented on every call to SlideShow.GetRandomMarker that does not return with an error",
	})

	slideShowGetRandomStackTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "slide_show_get_random_stack_total",
		Help: "Auto-generated metric incremented on every call to SlideShow.GetRandomStack",
	})
	slideShowGetRandomStackSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "slide_show_get_random_stack_success_total",
		Help: "Auto-generated metric incremented on every call to SlideShow.GetRandomStack that does not return with an error",
	})
)

type SlideShow interface {
	GetRandomMarker(context.Context, GetRandomMarker) (*Marker, error)
	GetRandomStack(context.Context, GetRandomStack) (*Stack, error)
}

type slideShowServer struct {
	server    *otohttp.Server
	slideShow SlideShow
}

func RegisterSlideShow(server *otohttp.Server, slideShow SlideShow) {
	handler := &slideShowServer{
		server:    server,
		slideShow: slideShow,
	}
	server.Register("SlideShow", "GetRandomMarker", handler.handleGetRandomMarker)
	server.Register("SlideShow", "GetRandomStack", handler.handleGetRandomStack)
}

func (s *slideShowServer) handleGetRandomMarker(w http.ResponseWriter, r *http.Request) {
	slideShowGetRandomMarkerTotal.Inc()
	var request GetRandomMarker
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.slideShow.GetRandomMarker(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	slideShowGetRandomMarkerSuccessTotal.Inc()
}

func (s *slideShowServer) handleGetRandomStack(w http.ResponseWriter, r *http.Request) {
	slideShowGetRandomStackTotal.Inc()
	var request GetRandomStack
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.slideShow.GetRandomStack(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	slideShowGetRandomStackSuccessTotal.Inc()
}

type GetRandomMarker struct {
	ProjectID string `json:"projectID"`
}

type GetRandomStack struct {
	ProjectID string `json:"projectID"`
	Size      int    `json:"size"`
}

type Marker struct {
	VideoID string `json:"videoID"`
	Time    int64  `json:"time"`
	Error   string `json:"error,omitempty"`
}

type Stack struct {
	Markers []*Marker `json:"markers"`
	Error   string    `json:"error,omitempty"`
}

type Void struct {
}
