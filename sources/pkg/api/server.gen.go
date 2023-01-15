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
	sourcesAddChannelTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_add_channel_total",
		Help: "Auto-generated metric incremented on every call to Sources.AddChannel",
	})
	sourcesAddChannelSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_add_channel_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.AddChannel that does not return with an error",
	})

	sourcesAddPlaylistTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_add_playlist_total",
		Help: "Auto-generated metric incremented on every call to Sources.AddPlaylist",
	})
	sourcesAddPlaylistSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_add_playlist_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.AddPlaylist that does not return with an error",
	})

	sourcesAddVideoTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_add_video_total",
		Help: "Auto-generated metric incremented on every call to Sources.AddVideo",
	})
	sourcesAddVideoSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_add_video_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.AddVideo that does not return with an error",
	})

	sourcesCreateProjectTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_create_project_total",
		Help: "Auto-generated metric incremented on every call to Sources.CreateProject",
	})
	sourcesCreateProjectSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_create_project_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.CreateProject that does not return with an error",
	})

	sourcesDeleteProjectTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_delete_project_total",
		Help: "Auto-generated metric incremented on every call to Sources.DeleteProject",
	})
	sourcesDeleteProjectSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_delete_project_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.DeleteProject that does not return with an error",
	})

	sourcesGetProjectTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProject",
	})
	sourcesGetProjectSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProject that does not return with an error",
	})

	sourcesGetProjectByNameTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_by_name_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectByName",
	})
	sourcesGetProjectByNameSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_by_name_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectByName that does not return with an error",
	})

	sourcesGetProjectIDsForChannelTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_ids_for_channel_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectIDsForChannel",
	})
	sourcesGetProjectIDsForChannelSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_ids_for_channel_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectIDsForChannel that does not return with an error",
	})

	sourcesGetProjectIDsForPlaylistTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_ids_for_playlist_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectIDsForPlaylist",
	})
	sourcesGetProjectIDsForPlaylistSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_ids_for_playlist_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectIDsForPlaylist that does not return with an error",
	})

	sourcesGetProjectIDsForVideoTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_ids_for_video_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectIDsForVideo",
	})
	sourcesGetProjectIDsForVideoSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_get_project_ids_for_video_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.GetProjectIDsForVideo that does not return with an error",
	})

	sourcesListChannelIDsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_channel_ids_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListChannelIDs",
	})
	sourcesListChannelIDsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_channel_ids_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListChannelIDs that does not return with an error",
	})

	sourcesListChannelsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_channels_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListChannels",
	})
	sourcesListChannelsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_channels_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListChannels that does not return with an error",
	})

	sourcesListPlaylistIDsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_playlist_ids_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListPlaylistIDs",
	})
	sourcesListPlaylistIDsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_playlist_ids_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListPlaylistIDs that does not return with an error",
	})

	sourcesListPlaylistsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_playlists_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListPlaylists",
	})
	sourcesListPlaylistsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_playlists_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListPlaylists that does not return with an error",
	})

	sourcesListProjectsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_projects_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListProjects",
	})
	sourcesListProjectsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_projects_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListProjects that does not return with an error",
	})

	sourcesListVideoIDsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_video_ids_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListVideoIDs",
	})
	sourcesListVideoIDsSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_video_ids_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListVideoIDs that does not return with an error",
	})

	sourcesListVideosTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_videos_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListVideos",
	})
	sourcesListVideosSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_list_videos_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.ListVideos that does not return with an error",
	})

	sourcesRemoveChannelTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_remove_channel_total",
		Help: "Auto-generated metric incremented on every call to Sources.RemoveChannel",
	})
	sourcesRemoveChannelSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_remove_channel_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.RemoveChannel that does not return with an error",
	})

	sourcesRemovePlaylistTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_remove_playlist_total",
		Help: "Auto-generated metric incremented on every call to Sources.RemovePlaylist",
	})
	sourcesRemovePlaylistSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_remove_playlist_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.RemovePlaylist that does not return with an error",
	})

	sourcesRemoveVideoTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_remove_video_total",
		Help: "Auto-generated metric incremented on every call to Sources.RemoveVideo",
	})
	sourcesRemoveVideoSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sources_remove_video_success_total",
		Help: "Auto-generated metric incremented on every call to Sources.RemoveVideo that does not return with an error",
	})
)

