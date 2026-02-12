package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
	defer func() {
		log.Println("Closing InfluxDB connection...")
		influxStore.Close()
	}()

	h := handler.NewHandler(influxStore)

	r := chi.NewRouter()

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
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

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
