package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/isaac-olle/gas-backoffice/internal/shared/app_errors"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/command"
	"github.com/isaac-olle/gas-backoffice/internal/shared/buses/query"
	shared_logger "github.com/isaac-olle/gas-backoffice/internal/shared/logger"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/isaac-olle/gas-backoffice/internal/server/middlewares"
)

type Server struct {
	httpAddr   string
	engine     *gin.Engine
	logger     shared_logger.Logger
	shutdown   time.Duration
	commandBus command.Bus
	queryBus   query.Bus
}

func New(
	ctx context.Context,
	host string,
	port int,
	shutdownTimeout time.Duration,
	commandBus command.Bus,
	queryBus query.Bus,
	logger shared_logger.Logger,
	presenter app_errors.ErrorPresenter,
	errorManagers []app_errors.ErrorManager,
	providers ...RouteProvider,
) (context.Context, Server) {
	router := gin.New()
	router.Use(gin.RecoveryWithWriter(os.Stderr))
	router.Use(middlewares.NewGinLogger())
	router.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	srv := Server{
		engine:     router,
		httpAddr:   fmt.Sprintf("%s:%d", host, port),
		logger:     logger,
		shutdown:   shutdownTimeout,
		commandBus: commandBus,
		queryBus:   queryBus,
	}

	srv.registerRoutes(presenter, errorManagers, providers...)

	return serverContext(ctx), srv
}

func (s *Server) registerRoutes(
	presenter app_errors.ErrorPresenter,
	errorManagers []app_errors.ErrorManager,
	providers ...RouteProvider,
) {
	v1 := s.engine.Group("/v1")
	v1.Use(middlewares.RequestID(), middlewares.HttpError(presenter, errorManagers...))

	for _, p := range providers {
		p.Register(v1)
	}
}

func (s *Server) Run(ctx context.Context, logger shared_logger.Logger) error {
	logger.Info(
		"Server running on",
		shared_logger.Field{Key: "port", Value: s.httpAddr},
	)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(
				"server shut down",
				shared_logger.Field{Key: "error", Value: err},
			)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdown)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()
	return ctx
}
