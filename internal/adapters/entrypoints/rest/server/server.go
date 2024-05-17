package server

import (
	"context"
	"errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/registry"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/routes"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	http            http.Server
	logLevel        log.Level
	router          *mux.Router
	doneChannel     chan os.Signal
	env             environment.Environment
	useCaseRegistry registry.UseCaseRegistry
}

func NewServer(env environment.Environment, useCaseRegistry registry.UseCaseRegistry, logLevel log.Level) (*Server, chan os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return &Server{
		env:             env,
		doneChannel:     done,
		logLevel:        logLevel,
		router:          mux.NewRouter(),
		useCaseRegistry: useCaseRegistry,
	}, done
}

func (s *Server) start() error {
	routes.ConfigureRoutes(s.router, s.env, s.useCaseRegistry)
	w := log.StandardLogger().WriterLevel(s.logLevel)
	h := handlers.LoggingHandler(w, s.router)
	defer func(w *io.PipeWriter) {
		_ = w.Close()
	}(w)
	s.http = http.Server{
		Addr:              ":" + strconv.FormatUint(uint64(s.env.Port), 10),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
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
