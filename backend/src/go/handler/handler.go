package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"health_app/api/model"
)

type Store interface {
	Ingest(metrics []model.Metric) error
	GetSummary(date string) (*model.Summary, error)
	GetVitalsHR(date string) ([]model.TimeSeriesValue, error)
	GetVitalsBP(endDate string) ([]model.BloodPressure, error)
	GetVitalsGlucose(endDate string) ([]model.Glucose, error)
	GetSleep(endDate string) ([]model.Sleep, error)
	GetWorkouts(date string) ([]model.Workout, error)
	GetDietaryTrends(endDate string) ([]model.DietaryTrend, error)
	GetDietaryMealsToday(date string) ([]model.Meal, error)
	GetBodyComposition(endDate string) ([]model.BodyComposition, error)
}

type Handler struct {
	store Store
}

func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) HandleIngest(w http.ResponseWriter, r *http.Request) {
	var req model.IngestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.Ingest(req.Metrics); err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handler) HandleGetSummary(w http.ResponseWriter, r *http.Request) {
	date := getDateQueryParam(r)
	log.Printf("Received request to comput summary data for %s", date)
	summary, err := h.store.GetSummary(date)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, summary)
}

func (h *Handler) HandleGetVitalsHR(w http.ResponseWriter, r *http.Request) {
	date := getDateQueryParam(r)
	hr, err := h.store.GetVitalsHR(date)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, hr)
}

func (h *Handler) HandleGetVitalsBP(w http.ResponseWriter, r *http.Request) {
	endDate := getEndDateQueryParam(r)
	bp, err := h.store.GetVitalsBP(endDate)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, bp)
}

func (h *Handler) HandleGetVitalsGlucose(w http.ResponseWriter, r *http.Request) {
	endDate := getEndDateQueryParam(r)
	glucose, err := h.store.GetVitalsGlucose(endDate)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, glucose)
}

func (h *Handler) HandleGetSleep(w http.ResponseWriter, r *http.Request) {
	endDate := getEndDateQueryParam(r)
	sleep, err := h.store.GetSleep(endDate)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, sleep)
}

func (h *Handler) HandleGetWorkouts(w http.ResponseWriter, r *http.Request) {
	date := getDateQueryParam(r)
	workouts, err := h.store.GetWorkouts(date)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, workouts)
}

func (h *Handler) HandleGetDietaryTrends(w http.ResponseWriter, r *http.Request) {
	endDate := getEndDateQueryParam(r)
	trends, err := h.store.GetDietaryTrends(endDate)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, trends)
}

func (h *Handler) HandleGetDietaryMealsToday(w http.ResponseWriter, r *http.Request) {
	date := getDateQueryParam(r)
	meals, err := h.store.GetDietaryMealsToday(date)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, meals)
}

func (h *Handler) HandleGetBodyComposition(w http.ResponseWriter, r *http.Request) {
	endDate := getEndDateQueryParam(r)
	bodyComp, err := h.store.GetBodyComposition(endDate)
	if err != nil {
		log.Printf("ERROR: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, bodyComp)
}

func getDateQueryParam(r *http.Request) string {
	date := r.URL.Query().Get("date")
	if date == "" {
		log.Printf("Date not found...Received url: %s", r.URL)
		date = time.Now().UTC().Format("2006-01-02")
	}
	return date
}

func getEndDateQueryParam(r *http.Request) string {
	date := r.URL.Query().Get("end_date")
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	return date
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
