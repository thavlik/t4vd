package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/gadget"
	"github.com/thavlik/t4vd/base/pkg/gadget/metadata"
	seer "github.com/thavlik/t4vd/seer/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"
	"go.uber.org/zap"
)

type Server struct {
	sources      sources.Sources
	seer         base.ServiceOptions
	slideshow    base.ServiceOptions
	compiler     base.ServiceOptions
	gadgetID     string
	projectID    string
	maxBatchSize int
	log          *zap.Logger
}

func NewServer(
	sources sources.Sources,
	seer base.ServiceOptions,
	slideshow base.ServiceOptions,
	compiler base.ServiceOptions,
	gadgetID string,
	projectID string,
	maxBatchSize int,
	log *zap.Logger,
) *Server {
	return &Server{
		sources,
		seer,
		slideshow,
		compiler,
		gadgetID,
		projectID,
		maxBatchSize,
		log,
	}
}

func (s *Server) ListenAndServe(port int) error {
	seerClient := seer.NewSeerClientFromOptions(s.seer)
	meta := &metadata.Metadata{
		Name:         "video_dataset",
		MaxBatchSize: s.maxBatchSize,
		Inputs: []*metadata.Channel{{
			Name: "videos",
		}, {
			Name: "playlists",
		}, {
			Name: "channels",
		}},
		Outputs: []*metadata.Channel{{
			Name: "frames",
		}, {
			Name: "videos",
		}},
	}

	router := mux.NewRouter()
	router.HandleFunc("/metadata", gadget.HandleGetMetadata(
		meta,
		s.log,
	))

	// note: input x is handled by seer under the hood.
	// you don't PUT to the input channel x's directly
	// as the video cache is global and not specific to
	// a single gadget instance

	// we create different input channels for the different
	// types of data that comprise a video dataset
	router.HandleFunc("/input/videos/y", handleVideos(
		s.sources,
		s.projectID,
		s.log,
	))
	router.HandleFunc("/input/channels/y", handleChannels(
		s.sources,
		s.projectID,
		s.log,
	))
	router.HandleFunc("/input/playlists/y", handlePlaylists(
		s.sources,
		s.projectID,
		s.log,
	))

	// handler for retrieving a specific frame from a specific video
	router.HandleFunc("/output/frames/x", handleGetFrameData(
		s.slideshow,
		s.log,
	))
	router.HandleFunc("/output/frames/y", handleGetFrameMeta(
		seerClient,
		s.gadgetID,
		s.log,
	))

	// handler for retrieving a specific video
	router.HandleFunc("/output/videos/x", handleGetVideoData(
		"bucketName",
		s.log,
	))
	router.HandleFunc("/output/videos/y", handleGetVideoMeta(
		seerClient,
		s.gadgetID,
		s.log,
	))

	// handler for sampling random frames from random videos
	router.HandleFunc("/sample/output/frames/y", handleSampleRandomFrames(
		s.compiler,
		s.slideshow,
		s.gadgetID,
		s.projectID,
		s.maxBatchSize,
		s.log,
	))

	// handler for sampling random videos
	router.HandleFunc("/sample/output/videos/y", handleSampleRandomVideos(
		s.compiler,
		s.log,
	))

	return (&http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}).ListenAndServe()
}
