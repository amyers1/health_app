import React, { useState, useEffect, useRef } from "react";
import {
    Activity,
    Heart,
    User,
    Flame,
    Zap,
    Footprints,
    Moon,
    Scale,
    Clock,
    Calendar,
    Utensils,
    ChevronRight,
    TrendingUp,
    Droplets,
    Dumbbell,
    Coffee,
    PieChart as PieChartIcon,
    Syringe,
    ChevronLeft,
} from "lucide-react";
import {
    DashboardTab,
    DailySummary,
    TimeSeriesData,
    BloodPressureData,
    Workout,
    SleepData,
    DietaryData,
    BodyComposition,
} from "./types";
import MetricCard from "./components/MetricCard";
import WorkoutList from "./components/WorkoutList";
import {
    HeartRateChart,
    StepBarChart,
    BloodPressureChart,
    SleepStagesChart,
    MacroTrendsChart,
    MonthlyCalorieChart,
    BloodGlucoseChart,
    BodyWeightChart,
    BodyFatPercentageChart,
} from "./components/Charts";

const App: React.FC = () => {
    const [activeTab, setActiveTab] = useState<DashboardTab>(
        DashboardTab.OVERVIEW,
    );
    const [selectedDate, setSelectedDate] = useState<string>(
        new Date().toISOString().split("T")[0],
    );
    const [summary, setSummary] = useState<DailySummary | null>(null);
    const [hrData, setHrData] = useState<TimeSeriesData[]>([]);
    const [bpData, setBpData] = useState<BloodPressureData[]>([]);
    const [glucoseData, setGlucoseData] = useState<TimeSeriesData[]>([]);
    const [workouts, setWorkouts] = useState<Workout[]>([]);
    const [bodyTrends, setBodyTrends] = useState<BodyComposition[]>([]);
    const [sleepHistory, setSleepHistory] = useState<SleepData[]>([]);
    const [dietaryTrends, setDietaryTrends] = useState<
        (DietaryData & { trend?: number })[]
    >([]);

    const dateInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const safeJson = (res: Response) => {
                    if (!res.ok) {
                        console.error("API Error:", res.status, res.statusText);
                        return null;
                    }
                    return res.json();
                };

                const summaryPromise = fetch(
                    `/api/v1/summary?date=${selectedDate}`,
                ).then(safeJson);
                const hrPromise = fetch(
                    `/api/v1/vitals/hr?date=${selectedDate}`,
                ).then(safeJson);
                const bpPromise = fetch(
                    `/api/v1/vitals/bp?end_date=${selectedDate}`,
                ).then(safeJson);
                const glucosePromise = fetch(
                    `/api/v1/vitals/glucose?end_date=${selectedDate}`,
                ).then(safeJson);
                const workoutsPromise = fetch(
                    `/api/v1/workouts?date=${selectedDate}`,
                ).then(safeJson);
                const bodyPromise = fetch(
                    `/api/v1/body/composition?end_date=${selectedDate}`,
                ).then(safeJson);
                const sleepPromise = fetch(
                    `/api/v1/sleep?end_date=${selectedDate}`,
                ).then(safeJson);
                const dietPromise = fetch(
                    `/api/v1/dietary/trends?end_date=${selectedDate}`,
                ).then(safeJson);

                const [summary, hr, bp, glucose, workouts, body, sleep, diet] =
                    await Promise.all([
                        summaryPromise,
                        hrPromise,
                        bpPromise,
                        glucosePromise,
                        workoutsPromise,
                        bodyPromise,
                        sleepPromise,
                        dietPromise,
                    ]);

                setSummary(summary || null);
                setHrData(hr || []);
                setBpData(bp || []);
                setGlucoseData(glucose || []);
                setWorkouts(workouts || []);
                setBodyTrends(body || []);
                setSleepHistory(sleep || []);
                setDietaryTrends(diet || []);
            } catch (error) {
                console.error("Failed to fetch data:", error);
            }
        };

        fetchData();
    }, [selectedDate]);

    const handleDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSelectedDate(e.target.value);
    };

    const shiftDate = (days: number) => {
        const date = new Date(selectedDate);
        date.setDate(date.getDate() + days);
        setSelectedDate(date.toISOString().split("T")[0]);
    };

    const formattedDisplayDate = new Date(selectedDate).toLocaleDateString(
        undefined,
        {
            month: "short",
            day: "numeric",
            year: "numeric",
        },
    );

    const isToday = selectedDate === new Date().toISOString().split("T")[0];

    const renderOverview = () => (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <MetricCard
                    title="Steps"
                    value={summary?.steps.toLocaleString() || 0}
                    unit=""
                    icon={<Footprints size={20} />}
                    colorClass="text-emerald-400 bg-emerald-400"
                    subtitle="Goal: 10,000"
                    trend={{ value: 12, isUp: true }}
                />
                <MetricCard
                    title="Calories"
                    value={summary?.activeCalories || 0}
                    unit="kcal"
                    icon={<Flame size={20} />}
                    colorClass="text-orange-400 bg-orange-400"
                    subtitle="Goal: 600"
                />
                <MetricCard
                    title="Heart Rate"
                    value={hrData[hrData.length - 1]?.value || "--"}
                    unit="bpm"
                    icon={<Heart size={20} />}
                    colorClass="text-rose-400 bg-rose-400"
                    subtitle="Resting: 62 bpm"
                />
                <MetricCard
                    title="Sleep"
                    value="7.5"
                    unit="hrs"
                    icon={<Moon size={20} />}
                    colorClass="text-indigo-400 bg-indigo-400"
                    subtitle="92% Quality Score"
                />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2 space-y-6">
                    <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                        <div className="flex justify-between items-center mb-6">
                            <h3 className="font-semibold text-slate-100 flex items-center gap-2">
                                <Clock size={18} className="text-rose-400" />
                                Heart Rate Intensity
                            </h3>
                        </div>
                        <HeartRateChart data={hrData} />
                    </section>
                </div>

                <div className="space-y-6">
                    <WorkoutList workouts={workouts.slice(0, 4)} />

                    <div className="bg-gradient-to-br from-indigo-600/20 to-violet-600/20 border border-indigo-500/20 rounded-2xl p-4">
                        <h3 className="text-sm font-semibold text-indigo-100 mb-2">
                            Daily Insight
                        </h3>
                        <p className="text-xs text-indigo-200/70 mb-4 leading-relaxed">
                            {isToday
                                ? "You're on track to hit your step goal. Keep moving!"
                                : `On ${formattedDisplayDate}, you achieved 92% of your activity target.`}
                        </p>
                        <div className="flex justify-between items-end">
                            <div>
                                <span className="text-2xl font-bold text-slate-50">
                                    92%
                                </span>
                                <p className="text-[10px] text-indigo-300">
                                    Goal Completion
                                </p>
                            </div>
                            <button className="bg-indigo-500 hover:bg-indigo-400 text-white p-2 rounded-xl transition-colors">
                                <ChevronRight size={18} />
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

    const renderSleep = () => (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <MetricCard
                    title="Selected Night"
                    value="7h 45m"
                    unit=""
                    icon={<Moon size={20} />}
                    colorClass="text-indigo-400 bg-indigo-400"
                    subtitle={`Morning of ${formattedDisplayDate}`}
                />
                <MetricCard
                    title="Deep Sleep"
                    value="1h 22m"
                    unit=""
                    icon={<Zap size={20} />}
                    colorClass="text-blue-400 bg-blue-400"
                    subtitle="Above average"
                />
                <MetricCard
                    title="Efficiency"
                    value="94"
                    unit="%"
                    icon={<Activity size={20} />}
                    colorClass="text-emerald-400 bg-emerald-400"
                    subtitle="8% Restlessness"
                />
            </div>

            <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                <h3 className="font-semibold text-slate-100 mb-6 flex items-center gap-2">
                    <Calendar size={18} className="text-indigo-400" />
                    Weekly Sleep Analysis
                </h3>
                <SleepStagesChart data={sleepHistory} />
            </section>

            <div className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6">
                <h4 className="text-sm font-semibold text-slate-300 mb-4">
                    Quality Metrics
                </h4>
                <div className="space-y-4">
                    <div className="flex gap-4 items-start">
                        <div className="w-8 h-8 rounded-lg bg-indigo-500/10 flex items-center justify-center text-indigo-400 flex-shrink-0">
                            <Clock size={16} />
                        </div>
                        <div>
                            <p className="text-sm text-slate-200">
                                Stability Score: 8.5/10
                            </p>
                            <p className="text-xs text-slate-500">
                                Your sleep onset on this day was highly
                                consistent with your monthly average.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );

    const renderWorkouts = () => (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-6">
                    <WorkoutList workouts={workouts} />
                </div>
                <div className="space-y-6">
                    <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6">
                        <h3 className="font-semibold text-slate-100 mb-4 flex items-center gap-2">
                            <Dumbbell size={18} className="text-blue-400" />
                            Volume Analysis (Week of {formattedDisplayDate})
                        </h3>
                        <div className="space-y-6">
                            {[
                                {
                                    type: "Running",
                                    hrs: 3.5,
                                    color: "bg-emerald-500",
                                },
                                {
                                    type: "Strength",
                                    hrs: 4.2,
                                    color: "bg-blue-500",
                                },
                                {
                                    type: "Yoga",
                                    hrs: 1.5,
                                    color: "bg-indigo-500",
                                },
                            ].map((item) => (
                                <div key={item.type} className="space-y-2">
                                    <div className="flex justify-between text-xs">
                                        <span className="text-slate-400">
                                            {item.type}
                                        </span>
                                        <span className="text-slate-100 font-medium">
                                            {item.hrs} hrs/wk
                                        </span>
                                    </div>
                                    <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                                        <div
                                            className={`${item.color} h-full`}
                                            style={{
                                                width: `${(item.hrs / 5) * 100}%`,
                                            }}
                                        ></div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </section>
                </div>
            </div>
        </div>
    );

    const renderDietary = () => (
        <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <MetricCard
                    title="Daily Cal"
                    value={summary?.dietaryCalories || 0}
                    unit="kcal"
                    icon={<Coffee size={20} />}
                    colorClass="text-yellow-400 bg-yellow-400"
                    subtitle="Net: -320 kcal"
                />
                <MetricCard
                    title="Protein"
                    value="142"
                    unit="g"
                    icon={<Zap size={20} />}
                    colorClass="text-rose-400 bg-rose-400"
                    subtitle="Goal: 160g"
                />
                <MetricCard
                    title="Carbs"
                    value="210"
                    unit="g"
                    icon={<PieChartIcon size={20} />}
                    colorClass="text-emerald-400 bg-emerald-400"
                    subtitle="Goal: 200g"
                />
                <MetricCard
                    title="Fat"
                    value="68"
                    unit="g"
                    icon={<Droplets size={20} />}
                    colorClass="text-blue-400 bg-blue-400"
                    subtitle="Goal: 70g"
                />
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                    <h3 className="font-semibold text-slate-100 mb-6 flex items-center gap-2">
                        <Calendar size={18} className="text-yellow-400" />
                        30-Day Calorie History
                    </h3>
                    <MonthlyCalorieChart data={dietaryTrends} />
                </section>

                <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                    <h3 className="font-semibold text-slate-100 mb-6 flex items-center gap-2">
                        <TrendingUp size={18} className="text-emerald-400" />
                        Weekly Macro Balance
                    </h3>
                    <MacroTrendsChart data={dietaryTrends.slice(-7)} />
                </section>
            </div>
        </div>
    );

    const renderContent = () => {
        switch (activeTab) {
            case DashboardTab.OVERVIEW:
                return renderOverview();
            case DashboardTab.SLEEP:
                return renderSleep();
            case DashboardTab.WORKOUTS:
                return renderWorkouts();
            case DashboardTab.DIETARY:
                return renderDietary();
            case DashboardTab.HEALTH:
                return (
                    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                                <h3 className="font-semibold text-slate-100 mb-6 flex items-center gap-2">
                                    <Droplets
                                        size={18}
                                        className="text-red-400"
                                    />
                                    Blood Pressure (30 Day Trend)
                                </h3>
                                <BloodPressureChart data={bpData} />
                            </section>
                            <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                                <h3 className="font-semibold text-slate-100 mb-6 flex items-center gap-2">
                                    <Heart
                                        size={18}
                                        className="text-rose-400"
                                    />
                                    Heart Rate Variation
                                </h3>
                                <HeartRateChart data={hrData} />
                            </section>
                        </div>
                        <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4">
                            <h3 className="font-semibold text-slate-100 mb-6 flex items-center gap-2">
                                <Syringe
                                    size={18}
                                    className="text-purple-400"
                                />
                                Blood Glucose (30 Day Trend)
                            </h3>
                            <BloodGlucoseChart data={glucoseData} />
                        </section>
                    </div>
                );
            case DashboardTab.BODY:
                const latestBodyData =
                    bodyTrends.length > 0
                        ? bodyTrends[bodyTrends.length - 1]
                        : null;
                return (
                    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-2 duration-500">
                        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                            <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6">
                                <BodyWeightChart data={bodyTrends} />
                            </section>
                            <section className="bg-slate-900/60 border border-slate-800 rounded-2xl p-6">
                                <BodyFatPercentageChart data={bodyTrends} />
                            </section>
                        </div>
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                            <MetricCard
                                title="Weight"
                                value={latestBodyData?.weight || 0}
                                unit="lbs"
                                icon={<Scale size={24} />}
                                colorClass="text-blue-400 bg-blue-400"
                                subtitle={`As of ${formattedDisplayDate}`}
                            />
                            <MetricCard
                                title="Body Fat"
                                value={latestBodyData?.body_fat || 0}
                                unit="%"
                                icon={<TrendingUp size={24} />}
                                colorClass="text-purple-400 bg-purple-400"
                                subtitle="Trend: Stable"
                            />
                            <MetricCard
                                title="BMI"
                                value="24.1"
                                unit=""
                                icon={<User size={24} />}
                                colorClass="text-emerald-400 bg-emerald-400"
                                subtitle="Standard range"
                            />
                        </div>
                    </div>
                );
            default:
                return renderOverview();
        }
    };

    return (
        <div className="min-h-screen bg-slate-950 pb-24 md:pb-6 flex flex-col">
            {/* Header */}
            <header className="sticky top-0 z-30 bg-slate-950/80 backdrop-blur-md border-b border-slate-800 px-4 py-3 md:px-8">
                <div className="max-w-7xl mx-auto flex justify-between items-center">
                    <div>
                        <h1 className="text-xl font-bold text-slate-50 flex items-center gap-2">
                            <Zap
                                className="text-indigo-500"
                                fill="currentColor"
                                size={24}
                            />
                            VitalStream
                        </h1>
                    </div>

                    {/* Enhanced Date Selector */}
                    <div className="flex items-center gap-1 sm:gap-4">
                        <div className="flex items-center bg-slate-900 border border-slate-800 rounded-xl overflow-hidden shadow-inner">
                            <button
                                onClick={() => shiftDate(-1)}
                                className="p-2 hover:bg-slate-800 text-slate-400 transition-colors border-r border-slate-800"
                            >
                                <ChevronLeft size={16} />
                            </button>

                            <div className="relative group">
                                <button
                                    onClick={() =>
                                        dateInputRef.current?.showPicker()
                                    }
                                    className="flex items-center gap-2 px-3 py-2 text-slate-100 hover:text-white transition-colors text-sm font-semibold"
                                >
                                    <Calendar
                                        size={14}
                                        className="text-indigo-400"
                                    />
                                    <span className="min-w-[100px] text-center">
                                        {isToday
                                            ? "Today"
                                            : formattedDisplayDate}
                                    </span>
                                </button>
                                <input
                                    ref={dateInputRef}
                                    type="date"
                                    value={selectedDate}
                                    onChange={handleDateChange}
                                    className="absolute inset-0 opacity-0 pointer-events-none"
                                    max={new Date().toISOString().split("T")[0]}
                                />
                            </div>

                            <button
                                onClick={() => shiftDate(1)}
                                disabled={isToday}
                                className={`p-2 hover:bg-slate-800 transition-colors border-l border-slate-800 ${isToday ? "text-slate-700" : "text-slate-400"}`}
                            >
                                <ChevronRight size={16} />
                            </button>
                        </div>

                        <div className="hidden sm:flex w-9 h-9 rounded-full bg-slate-800 border border-slate-700 items-center justify-center text-slate-300">
                            <User size={18} />
                        </div>
                    </div>
                </div>
            </header>

            {/* Main Content */}
            <main className="flex-1 w-full max-w-7xl mx-auto px-4 py-6 md:px-8 md:py-10">
                <div className="mb-8 overflow-x-auto custom-scrollbar">
                    <div className="flex bg-slate-900 p-1 rounded-xl w-max">
                        {[
                            {
                                id: DashboardTab.OVERVIEW,
                                icon: <Activity size={16} />,
                                label: "Summary",
                            },
                            {
                                id: DashboardTab.SLEEP,
                                icon: <Moon size={16} />,
                                label: "Sleep",
                            },
                            {
                                id: DashboardTab.WORKOUTS,
                                icon: <Dumbbell size={16} />,
                                label: "Training",
                            },
                            {
                                id: DashboardTab.DIETARY,
                                icon: <Utensils size={16} />,
                                label: "Dietary",
                            },
                            {
                                id: DashboardTab.HEALTH,
                                icon: <Heart size={16} />,
                                label: "Vitals",
                            },
                            {
                                id: DashboardTab.BODY,
                                icon: <Scale size={16} />,
                                label: "Composition",
                            },
                        ].map((tab) => (
                            <button
                                key={tab.id}
                                onClick={() =>
                                    setActiveTab(tab.id as DashboardTab)
                                }
                                className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all whitespace-nowrap ${
                                    activeTab === tab.id
                                        ? "bg-slate-800 text-slate-50 shadow-sm shadow-black/20"
                                        : "text-slate-500 hover:text-slate-300"
                                }`}
                            >
                                {tab.icon}
                                <span>{tab.label}</span>
                            </button>
                        ))}
                    </div>
                </div>

                {renderContent()}
            </main>

            {/* Mobile Bottom Navigation */}
            <nav className="fixed bottom-0 left-0 right-0 z-50 bg-slate-950/90 backdrop-blur-lg border-t border-slate-800 px-6 py-4 md:hidden">
                <div className="flex justify-between items-center max-w-md mx-auto">
                    <button
                        onClick={() => setActiveTab(DashboardTab.OVERVIEW)}
                        className={`flex flex-col items-center gap-1 ${activeTab === DashboardTab.OVERVIEW ? "text-indigo-400" : "text-slate-500"}`}
                    >
                        <Activity size={24} />
                        <span className="text-[10px] font-bold">Summary</span>
                    </button>
                    <button
                        onClick={() => setActiveTab(DashboardTab.SLEEP)}
                        className={`flex flex-col items-center gap-1 ${activeTab === DashboardTab.SLEEP ? "text-indigo-400" : "text-slate-500"}`}
                    >
                        <Moon size={24} />
                        <span className="text-[10px] font-bold">Sleep</span>
                    </button>
                    <button
                        className="flex flex-col items-center gap-1 text-slate-500"
                        onClick={() => setActiveTab(DashboardTab.WORKOUTS)}
                    >
                        <div className="bg-indigo-600 text-white p-3 rounded-full -mt-12 shadow-lg shadow-indigo-600/40 border-4 border-slate-950">
                            <Zap size={24} />
                        </div>
                    </button>
                    <button
                        onClick={() => setActiveTab(DashboardTab.DIETARY)}
                        className={`flex flex-col items-center gap-1 ${activeTab === DashboardTab.DIETARY ? "text-indigo-400" : "text-slate-500"}`}
                    >
                        <Utensils size={24} />
                        <span className="text-[10px] font-bold">Diet</span>
                    </button>
                    <button
                        onClick={() => setActiveTab(DashboardTab.HEALTH)}
                        className={`flex flex-col items-center gap-1 ${activeTab === DashboardTab.HEALTH ? "text-indigo-400" : "text-slate-500"}`}
                    >
                        <Heart size={24} />
                        <span className="text-[10px] font-bold">Vitals</span>
                    </button>
                </div>
            </nav>
        </div>
    );
};

export default App;
