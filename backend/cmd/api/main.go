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
	"github.com/gowthamd/go-crud-app/internal/config"
	"github.com/gowthamd/go-crud-app/internal/db"
	"github.com/gowthamd/go-crud-app/internal/handler"
	"github.com/gowthamd/go-crud-app/internal/repository"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Database
	database, err := db.New(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// 3. Initialize Repositories & Handlers
	itemRepo := repository.NewItemRepository(database.Pool)
	itemHandler := handler.NewItemHandler(itemRepo)
	healthHandler := handler.NewHealthHandler(database)

	// 4. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS Config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", healthHandler.Liveness)
	r.Get("/health/ready", healthHandler.Readiness)

	r.Route("/api", func(r chi.Router) {
		r.Route("/items", func(r chi.Router) {
			r.Post("/", itemHandler.CreateItem)
			r.Get("/", itemHandler.ListItems)
			r.Get("/{id}", itemHandler.GetItem)
			r.Put("/{id}", itemHandler.UpdateItem)
			r.Delete("/{id}", itemHandler.DeleteItem)
		})
	})

	// 5. Start Server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Graceful Shutdown
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
