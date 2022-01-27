package server

import (
	"context"
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func (s *Server) Start(port string, handler http.Handler, log logger.Logger, logMsg config.LogMsg) error {

	log.Infof(logMsg.Format, log.PackageAndFileNames(), logMsg.Init)

	s.server = &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