type Sources interface {
	AddChannel(context.Context, AddChannelRequest) (*Channel, error)
	AddPlaylist(context.Context, AddPlaylistRequest) (*Playlist, error)
	AddVideo(context.Context, AddVideoRequest) (*Video, error)
	CreateProject(context.Context, Project) (*Project, error)
	DeleteProject(context.Context, DeleteProject) (*Void, error)
	GetProject(context.Context, GetProject) (*Project, error)
	GetProjectByName(context.Context, GetProjectByName) (*Project, error)
	GetProjectIDsForChannel(context.Context, GetProjectIDsForChannelRequest) (*GetProjectIDsForChannelResponse, error)
	GetProjectIDsForPlaylist(context.Context, GetProjectIDsForPlaylistRequest) (*GetProjectIDsForPlaylistResponse, error)
	GetProjectIDsForVideo(context.Context, GetProjectIDsForVideoRequest) (*GetProjectIDsForVideoResponse, error)
	ListChannelIDs(context.Context, ListChannelIDsRequest) (*ListChannelIDsResponse, error)
	ListChannels(context.Context, ListChannelsRequest) (*ListChannelsResponse, error)
	ListPlaylistIDs(context.Context, ListPlaylistIDsRequest) (*ListPlaylistIDsResponse, error)
	ListPlaylists(context.Context, ListPlaylistsRequest) (*ListPlaylistsResponse, error)
	ListProjects(context.Context, ListProjectsRequest) (*ListProjectsResponse, error)
	ListVideoIDs(context.Context, ListVideoIDsRequest) (*ListVideoIDsResponse, error)
	ListVideos(context.Context, ListVideosRequest) (*ListVideosResponse, error)
	RemoveChannel(context.Context, RemoveChannelRequest) (*Void, error)
	RemovePlaylist(context.Context, RemovePlaylistRequest) (*Void, error)
	RemoveVideo(context.Context, RemoveVideoRequest) (*Void, error)
}

type sourcesServer struct {
	server  *otohttp.Server
	sources Sources
}

func RegisterSources(server *otohttp.Server, sources Sources) {
	handler := &sourcesServer{
		server:  server,
		sources: sources,
	}
	server.Register("Sources", "AddChannel", handler.handleAddChannel)
	server.Register("Sources", "AddPlaylist", handler.handleAddPlaylist)
	server.Register("Sources", "AddVideo", handler.handleAddVideo)
	server.Register("Sources", "CreateProject", handler.handleCreateProject)
	server.Register("Sources", "DeleteProject", handler.handleDeleteProject)
	server.Register("Sources", "GetProject", handler.handleGetProject)
	server.Register("Sources", "GetProjectByName", handler.handleGetProjectByName)
	server.Register("Sources", "GetProjectIDsForChannel", handler.handleGetProjectIDsForChannel)
	server.Register("Sources", "GetProjectIDsForPlaylist", handler.handleGetProjectIDsForPlaylist)
	server.Register("Sources", "GetProjectIDsForVideo", handler.handleGetProjectIDsForVideo)
	server.Register("Sources", "ListChannelIDs", handler.handleListChannelIDs)
	server.Register("Sources", "ListChannels", handler.handleListChannels)
	server.Register("Sources", "ListPlaylistIDs", handler.handleListPlaylistIDs)
	server.Register("Sources", "ListPlaylists", handler.handleListPlaylists)
	server.Register("Sources", "ListProjects", handler.handleListProjects)
	server.Register("Sources", "ListVideoIDs", handler.handleListVideoIDs)
	server.Register("Sources", "ListVideos", handler.handleListVideos)
	server.Register("Sources", "RemoveChannel", handler.handleRemoveChannel)
	server.Register("Sources", "RemovePlaylist", handler.handleRemovePlaylist)
	server.Register("Sources", "RemoveVideo", handler.handleRemoveVideo)
}

