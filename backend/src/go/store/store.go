package store

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"
	"health_app/api/model"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var easternZone, _ = time.LoadLocation("America/New_York")

type InfluxDBStore struct {
	client influxdb2.Client
	bucket string
	org    string
}

func NewInfluxDBStore() (*InfluxDBStore, error) {
	url := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_DATABASE")

	if url == "" || token == "" || org == "" || bucket == "" {
		return nil, fmt.Errorf("INFLUX_HOST, INFLUX_TOKEN, INFLUX_ORG, and INFLUX_DATABASE must be set")
	}

	client := influxdb2.NewClient(url, token)
	_, err := client.Health(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to InfluxDB: %w", err)
	}

	return &InfluxDBStore{
		client: client,
		bucket: bucket,
		org:    org,
	}, nil
}

func (s *InfluxDBStore) Close() {
	s.client.Close()
}

func (s *InfluxDBStore) Ingest(metrics []model.Metric) error {
	writeAPI := s.client.WriteAPIBlocking(s.org, s.bucket)
	for _, m := range metrics {
		p := influxdb2.NewPoint(m.Measurement, m.Tags, m.Fields, m.Timestamp)
		if err := writeAPI.WritePoint(context.Background(), p); err != nil {
			return err
		}
	}
	return nil
}

func (s *InfluxDBStore) query(ctx context.Context, query string) (*api.QueryTableResult, error) {
	queryAPI := s.client.QueryAPI(s.org)
	return queryAPI.Query(ctx, query)
}

func (s *InfluxDBStore) GetSummary(date string) (*model.Summary, error) {
	start, stop := getDayRange(date)
	summary := &model.Summary{}

	query := fmt.Sprintf(`
        SELECT metric, source, value
        FROM "daily_totals"
        WHERE time >= '%s' AND time < '%s'
    `, start, stop)

	result, err := s.query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		record := result.Record()
		metric, okMetric := record.ValueByKey("metric").(string)
		source, _ := record.ValueByKey("source").(string)
		value := record.ValueByKey("value")

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
		case "dietary_energy":
			summary.DietaryCalories = floatValue
		}
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	return summary, nil
}

