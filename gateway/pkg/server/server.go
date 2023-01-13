package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pacedotdev/oto/otohttp"
	remoteiam "github.com/thavlik/t4vd/base/pkg/iam/api"
	"github.com/thavlik/t4vd/base/pkg/pubsub"
	gateway "github.com/thavlik/t4vd/gateway/pkg/api"

	"github.com/thavlik/t4vd/base/pkg/base"
	"github.com/thavlik/t4vd/base/pkg/iam"
	compiler "github.com/thavlik/t4vd/compiler/pkg/api"
	filter "github.com/thavlik/t4vd/filter/pkg/api"
	sources "github.com/thavlik/t4vd/sources/pkg/api"

	"go.uber.org/zap"
)

type websockMessageHandler func(ctx context.Context, userID string, msg map[string]interface{}, c *websocket.Conn) error

type Server struct {
	iam        iam.IAM
	seerOpts   base.ServiceOptions
	sources    sources.Sources
	compiler   compiler.Compiler
	filter     filter.Filter
	slideshow  base.ServiceOptions
	corsHeader string
	subsL      sync.Mutex
	subs       map[string][]*Subscription
	pub        pubsub.Publisher
	log        *zap.Logger
	wsHandlers map[string]websockMessageHandler
	wsSubs     map[*websocket.Conn][]*Subscription
	wsSubsL    chan struct{}
}

func NewServer(
	iam iam.IAM,
	seerOpts base.ServiceOptions,
	sources sources.Sources,
	compiler compiler.Compiler,
	filter filter.Filter,
	slideshow base.ServiceOptions,
	pub pubsub.Publisher,
	corsHeader string,
	log *zap.Logger,
) *Server {
	s := &Server{
		iam,
		seerOpts,
		sources,
		compiler,
		filter,
		slideshow,
		corsHeader,
		sync.Mutex{},
		make(map[string][]*Subscription),
		pub,
		log,
		make(map[string]websockMessageHandler),
		make(map[*websocket.Conn][]*Subscription),
		make(chan struct{}, 1),
	}
	s.setupWebSockHandlers()
	return s
}

func (s *Server) AdminListenAndServe(port int) error {
	otoServer := otohttp.NewServer()
	remoteiam.RegisterRemoteIAM(otoServer, s)
	gateway.RegisterGateway(otoServer, s)
	mux := http.NewServeMux()
	mux.Handle("/", otoServer)
	mux.HandleFunc("/readyz", base.ReadyHandler)
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}

func (s *Server) ListenAndServe(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/readyz", base.ReadyHandler)
	mux.HandleFunc("/project", s.handleGetProject())
	mux.HandleFunc("/project/create", s.handleCreateProject())
	mux.HandleFunc("/project/delete", s.handleDeleteProject())
	mux.HandleFunc("/project/list", s.handleListProjects())
	mux.HandleFunc("/project/exists", s.handleProjectExists())
	mux.HandleFunc("/project/collaborators/add", s.handleProjectAddCollaborator())
	mux.HandleFunc("/project/collaborators/remove", s.handleProjectRemoveCollaborator())
	mux.HandleFunc("/channel/add", s.handleAddChannel())
	mux.HandleFunc("/channel/remove", s.handleRemoveChannel())
	mux.HandleFunc("/channel/list", s.handleListChannels())
	mux.HandleFunc("/channel/avatar", s.handleGetChannelAvatar())
	mux.HandleFunc("/playlist/add", s.handleAddPlaylist())
	mux.HandleFunc("/playlist/remove", s.handleRemovePlaylist())
	mux.HandleFunc("/playlist/list", s.handleListPlaylists())
	mux.HandleFunc("/playlist/thumbnail", s.handleGetPlaylistThumbnail())
	mux.HandleFunc("/video/add", s.handleAddVideo())
	mux.HandleFunc("/video/remove", s.handleRemoveVideo())
	mux.HandleFunc("/video/list", s.handleListVideos())
	mux.HandleFunc("/video/thumbnail", s.handleGetVideoThumbnail())
	mux.HandleFunc("/dataset", s.handleGetDataset())
	mux.HandleFunc("/filter/stack", s.handleGetFilterStack())
	mux.HandleFunc("/filter/classify", s.handleFilterClassify())
	mux.HandleFunc("/randmarker", s.handleGetRandomMarker())
	mux.HandleFunc("/frame", s.handleGetFrame())
	mux.HandleFunc("/sse", s.handleServerSentEvents())
	mux.HandleFunc("/ws", s.handleWebSock())
	if s.iam != nil {
		mux.HandleFunc("/user/login", s.handleLogin())
		mux.HandleFunc("/user/search", s.handleUserSearch())
		mux.HandleFunc("/user/signout", s.handleSignOut())
		mux.HandleFunc("/user/register", s.handleRegister())
		mux.HandleFunc("/user/resetpassword", s.handleSetPassword())
		mux.HandleFunc("/user/exists", s.handleUserExists())
	}
	s.log.Info("listening forever", zap.Int("port", port))
	return (&http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}).ListenAndServe()
}

func addPreflightHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Access-Control-Max-Age", "1728000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "AccessToken,DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) rbacHandler(
	method string,
	permissions []string,
	f func(userID string, w http.ResponseWriter, r *http.Request) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := func() (err error) {
			w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
			if r.Method == http.MethodOptions {
				addPreflightHeaders(w)
				return nil
			}
			if method != "" && r.Method != method {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			var userID string
			if permissions != nil {
				// empty slice of permissions checks login
				// without requiring any specific permission
				userID, err = s.rbac(r.Context(), r, permissions)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					s.log.Error("auth failure",
						zap.String("r.RequestURI", r.RequestURI),
						zap.Error(err))
					return nil
				}
			}
			return f(userID, w, r)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("handler error",
				zap.String("r.RequestURI", r.RequestURI),
				zap.Error(err))
		}
	}
}

func (s *Server) handler(
	method string,
	f func(w http.ResponseWriter, r *http.Request) error,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.corsHeader)
		if r.Method == http.MethodOptions {
			addPreflightHeaders(w)
			return
		}
		if err := func() (err error) {
			if method != "" && r.Method != method {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
			return f(w, r)
		}(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error("handler error",
				zap.Error(err),
				zap.String("r.RequestURI", r.RequestURI))
		}
	}
}