func (s *sourcesServer) handleAddChannel(w http.ResponseWriter, r *http.Request) {
	sourcesAddChannelTotal.Inc()
	var request AddChannelRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.AddChannel(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesAddChannelSuccessTotal.Inc()
}

func (s *sourcesServer) handleAddPlaylist(w http.ResponseWriter, r *http.Request) {
	sourcesAddPlaylistTotal.Inc()
	var request AddPlaylistRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.AddPlaylist(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesAddPlaylistSuccessTotal.Inc()
}

func (s *sourcesServer) handleAddVideo(w http.ResponseWriter, r *http.Request) {
	sourcesAddVideoTotal.Inc()
	var request AddVideoRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.AddVideo(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesAddVideoSuccessTotal.Inc()
}

func (s *sourcesServer) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	sourcesCreateProjectTotal.Inc()
	var request Project
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.CreateProject(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesCreateProjectSuccessTotal.Inc()
}

func (s *sourcesServer) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	sourcesDeleteProjectTotal.Inc()
	var request DeleteProject
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.DeleteProject(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesDeleteProjectSuccessTotal.Inc()
}

func (s *sourcesServer) handleGetProject(w http.ResponseWriter, r *http.Request) {
	sourcesGetProjectTotal.Inc()
	var request GetProject
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.GetProject(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesGetProjectSuccessTotal.Inc()
}

func (s *sourcesServer) handleGetProjectByName(w http.ResponseWriter, r *http.Request) {
	sourcesGetProjectByNameTotal.Inc()
	var request GetProjectByName
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.GetProjectByName(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesGetProjectByNameSuccessTotal.Inc()
}

func (s *sourcesServer) handleGetProjectIDsForChannel(w http.ResponseWriter, r *http.Request) {
	sourcesGetProjectIDsForChannelTotal.Inc()
	var request GetProjectIDsForChannelRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.GetProjectIDsForChannel(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesGetProjectIDsForChannelSuccessTotal.Inc()
}

func (s *sourcesServer) handleGetProjectIDsForPlaylist(w http.ResponseWriter, r *http.Request) {
	sourcesGetProjectIDsForPlaylistTotal.Inc()
	var request GetProjectIDsForPlaylistRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.GetProjectIDsForPlaylist(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesGetProjectIDsForPlaylistSuccessTotal.Inc()
}

func (s *sourcesServer) handleGetProjectIDsForVideo(w http.ResponseWriter, r *http.Request) {
	sourcesGetProjectIDsForVideoTotal.Inc()
	var request GetProjectIDsForVideoRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.GetProjectIDsForVideo(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesGetProjectIDsForVideoSuccessTotal.Inc()
}

func (s *sourcesServer) handleListChannelIDs(w http.ResponseWriter, r *http.Request) {
	sourcesListChannelIDsTotal.Inc()
	var request ListChannelIDsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListChannelIDs(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListChannelIDsSuccessTotal.Inc()
}

func (s *sourcesServer) handleListChannels(w http.ResponseWriter, r *http.Request) {
	sourcesListChannelsTotal.Inc()
	var request ListChannelsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListChannels(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListChannelsSuccessTotal.Inc()
}

func (s *sourcesServer) handleListPlaylistIDs(w http.ResponseWriter, r *http.Request) {
	sourcesListPlaylistIDsTotal.Inc()
	var request ListPlaylistIDsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListPlaylistIDs(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListPlaylistIDsSuccessTotal.Inc()
}

func (s *sourcesServer) handleListPlaylists(w http.ResponseWriter, r *http.Request) {
	sourcesListPlaylistsTotal.Inc()
	var request ListPlaylistsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListPlaylists(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListPlaylistsSuccessTotal.Inc()
}

func (s *sourcesServer) handleListProjects(w http.ResponseWriter, r *http.Request) {
	sourcesListProjectsTotal.Inc()
	var request ListProjectsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListProjects(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListProjectsSuccessTotal.Inc()
}

func (s *sourcesServer) handleListVideoIDs(w http.ResponseWriter, r *http.Request) {
	sourcesListVideoIDsTotal.Inc()
	var request ListVideoIDsRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListVideoIDs(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListVideoIDsSuccessTotal.Inc()
}

func (s *sourcesServer) handleListVideos(w http.ResponseWriter, r *http.Request) {
	sourcesListVideosTotal.Inc()
	var request ListVideosRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.ListVideos(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesListVideosSuccessTotal.Inc()
}

func (s *sourcesServer) handleRemoveChannel(w http.ResponseWriter, r *http.Request) {
	sourcesRemoveChannelTotal.Inc()
	var request RemoveChannelRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.RemoveChannel(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesRemoveChannelSuccessTotal.Inc()
}

func (s *sourcesServer) handleRemovePlaylist(w http.ResponseWriter, r *http.Request) {
	sourcesRemovePlaylistTotal.Inc()
	var request RemovePlaylistRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.RemovePlaylist(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesRemovePlaylistSuccessTotal.Inc()
}

func (s *sourcesServer) handleRemoveVideo(w http.ResponseWriter, r *http.Request) {
	sourcesRemoveVideoTotal.Inc()
	var request RemoveVideoRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.sources.RemoveVideo(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	sourcesRemoveVideoSuccessTotal.Inc()
}

type AddChannelRequest struct {
	ProjectID   string `json:"projectID"`
	Input       string `json:"input"`
	Blacklist   bool   `json:"blacklist"`
	SubmitterID string `json:"submitterID"`
}

type AddPlaylistRequest struct {
	ProjectID   string `json:"projectID"`
	Input       string `json:"input"`
	Blacklist   bool   `json:"blacklist"`
	SubmitterID string `json:"submitterID"`
}

type AddVideoRequest struct {
	ProjectID   string `json:"projectID"`
	Input       string `json:"input"`
	Blacklist   bool   `json:"blacklist"`
	SubmitterID string `json:"submitterID"`
}

type Channel struct {
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
	Error     string `json:"error,omitempty"`
}

type Collaborator struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type DeleteProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetProject struct {
	ID string `json:"id"`
}

type GetProjectByName struct {
	Name string `json:"name"`
}

type GetProjectIDsForChannelRequest struct {
	ChannelID string `json:"channelID"`
}

type GetProjectIDsForChannelResponse struct {
	ProjectIDs []string `json:"projectIDs"`
	Error      string   `json:"error,omitempty"`
}

type GetProjectIDsForPlaylistRequest struct {
	PlaylistID string `json:"playlistID"`
}

type GetProjectIDsForPlaylistResponse struct {
	ProjectIDs []string `json:"projectIDs"`
	Error      string   `json:"error,omitempty"`
}

type GetProjectIDsForVideoRequest struct {
	VideoID string `json:"videoID"`
}

type GetProjectIDsForVideoResponse struct {
	ProjectIDs []string `json:"projectIDs"`
	Error      string   `json:"error,omitempty"`
}

type ListChannelIDsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type ListChannelIDsResponse struct {
	IDs   []string `json:"iDS"`
	Error string   `json:"error,omitempty"`
}

type ListChannelsRequest struct {
	ProjectID string `json:"projectID"`
}

type ListChannelsResponse struct {
	Channels []*Channel `json:"channels"`
	Error    string     `json:"error,omitempty"`
}

type ListPlaylistIDsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type ListPlaylistIDsResponse struct {
	IDs   []string `json:"iDS"`
	Error string   `json:"error,omitempty"`
}

type ListPlaylistsRequest struct {
	ProjectID string `json:"projectID"`
}

type ListPlaylistsResponse struct {
	Playlists []*Playlist `json:"playlists"`
	Error     string      `json:"error,omitempty"`
}

type ListProjectsRequest struct {
	CreatedByUserID string `json:"createdByUserID"`
	VisibleToUserID string `json:"visibleToUserID"`
}

type ListProjectsResponse struct {
	Projects []*Project `json:"projects"`
	Error    string     `json:"error,omitempty"`
}

type ListVideoIDsRequest struct {
	ProjectID string `json:"projectID"`
	Blacklist bool   `json:"blacklist"`
}

type ListVideoIDsResponse struct {
	IDs   []string `json:"iDS"`
	Error string   `json:"error,omitempty"`
}

type ListVideosRequest struct {
	ProjectID string `json:"projectID"`
}

type ListVideosResponse struct {
	Videos []*Video `json:"videos"`
	Error  string   `json:"error,omitempty"`
}

type Playlist struct {
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
	Error     string `json:"error,omitempty"`
}

type Project struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	CreatorID   string   `json:"creatorID"`
	GroupID     string   `json:"groupID"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Error       string   `json:"error,omitempty"`
}

type RemoveChannelRequest struct {
	ProjectID string `json:"projectID"`
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
}

type RemovePlaylistRequest struct {
	ProjectID string `json:"projectID"`
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
}

type RemoveVideoRequest struct {
	ProjectID string `json:"projectID"`
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
}

type Video struct {
	ID        string `json:"id"`
	Blacklist bool   `json:"blacklist"`
	Error     string `json:"error,omitempty"`
}

type Void struct {
	Error string `json:"error,omitempty"`
}
