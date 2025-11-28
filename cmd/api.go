package main

import (
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/FelipePn10/crispaybackend/config"
	"github.com/FelipePn10/crispaybackend/internal/database"
	"github.com/FelipePn10/crispaybackend/internal/didit"
	"github.com/FelipePn10/crispaybackend/internal/email/service"
	"github.com/FelipePn10/crispaybackend/internal/handlers"
	"github.com/FelipePn10/crispaybackend/internal/repository"

	"github.com/gin-gonic/gin"
)

type application struct {
	config       *config.Config
	logger       *slog.Logger
	db           *database.DB
	emailService *service.EmailService
}

func (app *application) traceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		app.logger.Info("request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
			slog.Int("status", c.Writer.Status()),
		)
	}
}

// // Middleware to validate Didit webhook
// func (app *application) diditMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		signature := c.GetHeader("Didit-Signature")
// 		if signature == "" {
// 			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing Didit signature"})
// 			return
// 		}
// 		c.Next()
// 	}
// }

func (app *application) mount() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(app.traceMiddleware())

	// Health check global
	r.GET("/health", app.healthHandler)

	api := r.Group("/api")
	{
		app.diditRoutes(api)
	}

	return r
}

func (app *application) healthHandler(c *gin.Context) {
	resp := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "core-api",
	}
	c.JSON(http.StatusOK, resp)
}

func (app *application) diditRoutes(rg *gin.RouterGroup) {
	repo := repository.NewVerificationRepository(app.db.Queries())
	diditClient := didit.NewClient(app.config)

	webhookHandler := handlers.NewWebhookHandler(diditClient, app.config, repo, app.emailService)

	// Rotas Didit
	rg.POST("/webhooks/didit", webhookHandler.HandleVerificationWebhook)
	rg.POST("/verification/start", webhookHandler.StartVerification)
	rg.GET("/verification/status/:sessionId", webhookHandler.GetVerificationStatus)
	rg.GET("/verification/user/:userId", webhookHandler.GetUserVerifications)
}

func (app *application) run(h *gin.Engine) error {
	addr := app.config.ServerPort
	if addr == "" {
		addr = "6000"
	}
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", addr)
	return srv.ListenAndServe()
}
