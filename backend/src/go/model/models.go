package model

import "time"

// IngestRequest is the structure for the /api/v1/ingest endpoint
type IngestRequest struct {
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Summary is the structure for the /api/v1/summary endpoint
type Summary struct {
	Steps           int     `json:"steps"`
	Distance        float64 `json:"distance"`
	ActiveCalories  float64 `json:"activeCalories"`
	BasalCalories   float64 `json:"basalCalories"`
	DietaryCalories float64 `json:"dietaryCalories"`
}

// TimeSeriesValue is a generic struct for time series data
type TimeSeriesValue struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

// BloodPressure is the structure for blood pressure data
type BloodPressure struct {
	Time      string `json:"time"`
	Systolic  int    `json:"systolic"`
	Diastolic int    `json:"diastolic"`
	Category  string `json:"category"`
}

// Glucose is the structure for glucose data
type Glucose struct {
	Time  string  `json:"time"`
	Value float64 `json:"value"`
}

// Sleep is the structure for sleep data
type Sleep struct {
	Date            string  `json:"date"`
	TotalDuration   float64 `json:"totalDuration"`
	DeepSleep       float64 `json:"deepSleep"`
	RemSleep        float64 `json:"remSleep"`
	LightSleep      float64 `json:"lightSleep"`
	Awake           float64 `json:"awake"`
	Efficiency      float64 `json:"efficiency"`
}

// Workout is the structure for workout data
type Workout struct {
	ID        string `json:"id"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Duration  int    `json:"duration"`
	Calories  float64 `json:"calories"`
	Type      string `json:"type"`
	AvgHr     int    `json:"avgHr"`
}

// DietaryTrend is the structure for dietary trend data
type DietaryTrend struct {
	Date      string  `json:"date"`
	Calories  float64 `json:"calories"`
	Protein   float64 `json:"protein"`
	Carbs     float64 `json:"carbs"`
	Fat       float64 `json:"fat"`
	Trend     float64 `json:"trend"`
}

// Meal is the structure for meal data
type Meal struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Cal  int    `json:"cal"`
}

// BodyComposition is the structure for body composition data
type BodyComposition struct {
	Time    string  `json:"time"`
	Weight  float64 `json:"weight"`
	BodyFat float64 `json:"bodyFat"`
}
