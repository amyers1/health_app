package main

import (
	"context"
	"fmt"
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
	"health_app/api/handler"
	"health_app/api/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	influxStore, err := store.NewInfluxDBStore()
	if err != nil {
		log.Fatalf("Failed to create InfluxDB store: %v", err)
	}

	h := handler.NewHandler(influxStore)

	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/ingest", h.HandleIngest)
		r.Get("/summary", h.HandleGetSummary)
		r.Get("/vitals/hr", h.HandleGetVitalsHR)
		r.Get("/vitals/bp", h.HandleGetVitalsBP)
		r.Get("/vitals/glucose", h.HandleGetVitalsGlucose)
		r.Get("/sleep", h.HandleGetSleep)
		r.Get("/workouts", h.HandleGetWorkouts)
		r.Get("/dietary/trends", h.HandleGetDietaryTrends)
		r.Get("/dietary/meals/today", h.HandleGetDietaryMealsToday)
		r.Get("/body/composition", h.HandleGetBodyComposition)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "13001"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server starting on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Block until we receive a signal
	<-quit
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Clean up InfluxDB connection
	log.Println("Closing InfluxDB connection...")
	influxStore.Close()

	log.Println("Server exited")
}
