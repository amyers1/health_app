package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
	"health_app/api/model"

	"github.com/InfluxCommunity/influxdb3-go/v2/influxdb3"
	"github.com/joho/godotenv"
)

var easternZone, _ = time.LoadLocation("America/New_York")

type InfluxDBStore struct {
	client *influxdb3.Client
	bucket string
	org    string
}

func NewInfluxDBStore() (*InfluxDBStore, error) {
	url := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_DATABASE")

	if url == "" || token == "" || org == "" || bucket == "" {
		log.Printf("Could not load environment variables...attempting to load from .env file")

		err := godotenv.Load("../../../.env")
		if err != nil {
			return nil, fmt.Errorf("INFLUX_HOST, INFLUX_TOKEN, INFLUX_ORG, and INFLUX_DATABASE must be set")
		}
		url = "http://10.0.0.9:8181"
		token = os.Getenv("INFLUX_TOKEN")
		org = os.Getenv("INFLUX_ORG")
		bucket = os.Getenv("INFLUX_DATABASE")
	}

	// For Debug
	log.Printf("Connecting to InfluxDB at: %s (org: %s, bucket: %s)", url, org, bucket)

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:         url,
		Token:        token,
		Database:     bucket,
		Organization: org,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create InfluxDB client: %w", err)
	}

	return &InfluxDBStore{
		client: client,
		bucket: bucket,
		org:    org,
	}, nil
}

func (s *InfluxDBStore) Close() {
	if s.client != nil {
		log.Println("Closing InfluxDB client...")
		s.client.Close()
		log.Println("InfluxDB client closed successfully")
	}
}

func (s *InfluxDBStore) Ingest(metrics []model.Metric) error {
	// Convert metrics to line protocol format
	var lineProtocol string
	for _, m := range metrics {
		// Build tags string
		tagStr := ""
		for k, v := range m.Tags {
			if tagStr != "" {
				tagStr += ","
			}
			tagStr += fmt.Sprintf("%s=%s", k, v)
		}

		// Build fields string
		fieldStr := ""
		for k, v := range m.Fields {
			if fieldStr != "" {
				fieldStr += ","
			}
			switch val := v.(type) {
			case string:
				fieldStr += fmt.Sprintf(`%s="%s"`, k, val)
			case float64:
				fieldStr += fmt.Sprintf("%s=%f", k, val)
			case int64:
				fieldStr += fmt.Sprintf("%s=%di", k, val)
			case int:
				fieldStr += fmt.Sprintf("%s=%di", k, val)
			case bool:
				fieldStr += fmt.Sprintf("%s=%t", k, val)
			default:
				fieldStr += fmt.Sprintf("%s=%v", k, val)
			}
		}

		// Build line protocol: measurement[,tag=value...] field=value[,field=value...] [timestamp]
		line := m.Measurement
		if tagStr != "" {
			line += "," + tagStr
		}
		line += " " + fieldStr
		if !m.Timestamp.IsZero() {
			line += fmt.Sprintf(" %d", m.Timestamp.UnixNano())
		}
		lineProtocol += line + "\n"
	}

	return s.client.Write(context.Background(), []byte(lineProtocol))
}

func (s *InfluxDBStore) query(ctx context.Context, query string) (*influxdb3.QueryIterator, error) {
	return s.client.Query(ctx, query)
}

