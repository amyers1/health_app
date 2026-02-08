export interface DailySummary {
    steps: number;
    distance: number;
    activeCalories: number;
    basalCalories: number;
    dietaryCalories: number;
}

export interface TimeSeriesData {
    time: string;
    value: number;
    [key: string]: string | number;
}

export interface BloodPressureData {
    time: string;
    systolic: number;
    diastolic: number;
    category: string;
}

export interface Workout {
    id: string;
    time: string;
    name: string;
    duration: number; // minutes
    calories: number;
    type: "run" | "cycle" | "weights" | "yoga" | "swim";
    avgHr?: number;
}

export interface SleepData {
    date: string;
    totalDuration: number; // hours
    deepSleep: number; // hours
    remSleep: number; // hours
    lightSleep: number; // hours
    awake: number; // hours
    efficiency: number; // percentage
}

export interface DietaryData {
    date: string;
    calories: number;
    protein: number; // grams
    carbs: number; // grams
    fat: number; // grams
}

export interface BodyComposition {
    time: string;
    weight: number;
    body_fat: number;
    muscle_mass: number;
}

export enum DashboardTab {
    OVERVIEW = "overview",
    ACTIVITY = "activity",
    HEALTH = "health",
    BODY = "body",
    SLEEP = "sleep",
    WORKOUTS = "workouts",
    DIETARY = "dietary",
}
