package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/config"
	"github.com/FallCatsinSeng/SWU_OSR/backend/internal/handler"
	mw "github.com/FallCatsinSeng/SWU_OSR/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Connect to PostgreSQL
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to PostgreSQL", zap.Error(err))
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("failed to ping PostgreSQL", zap.Error(err))
	}
	logger.Info("connected to PostgreSQL")

	// Connect to Redis
	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.Fatal("failed to parse Redis URL", zap.Error(err))
	}
	rdb := redis.NewClient(redisOpts)
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("failed to ping Redis", zap.Error(err))
	}
	logger.Info("connected to Redis")

	// Set up router
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.CORS(cfg.CORSOrigin))

	// Rate limiter
	rateLimiter := mw.NewRateLimiter(rdb, cfg.RateLimitIP, cfg.RateLimitUser)
	r.Use(rateLimiter.Middleware)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		handler.RespondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API route groups
	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Get("/feed", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
			r.Get("/profiles/{alias}", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
		})

		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/siakad-login", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
			r.Post("/github-callback", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
			r.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(mw.JWTAuth(cfg.JWTSecret))

			r.Put("/profile", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
			r.Get("/repos/available", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
			r.Post("/showcase", func(w http.ResponseWriter, r *http.Request) {
				handler.RespondJSON(w, http.StatusOK, nil)
			})
		})

		// Webhook endpoint
		r.Post("/webhooks/github", func(w http.ResponseWriter, r *http.Request) {
			handler.RespondJSON(w, http.StatusOK, nil)
		})
	})

	// Keep references to pool and rdb for future use
	_ = pool
	_ = rdb

	// Start HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("server starting", zap.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	<-done
	logger.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server shutdown failed", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}