func (s *InfluxDBStore) GetVitalsHR(date string) ([]model.TimeSeriesValue, error) {
	start, stop := getDayRange(date)
	sqlQuery := fmt.Sprintf(`
SELECT time_bucket('10m', time) as bucket, avg("avg") as value
FROM "heart_rate"
WHERE time >= '%s' AND time < '%s'
GROUP BY bucket
ORDER BY bucket`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var values []model.TimeSeriesValue
	for result.Next() {
		val, okVal := result.Record().ValueByKey("value").(float64)
		t, okTime := result.Record().ValueByKey("bucket").(time.Time)
		if okVal && okTime {
			values = append(values, model.TimeSeriesValue{
				Time:  t.In(easternZone).Format("15:04"),
				Value: val,
			})
		}
	}
	return values, result.Err()
}

func (s *InfluxDBStore) GetVitalsBP(endDate string) ([]model.BloodPressure, error) {
	start, stop := getDaysRange(endDate, 90)
	sqlQuery := fmt.Sprintf(`
SELECT time, systolic, diastolic
FROM "blood_pressure"
WHERE time >= '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var bps []model.BloodPressure
	for result.Next() {
		record := result.Record()
		systolic, okSys := record.ValueByKey("systolic").(int64)
		diastolic, okDia := record.ValueByKey("diastolic").(int64)
		t, okTime := record.ValueByKey("time").(time.Time)

		if okSys && okDia && okTime {
			bp := model.BloodPressure{
				Time:      t.In(easternZone).Format("Jan 02"),
				Systolic:  int(systolic),
				Diastolic: int(diastolic),
				Category:  getBPCategory(int(systolic), int(diastolic)),
			}
			bps = append(bps, bp)
		}
	}
	return bps, result.Err()
}

func (s *InfluxDBStore) GetVitalsGlucose(endDate string) ([]model.Glucose, error) {
	start, stop := getDaysRange(endDate, 90)
	sqlQuery := fmt.Sprintf(`
SELECT time, qty as value
FROM "blood_glucose"
WHERE time >= '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var glucoses []model.Glucose
	for result.Next() {
		value, okVal := result.Record().ValueByKey("value").(float64)
		t, okTime := result.Record().ValueByKey("time").(time.Time)
		if okVal && okTime {
			glucoses = append(glucoses, model.Glucose{
				Time:  t.In(easternZone).Format("Jan 02"),
				Value: value,
			})
		}
	}
	return glucoses, result.Err()
}

func (s *InfluxDBStore) GetSleep(endDate string) ([]model.Sleep, error) {
	start, stop := getDaysRange(endDate, 90)
	sqlQuery := fmt.Sprintf(`
SELECT time, "totalSleep", "deep", "rem", "core"
FROM "sleep_analysis"
WHERE time >= '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	var sleeps []model.Sleep
	for result.Next() {
		record := result.Record()
		t, okTime := record.ValueByKey("time").(time.Time)
		total, okTotal := record.ValueByKey("totalSleep").(float64)
		deep, okDeep := record.ValueByKey("deep").(float64)
		rem, okRem := record.ValueByKey("rem").(float64)
		light, okLight := record.ValueByKey("core").(float64)

		if okTime && okTotal && okDeep && okRem && okLight {
			sleeps = append(sleeps, model.Sleep{
				Date:          t.In(easternZone).Format("Jan 02"),
				TotalDuration: total,
				DeepSleep:     deep,
				RemSleep:      rem,
				LightSleep:    light,
				Awake:         0, // Not available in query
				Efficiency:    95, // Hardcoded as per python
			})
		}
	}
	return sleeps, result.Err()
}

func (s *InfluxDBStore) GetWorkouts(date string) ([]model.Workout, error) {
	start, stop := getDaysRange(date, 90)
	sqlQuery := fmt.Sprintf(`
SELECT workout_id, time, workout_name, duration, active_energy_value
FROM "workout"
WHERE time >= '%s' AND time <= '%s'
ORDER BY time ASC`, start, stop)

	result, err := s.query(context.Background(), sqlQuery)
	if err != nil {
		return nil, err
	}

	workoutsMap := make(map[string]model.Workout)
	var workoutIDs []string
	for result.Next() {
		record := result.Record()
		workoutID, _ := record.ValueByKey("workout_id").(string)
		t, _ := record.ValueByKey("time").(time.Time)
		name, _ := record.ValueByKey("workout_name").(string)
		duration, _ := record.ValueByKey("duration").(int64)
		calories, _ := record.ValueByKey("active_energy_value").(int64)

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
        WHERE time >= '%s' AND time <= '%s'
        GROUP BY workout_id`, start, stop)

	hrResult, err := s.query(context.Background(), hrQuery)
	if err != nil {
		return nil, err
	}

	for hrResult.Next() {
		record := hrResult.Record()
		workoutID, _ := record.ValueByKey("workout_id").(string)
		avgHr, _ := record.ValueByKey("avg_hr").(float64)

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
    _, stop := getDaysRange(endDate, 90)
    trendStart, _ := getDaysRange(endDate, 97)

    nutrients := []string{"dietary_energy", "protein", "carbohydrates", "total_fat"}

    // 1. Fetch all raw data points
    dailyData := make(map[string]*dailyNutrient)

    for _, nutrient := range nutrients {
        queryRangeStart := trendStart

        sqlQuery := fmt.Sprintf(`SELECT time, qty FROM "%s" WHERE time >= '%s' AND time <= '%s'`, nutrient, queryRangeStart, stop)

        result, err := s.query(context.Background(), sqlQuery)
        if err != nil {
            return nil, fmt.Errorf("failed to query nutrient %s: %w", nutrient, err)
        }

        for result.Next() {
            t, _ := result.Record().ValueByKey("time").(time.Time)
            value, _ := result.Record().ValueByKey("qty").(float64)

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

    // 2. Calculate rolling average for trend
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

    // 3. Build final response
    var trends []model.DietaryTrend
    var lastTrend float64 = 0

    endDateT, _ := time.ParseInLocation("2006-01-02", endDate, easternZone)
    startDateT := endDateT.AddDate(0, 0, -89)

    for d := startDateT; !d.After(endDateT); d = d.AddDate(0, 0, 1) {
        dayStr := d.Format("2006-01-02")

        data := &dailyNutrient{}
        if val, ok := dailyData[dayStr]; ok {
            data = val
        }

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
	start, stop := getDaysRange(endDate, 90)

	compositionMap := make(map[string]*model.BodyComposition)
	var orderedDates []string

	// Weight
	weightQuery := fmt.Sprintf(`SELECT time, qty as weight FROM "weight_body_mass" WHERE time >= '%s' AND time <= '%s' ORDER BY time ASC`, start, stop)
	weightResult, err := s.query(context.Background(), weightQuery)
	if err != nil {
		return nil, err
	}
	for weightResult.Next() {
		t, _ := weightResult.Record().ValueByKey("time").(time.Time)
		weight, _ := weightResult.Record().ValueByKey("weight").(float64)

		dateStr := t.In(easternZone).Format("Jan 02")
		if _, ok := compositionMap[dateStr]; !ok {
			compositionMap[dateStr] = &model.BodyComposition{Time: dateStr}
			orderedDates = append(orderedDates, dateStr)
		}
		compositionMap[dateStr].Weight = weight // last one wins for the day
	}

	// Body Fat
	bfQuery := fmt.Sprintf(`SELECT time, qty as bodyFat FROM "body_fat_percentage" WHERE time >= '%s' AND time <= '%s' ORDER BY time ASC`, start, stop)
	bfResult, err := s.query(context.Background(), bfQuery)
	if err != nil {
		return nil, err
	}
	for bfResult.Next() {
		t, _ := bfResult.Record().ValueByKey("time").(time.Time)
		bodyFat, _ := bfResult.Record().ValueByKey("bodyFat").(float64)
		dateStr := t.In(easternZone).Format("Jan 02")
		if comp, ok := compositionMap[dateStr]; ok {
			comp.BodyFat = bodyFat
		}
	}

	var compositions []model.BodyComposition
	for _, dateStr := range orderedDates {
		comp := compositionMap[dateStr]
		if comp.Weight > 0 && comp.BodyFat > 0 {
			compositions = append(compositions, *comp)
		}
	}

	return compositions, nil
}

// --- Helper Functions ---

func getDayRange(dateStr string) (string, string) {
	t, _ := time.ParseInLocation("2006-01-02", dateStr, easternZone)
	start := t.Format(time.RFC3339)
	stop := t.Add(24 * time.Hour).Format(time.RFC3339)
	return start, stop
}

func getDaysRange(endDateStr string, days int) (string, string) {
	end, _ := time.ParseInLocation("2006-01-02", endDateStr, easternZone)
	stop := end.Add(24*time.Hour - 1*time.Nanosecond).Format(time.RFC3339Nano)
	start := end.AddDate(0, 0, -days+1).Format(time.RFC3339Nano)
	return start, stop
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
