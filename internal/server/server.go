package server

import (
	"github.com/Aserose/ArchaicReverie/internal/config"
	"github.com/Aserose/ArchaicReverie/pkg/logger"
	"net/http"
)

const ServerPackageName = "server"

type Server struct {
	server *http.Server
}

func (s Server) Start(port string, handler http.Handler, log logger.Logger, logMsg config.LogMsg) error {

	log.Infof(logMsg.Format, ServerPackageName, logMsg.Init)

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return s.server.ListenAndServe()
}
