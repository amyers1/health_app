
import { DailySummary, TimeSeriesData, BloodPressureData, Workout, SleepData, DietaryData } from '../types';

export const getDailySummary = (): DailySummary => ({
  steps: 12458,
  distance: 6.2,
  activeCalories: 842,
  basalCalories: 1850,
  dietaryCalories: 2150,
});

export const getHeartRateData = (): TimeSeriesData[] => {
  const data = [];
  const now = new Date();
  for (let i = 24; i >= 0; i--) {
    const time = new Date(now.getTime() - i * 60 * 60 * 1000);
    data.push({
      time: time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
      value: Math.floor(60 + Math.random() * 40),
    });
  }
  return data;
};

export const getStepData = (): TimeSeriesData[] => {
  const data = [];
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
  for (const day of days) {
    data.push({
      time: day,
      value: Math.floor(5000 + Math.random() * 10000),
    });
  }
  return data;
};

export const getBloodPressureData = (): BloodPressureData[] => {
  const data: BloodPressureData[] = [];
  const now = new Date();
  for (let i = 29; i >= 0; i--) {
    const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000);
    const systolic = Math.floor(110 + Math.random() * 20);
    const diastolic = Math.floor(70 + Math.random() * 15);
    data.push({
      time: date.toLocaleDateString([], { month: 'short', day: 'numeric' }),
      systolic,
      diastolic,
      category: systolic < 120 && diastolic < 80 ? 'Normal' : 'Elevated'
    });
  }
  return data;
};

export const getBloodGlucoseData = (): TimeSeriesData[] => {
  const data: TimeSeriesData[] = [];
  const now = new Date();
  for (let i = 29; i >= 0; i--) {
    const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000);
    // Average daily values roughly between 90 and 130 mg/dL
    data.push({
      time: date.toLocaleDateString([], { month: 'short', day: 'numeric' }),
      value: Math.floor(90 + Math.random() * 40),
    });
  }
  return data;
};

export const getWorkouts = (): Workout[] => [
  { id: '1', time: 'Today, 8:30 AM', name: 'Morning Run', duration: 45, calories: 420, type: 'run', avgHr: 155 },
  { id: '2', time: 'Yesterday, 6:00 PM', name: 'Weightlifting', duration: 60, calories: 350, type: 'weights', avgHr: 125 },
  { id: '3', time: '2 days ago, 7:15 AM', name: 'HIIT Session', duration: 30, calories: 310, type: 'run', avgHr: 168 },
  { id: '4', time: '4 days ago, 5:30 PM', name: 'Evening Yoga', duration: 40, calories: 120, type: 'yoga', avgHr: 95 },
  { id: '5', time: '5 days ago, 6:00 AM', name: 'Swim Laps', duration: 45, calories: 480, type: 'swim', avgHr: 142 },
];

export const getBodyTrends = (): TimeSeriesData[] => {
  const data = [];
  const now = new Date();
  for (let i = 12; i >= 0; i--) {
    const date = new Date(now.getTime() - i * 7 * 24 * 60 * 60 * 1000);
    data.push({
      time: date.toLocaleDateString([], { month: 'short', day: 'numeric' }),
      weight: 185 - (i * 0.2) + Math.random(),
      bodyFat: 18 - (i * 0.05) + Math.random() * 0.5,
    });
  }
  return data;
};

export const getSleepHistory = (): SleepData[] => {
  const data = [];
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
  for (const day of days) {
    const total = 6 + Math.random() * 3;
    data.push({
      date: day,
      totalDuration: total,
      deepSleep: total * 0.2,
      remSleep: total * 0.25,
      lightSleep: total * 0.45,
      awake: total * 0.1,
      efficiency: 85 + Math.random() * 10,
    });
  }
  return data;
};

export const getDietaryTrends = (): (DietaryData & { trend?: number })[] => {
  const data: (DietaryData & { trend?: number })[] = [];
  const now = new Date();
  
  // Generate 30 days of data
  for (let i = 29; i >= 0; i--) {
    const date = new Date(now.getTime() - i * 24 * 60 * 60 * 1000);
    const calories = 1800 + Math.random() * 600;
    data.push({
      date: date.toLocaleDateString([], { month: 'short', day: 'numeric' }),
      calories: Math.round(calories),
      protein: 120 + Math.random() * 40,
      carbs: 200 + Math.random() * 50,
      fat: 60 + Math.random() * 20,
    });
  }

  // Calculate a simple 7-day moving average for the trendline
  for (let i = 0; i < data.length; i++) {
    const start = Math.max(0, i - 6);
    const subset = data.slice(start, i + 1);
    const sum = subset.reduce((acc, curr) => acc + curr.calories, 0);
    data[i].trend = Math.round(sum / subset.length);
  }

  return data;
};
