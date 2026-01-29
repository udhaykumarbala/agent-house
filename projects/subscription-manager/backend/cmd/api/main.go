package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/subtrack/backend/internal/config"
	"github.com/subtrack/backend/internal/database"
	"github.com/subtrack/backend/internal/handler"
	"github.com/subtrack/backend/internal/repository"
	"github.com/subtrack/backend/internal/service"
	"github.com/subtrack/backend/pkg/response"
)

func main() {
	godotenv.Load()

	cfg := config.Load()

	ctx := context.Background()
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db.Pool)
	sessionRepo := repository.NewSessionRepository(db.Pool)
	subRepo := repository.NewSubscriptionRepository(db.Pool)

	authService := service.NewAuthService(userRepo, sessionRepo)
	subService := service.NewSubscriptionService(subRepo)
	analyticsService := service.NewAnalyticsService(subRepo)

	authHandler := handler.NewAuthHandler(authService)
	subHandler := handler.NewSubscriptionHandler(subService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	authMiddleware := handler.NewAuthMiddleware(authService)

	authRateLimiter := handler.NewRateLimiter(10, time.Minute)
	apiRateLimiter := handler.NewRateLimiter(100, time.Minute)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(handler.SecurityHeaders)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           43200,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.HealthCheck(r.Context()); err != nil {
			response.Error(w, http.StatusServiceUnavailable, "UNHEALTHY", "Database connection failed")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Use(authRateLimiter.Middleware)
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)

			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.RequireAuth)
				r.Get("/me", authHandler.Me)
			})
		})

		r.Route("/subscriptions", func(r chi.Router) {
			r.Use(apiRateLimiter.Middleware)
			r.Use(authMiddleware.RequireAuth)

			r.Get("/", subHandler.List)
			r.Post("/", subHandler.Create)
			r.Get("/{id}", subHandler.Get)
			r.Put("/{id}", subHandler.Update)
			r.Delete("/{id}", subHandler.Delete)
		})

		r.Route("/analytics", func(r chi.Router) {
			r.Use(apiRateLimiter.Middleware)
			r.Use(authMiddleware.RequireAuth)

			r.Get("/summary", analyticsHandler.Summary)
			r.Get("/upcoming", analyticsHandler.Upcoming)
		})

		r.Route("/user", func(r chi.Router) {
			r.Use(apiRateLimiter.Middleware)
			r.Use(authMiddleware.RequireAuth)

			r.Put("/preferences", authHandler.UpdatePreferences)
			r.Delete("/account", authHandler.DeleteAccount)
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
