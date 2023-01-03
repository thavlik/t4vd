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
	compilerCompileTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compiler_compile_total",
		Help: "Auto-generated metric incremented on every call to Compiler.Compile",
	})
	compilerCompileSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compiler_compile_success_total",
		Help: "Auto-generated metric incremented on every call to Compiler.Compile that does not return with an error",
	})

	compilerGetDatasetTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compiler_get_dataset_total",
		Help: "Auto-generated metric incremented on every call to Compiler.GetDataset",
	})
	compilerGetDatasetSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compiler_get_dataset_success_total",
		Help: "Auto-generated metric incremented on every call to Compiler.GetDataset that does not return with an error",
	})

	compilerResolveProjectsForVideoTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compiler_resolve_projects_for_video_total",
		Help: "Auto-generated metric incremented on every call to Compiler.ResolveProjectsForVideo",
	})
	compilerResolveProjectsForVideoSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "compiler_resolve_projects_for_video_success_total",
		Help: "Auto-generated metric incremented on every call to Compiler.ResolveProjectsForVideo that does not return with an error",
	})
)

type Compiler interface {
	Compile(context.Context, Compile) (*Void, error)
	GetDataset(context.Context, GetDatasetRequest) (*Dataset, error)
	ResolveProjectsForVideo(context.Context, ResolveProjectsForVideoRequest) (*ResolveProjectsForVideoResponse, error)
}

type compilerServer struct {
	server   *otohttp.Server
	compiler Compiler
}

func RegisterCompiler(server *otohttp.Server, compiler Compiler) {
	handler := &compilerServer{
		server:   server,
		compiler: compiler,
	}
	server.Register("Compiler", "Compile", handler.handleCompile)
	server.Register("Compiler", "GetDataset", handler.handleGetDataset)
	server.Register("Compiler", "ResolveProjectsForVideo", handler.handleResolveProjectsForVideo)
}

func (s *compilerServer) handleCompile(w http.ResponseWriter, r *http.Request) {
	compilerCompileTotal.Inc()
	var request Compile
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.compiler.Compile(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	compilerCompileSuccessTotal.Inc()
}

func (s *compilerServer) handleGetDataset(w http.ResponseWriter, r *http.Request) {
	compilerGetDatasetTotal.Inc()
	var request GetDatasetRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.compiler.GetDataset(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	compilerGetDatasetSuccessTotal.Inc()
}

func (s *compilerServer) handleResolveProjectsForVideo(w http.ResponseWriter, r *http.Request) {
	compilerResolveProjectsForVideoTotal.Inc()
	var request ResolveProjectsForVideoRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.compiler.ResolveProjectsForVideo(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	compilerResolveProjectsForVideoSuccessTotal.Inc()
}

type Compile struct {
	ProjectID string `json:"projectID"`
	All       bool   `json:"all"`
}

type Void struct {
	Error string `json:"error,omitempty"`
}

type GetDatasetRequest struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectID"`
}

type Dataset struct {
	ID        string   `json:"id"`
	Timestamp int64    `json:"timestamp"`
	Complete  bool     `json:"complete"`
	Videos    []*Video `json:"videos"`
	Error     string   `json:"error,omitempty"`
}

type ResolveProjectsForVideoRequest struct {
	VideoID string `json:"videoID"`
}

type ResolveProjectsForVideoResponse struct {
	ProjectIDs []string `json:"projectIDs"`
	Error      string   `json:"error,omitempty"`
}

type Video struct {
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