func (s *InfluxDBStore) GetSummary(date string) (*model.Summary, error) {
	start, stop := getDayRangeUTC(date)
	summary := &model.Summary{}

	query := fmt.Sprintf(`
        SELECT metric, source, value
        FROM "daily_totals"
        WHERE time >= '%s' AND time < '%s'
    `, start, stop)

	query2 := fmt.Sprintf(`
        SELECT qty
        FROM "dietary_energy"
        WHERE time >= '%s' AND time < '%s'
    `, start, stop)

	result, err := s.query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		record := result.Value()
		metric, okMetric := record["metric"].(string)
		source, _ := record["source"].(string)
		value := record["value"]

		if !okMetric || value == nil {
			continue
		}

		var floatValue float64
		switch v := value.(type) {
		case float64:
			floatValue = v
		case int64:
			floatValue = float64(v)
		default:
			continue
		}

		switch metric {
		case "step_count":
			if source == "RingConn" {
				summary.Steps = int(floatValue)
			}
		case "walking_running_distance":
			summary.Distance = floatValue
		case "active_energy":
			if source == "RingConn" {
				summary.ActiveCalories = floatValue
			}
		case "basal_energy_burned":
			if source == "RingConn" {
				summary.BasalCalories = floatValue
			}
		}
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	result2, err := s.query(context.Background(), query2)
	if err != nil {
		return nil, err
	}

	summary.DietaryCalories = 0.0
	for result2.Next() {
		record := result2.Value()
		value := record["qty"]

		if value == nil {
			continue
		}

		var floatValue float64
		switch v := value.(type) {
		case float64:
			floatValue = v
		case int64:
			floatValue = float64(v)
		default:
			continue
		}

		summary.DietaryCalories += floatValue
	}

	if result2.Err() != nil {
		return nil, result2.Err()
	}

	return summary, nil
}

