package server

import (
	"context"

	compiler "github.com/thavlik/bjjvb/compiler/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) triggerRecompile(projectID string) {
	if s.compiler == nil {
		return
	}
	if _, err := s.compiler.Compile(context.Background(), compiler.Compile{
		ProjectID: projectID,
	}); err != nil {
		s.log.Error("failed to trigger recompile", zap.Error(err))
	}
}
