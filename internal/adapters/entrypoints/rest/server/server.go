package server

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/middlewares"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type Server struct {
	http            http.Server
	logLevel        log.Level
	doneChannel     chan os.Signal
	env             environment.Environment
	useCaseRegistry registry.UseCaseRegistry
	timeouts        environment.ApplicationTimeouts
}

func NewServer(
	env environment.Environment,
	useCaseRegistry registry.UseCaseRegistry,
	logLevel log.Level,
	timeouts environment.ApplicationTimeouts,
) (*Server, chan os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &Server{
		env:             env,
		doneChannel:     done,
		logLevel:        logLevel,
		useCaseRegistry: useCaseRegistry,
		timeouts:        timeouts,
	}, done
}

func (s *Server) start() error {
	router := routes.NewRouter()
	router.Use(middlewares.NewAccessLogMiddleware(log.StandardLogger(), s.logLevel))
	routes.ConfigureRoutes(router, s.env, s.useCaseRegistry, routes.NewEndpointFactory())
	s.http = http.Server{
		Addr:              ":" + strconv.FormatUint(uint64(s.env.Port), 10),
		Handler:           router.BuildHandler(),
		ReadHeaderTimeout: s.timeouts.ServerReadHeader.Seconds(),
		WriteTimeout:      s.timeouts.ServerWrite.Seconds(),
		IdleTimeout:       s.timeouts.ServerIdle.Seconds(),
	}

	log.Info("Server started at localhost:", s.http.Addr)
	return s.http.ListenAndServe()
}

// Start to be called inside goroutine
func (s *Server) Start() {
	if err := s.start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Error running server: ", err)
		s.doneChannel <- syscall.SIGTERM
	}
}

func (s *Server) Shutdown(closeChannel chan<- bool) {
	err := s.http.Shutdown(context.Background())
	if err != nil {
		log.Error("Error shutting down server", err)
	}
	closeChannel <- true
	log.Debug("Server shutdown completed")
}