func (s *InfluxDBStore) GetVitalsHR(date string) ([]model.TimeSeriesValue, error) {
	// Match Python behavior: use rolling 24-hour window from now
	now := time.Now().UTC()
	stop := now.Format(time.RFC3339)
	start := now.Add(-24 * time.Hour).Format(time.RFC3339)

	sqlQuery := fmt.Sprintf(`
SELECT time, "avg" as value
FROM "heart_rate"
WHERE time > '%s' AND time <= '%s'
ORDER BY time`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var values []model.TimeSeriesValue
	for result.Next() {
		record := result.Value()
		val, okVal := record["value"].(float64)
		t, okTime := record["time"].(time.Time)
		if okVal && okTime {
			values = append(values, model.TimeSeriesValue{
				Time:  t.UTC().Format("2006-01-02T15:04:05Z"),
				Value: val,
			})
		}
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	// Aggregate into 10-minute buckets
	buckets := make(map[time.Time][]float64)
	for _, v := range values {
		t, _ := time.Parse("2006-01-02T15:04:05Z", v.Time)
		bucketTime := t.Truncate(10 * time.Minute)
		buckets[bucketTime] = append(buckets[bucketTime], v.Value)
	}

	var aggregatedValues []model.TimeSeriesValue
	for t, vals := range buckets {
		var sum float64
		for _, v := range vals {
			sum += v
		}
		avg := sum / float64(len(vals))
		aggregatedValues = append(aggregatedValues, model.TimeSeriesValue{
			Time:  t.In(easternZone).Format("15:04"),
			Value: avg,
		})
	}

	sort.Slice(aggregatedValues, func(i, j int) bool {
		return aggregatedValues[i].Time < aggregatedValues[j].Time
	})

	return aggregatedValues, nil
}

func (s *InfluxDBStore) GetVitalsBP(endDate string) ([]model.BloodPressure, error) {
	start, stop := getDaysRangeUTC(endDate, 30)

	log.Printf("Querying blood pressure: start=%s, stop=%s", start, stop)

	sqlQuery := fmt.Sprintf(`
SELECT time, systolic, diastolic
FROM "blood_pressure"
WHERE time > '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		log.Printf("Blood pressure query error: %v", err)
		return nil, err
	}

	var bps []model.BloodPressure
	for result.Next() {
		record := result.Value()

		// Handle both int64 and float64 types from InfluxDB
		var systolic, diastolic int

		switch v := record["systolic"].(type) {
		case int64:
			systolic = int(v)
		case float64:
			systolic = int(v)
		default:
			log.Printf("Unexpected systolic type: %T", v)
			continue
		}

		switch v := record["diastolic"].(type) {
		case int64:
			diastolic = int(v)
		case float64:
			diastolic = int(v)
		default:
			log.Printf("Unexpected diastolic type: %T", v)
			continue
		}

		t, okTime := record["time"].(time.Time)
		if !okTime {
			log.Printf("Invalid time in blood pressure record")
			continue
		}

		bp := model.BloodPressure{
			Time:      t.In(easternZone).Format("Jan 02"),
			Systolic:  systolic,
			Diastolic: diastolic,
			Category:  getBPCategory(systolic, diastolic),
		}
		bps = append(bps, bp)
	}

	if result.Err() != nil {
		log.Printf("Blood pressure iteration error: %v", result.Err())
		return nil, result.Err()
	}

	log.Printf("Found %d blood pressure records", len(bps))
	return bps, nil
}

func (s *InfluxDBStore) GetVitalsGlucose(endDate string) ([]model.Glucose, error) {
	start, stop := getDaysRangeUTC(endDate, 30)
	sqlQuery := fmt.Sprintf(`
SELECT time, qty as value
FROM "blood_glucose"
WHERE time > '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var glucoses []model.Glucose
	for result.Next() {
		record := result.Value()
		value, okVal := record["value"].(float64)
		t, okTime := record["time"].(time.Time)
		if okVal && okTime {
			glucoses = append(glucoses, model.Glucose{
				Time:  t.In(easternZone).Format("Jan 02"),
				Value: value,
			})
		}
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return glucoses, nil
}

func (s *InfluxDBStore) GetSleep(endDate string) ([]model.Sleep, error) {
	start, stop := getDaysRangeUTC(endDate, 7)
	sqlQuery := fmt.Sprintf(`
SELECT time, "totalSleep", "deep", "rem", "core", "awake"
FROM "sleep_analysis"
WHERE time > '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var sleeps []model.Sleep
	for result.Next() {
		record := result.Value()
		t, okTime := record["time"].(time.Time)
		total, okTotal := record["totalSleep"].(float64)
		deep, okDeep := record["deep"].(float64)
		rem, okRem := record["rem"].(float64)
		light, okLight := record["core"].(float64)
		awake, okAwake := record["awake"].(float64)

		if okTime && okTotal && okDeep && okRem && okLight && okAwake {
			sleeps = append(sleeps, model.Sleep{
				Date:          t.In(easternZone).Format("Jan 02"),
				TotalDuration: total,
				DeepSleep:     deep,
				RemSleep:      rem,
				LightSleep:    light,
				Awake:         awake,
				Efficiency:    95, // Hardcoded as per python
			})
		}
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return sleeps, nil
}

func (s *InfluxDBStore) GetWorkouts(date string) ([]model.Workout, error) {
	start, stop := getDaysRangeUTC(date, 90)
	sqlQuery := fmt.Sprintf(`
SELECT workout_id, time, workout_name, duration, active_energy_value
FROM "workout"
WHERE time > '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	workoutsMap := make(map[string]model.Workout)
	var workoutIDs []string
	for result.Next() {
		record := result.Value()
		workoutID, _ := record["workout_id"].(string)
		t, _ := record["time"].(time.Time)
		name, _ := record["workout_name"].(string)
		duration, _ := record["duration"].(int64)
		calories, _ := record["active_energy_value"].(int64)

		workoutsMap[workoutID] = model.Workout{
			ID:       workoutID,
			Time:     t.In(easternZone).Format("2006-01-02 15:04"),
			Name:     name,
			Duration: int(duration / 60),
			Calories: float64(calories),
			Type:     name,
		}
		workoutIDs = append(workoutIDs, workoutID)
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	hrQuery := fmt.Sprintf(`
        SELECT workout_id, avg("avg") as avg_hr
        FROM "workout_heart_rate"
        WHERE time > '%s' AND time <= '%s'
        GROUP BY workout_id`, start, stop)

	hrResult, err := s.query(context.Background(), hrQuery)
	if err != nil {
		return nil, err
	}

	for hrResult.Next() {
		record := hrResult.Value()
		workoutID, _ := record["workout_id"].(string)
		avgHr, _ := record["avg_hr"].(float64)

		if workout, ok := workoutsMap[workoutID]; ok {
			workout.AvgHr = int(avgHr)
			workoutsMap[workoutID] = workout
		}
	}
	if hrResult.Err() != nil {
		return nil, hrResult.Err()
	}

	var workouts []model.Workout
	for _, id := range workoutIDs {
		workouts = append(workouts, workoutsMap[id])
	}

	return workouts, nil
}

type dailyNutrient struct {
	calories float64
	protein  float64
	carbs    float64
	fat      float64
}

func (s *InfluxDBStore) GetDietaryTrends(endDate string) ([]model.DietaryTrend, error) {
	_, stop := getDaysRangeUTC(endDate, 30)
	trendStart, _ := getDaysRangeUTC(endDate, 37)

	nutrients := []string{"dietary_energy", "protein", "carbohydrates", "total_fat"}

	// 1. Fetch all raw data points
	dailyData := make(map[string]*dailyNutrient)

	for _, nutrient := range nutrients {
		queryRangeStart := trendStart

		sqlQuery := fmt.Sprintf(`SELECT time, qty FROM "%s" WHERE time > '%s' AND time <= '%s'`, nutrient, queryRangeStart, stop)

		result, err := s.query(context.Background(), sqlQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to query nutrient %s: %w", nutrient, err)
		}

		for result.Next() {
			record := result.Value()
			t, _ := record["time"].(time.Time)
			value, _ := record["qty"].(float64)

			dayStr := t.In(easternZone).Format("2006-01-02")
			if _, ok := dailyData[dayStr]; !ok {
				dailyData[dayStr] = &dailyNutrient{}
			}

			switch nutrient {
			case "dietary_energy":
				dailyData[dayStr].calories += value
			case "protein":
				dailyData[dayStr].protein += value
			case "carbohydrates":
				dailyData[dayStr].carbs += value
			case "total_fat":
				dailyData[dayStr].fat += value
			}
		}
		if result.Err() != nil {
			return nil, result.Err()
		}
	}

	// 2. Calculate rolling average for trend (matching Python's behavior)
	var sortedDays []string
	for dayStr := range dailyData {
		sortedDays = append(sortedDays, dayStr)
	}
	sort.Strings(sortedDays)

	trendValues := make(map[string]float64)
	calorieHistory := []float64{}
	dayHistory := []string{}

	for _, dayStr := range sortedDays {
		calorieHistory = append(calorieHistory, dailyData[dayStr].calories)
		dayHistory = append(dayHistory, dayStr)
		if len(calorieHistory) > 7 {
			calorieHistory = calorieHistory[1:]
			dayHistory = dayHistory[1:]
		}

		sum := 0.0
		for _, v := range calorieHistory {
			sum += v
		}

		if len(calorieHistory) >= 3 {
			trendValues[dayStr] = sum / float64(len(calorieHistory))
		}
	}

	// 3. Build final response with forward-fill for missing trend values (matching Python)
	var trends []model.DietaryTrend
	var lastTrend float64 = 0

	endDateT, _ := time.ParseInLocation("2006-01-02", endDate, easternZone)
	startDateT := endDateT.AddDate(0, 0, -29)

	for d := startDateT; !d.After(endDateT); d = d.AddDate(0, 0, 1) {
		dayStr := d.Format("2006-01-02")

		data := &dailyNutrient{}
		if val, ok := dailyData[dayStr]; ok {
			data = val
		}

		// Forward-fill trend values (matching Python's fill_null(strategy='forward'))
		if trend, ok := trendValues[dayStr]; ok {
			lastTrend = trend
		}

		trends = append(trends, model.DietaryTrend{
			Date:     d.Format("Jan 02"),
			Calories: data.calories,
			Protein:  data.protein,
			Carbs:    data.carbs,
			Fat:      data.fat,
			Trend:    lastTrend,
		})
	}

	return trends, nil
}

func (s *InfluxDBStore) GetDietaryMealsToday(date string) ([]model.Meal, error) {
	// The schema does not clearly support this query. Returning placeholder data.
	return []model.Meal{
		{Name: "Breakfast", Desc: "Oatmeal, Berries, Whey", Cal: 420},
		{Name: "Lunch", Desc: "Chicken Salad, Quinoa", Cal: 580},
	}, nil
}

func (s *InfluxDBStore) GetBodyComposition(endDate string) ([]model.BodyComposition, error) {
	start, stop := getDaysRangeUTC(endDate, 30)

	// 1. Fetch weight data into a map keyed by timestamp
	weightMap := make(map[time.Time]float64)
	weightQuery := fmt.Sprintf(`SELECT time, qty as weight FROM "weight_body_mass" WHERE time > '%s' AND time <= '%s'`, start, stop)
	weightResult, err := s.query(context.Background(), weightQuery)
	if err != nil {
		return nil, fmt.Errorf("weight query error: %w", err)
	}

	for weightResult.Next() {
		record := weightResult.Value()
		t, okTime := record["time"].(time.Time)
		weight, okWeight := record["weight"].(float64)
		if okTime && okWeight {
			weightMap[t] = weight
		}
	}
	if weightResult.Err() != nil {
		return nil, weightResult.Err()
	}

	log.Printf("Found %d weight records", len(weightMap))

	// 2. Fetch body fat data and perform an inner join with weight data
	var compositions []model.BodyComposition
	bfQuery := fmt.Sprintf(`SELECT time, qty as bodyFat FROM "body_fat_percentage" WHERE time > '%s' AND time <= '%s'`, start, stop)
	bfResult, err := s.query(context.Background(), bfQuery)
	if err != nil {
		return nil, fmt.Errorf("body fat query error: %w", err)
	}

	for bfResult.Next() {
		record := bfResult.Value()
		t, okTime := record["time"].(time.Time)
		bodyFat, okBF := record["bodyfat"].(float64)

		// Check for matching weight measurement at the same timestamp (inner join)
		if okTime && okBF {
			if weight, ok := weightMap[t]; ok {
				compositions = append(compositions, model.BodyComposition{
					T:       t,
					Time:    t.In(easternZone).Format("Jan 02"),
					Weight:  weight,
					BodyFat: bodyFat,
				})
			}
		}
	}

	if bfResult.Err() != nil {
		log.Printf("Body fat iteration error: %v", bfResult.Err())
		return nil, bfResult.Err()
	}

	log.Printf("Found %d composition records", len(compositions))

	// Sort results by time ascending
	sort.Slice(compositions, func(i, j int) bool {
		return compositions[i].T.Before(compositions[j].T)
	})

	for _, value := range compositions {
		t := value.T.Format(time.RFC3339)
		w := value.Weight
		b := value.BodyFat
		log.Printf("comp records: time = %s, weight= %f, bf= %f\n", t, w, b)
	}

	return compositions, nil
}

// --- Helper Functions ---

// getDayRangeUTC returns UTC timestamps for the start and end of a day in Eastern time
func getDayRangeUTC(dateStr string) (string, string) {
	// Parse date in Eastern timezone
	t, _ := time.ParseInLocation("2006-01-02", dateStr, easternZone)

	// Create start and end times in Eastern
	startEastern := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, easternZone)
	endEastern := startEastern.Add(24 * time.Hour)

	// Convert to UTC for InfluxDB query
	startUTC := startEastern.UTC().Format(time.RFC3339)
	endUTC := endEastern.UTC().Format(time.RFC3339)

	return startUTC, endUTC
}

// getDaysRangeUTC returns UTC timestamps for a range of days ending on endDate
func getDaysRangeUTC(endDateStr string, days int) (string, string) {
	// Parse end date in Eastern timezone
	endDate, _ := time.ParseInLocation("2006-01-02", endDateStr, easternZone)

	// Create end of day in Eastern
	endEastern := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, easternZone)

	// Calculate start date
	startDate := endDate.AddDate(0, 0, -days+1)
	startEastern := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, easternZone)

	// Convert to UTC for InfluxDB query
	startUTC := startEastern.UTC().Format(time.RFC3339)
	stopUTC := endEastern.UTC().Format(time.RFC3339)

	return startUTC, stopUTC
}

func getBPCategory(systolic, diastolic int) string {
	if systolic > 180 || diastolic > 120 {
		return "Hypertensive Crisis"
	}
	if systolic >= 140 || diastolic >= 90 {
		return "Hypertension Stage 2"
	}
	if (systolic >= 130 && systolic <= 139) || (diastolic >= 80 && diastolic <= 89) {
		return "Hypertension Stage 1"
	}
	if systolic >= 120 && systolic <= 129 && diastolic < 80 {
		return "Elevated"
	}
	if systolic < 120 && diastolic < 80 {
		return "Normal"
	}
	return "Unknown"
}
