package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/p1shiA/ezLink/config"
	"github.com/p1shiA/ezLink/internal/api/handlers"
	"github.com/p1shiA/ezLink/internal/core/domain"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func NewServer(
	cfg *config.Config,
	txService domain.TransactionService,
	notifier handlers.TelegramNotifier,
	botUsername string,
	logger *zap.Logger,
) *Server {
	if cfg.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	// Custom zap logger middleware for Gin
	r.Use(func(c *gin.Context) {
		c.Next()
		logger.Debug("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
		)
	})

	// Register handlers
	r.GET("/health", handlers.HealthHandler)

	callbackHandler := handlers.NewCallbackHandler(txService, notifier, botUsername, logger)
	r.GET("/payment/callback", callbackHandler.Handle)

	port := cfg.ServerConfig.Port
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger.Named("HTTPServer"),
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server", zap.String("addr", s.httpServer.Addr))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
