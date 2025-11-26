package main

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	r.Use(app.traceMiddleware)
	r.Get("/health", app.healthHandler)
	//r.Post("/didit/section/create", app.diditSectionCreate)

	// r.Route("/v1", func(v1 chi.Router) {
	// 	v1.Mount("/users", app.userRoutes())
	// 	// v1.Mount("/billing", app.billingRoutes())
	// 	// v1.Mount("/auth", app.authRoutes())
	// })

	return r
}

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "core-api",
	}
	app.writeJSON(w, http.StatusOK, resp)
}

// func (app *application) diditSectionCreate(w http.ResponseWriter, r *http.Request) {
//
// }

// run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data map[string]any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(data); err != nil {
		http.Error(w, `{"error":"failed to write response"}`, http.StatusInternalServerError)
		return
	}
}

func (app *application) traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		app.logger.Info("request completed",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
	})
}

type application struct {
	config config
	logger *slog.Logger
}

type config struct {
	addr string
}
