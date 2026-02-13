import React from "react";
import {
    AreaChart,
    Area,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer,
    BarChart,
    Bar,
    Cell,
    LineChart,
    Line,
    Legend,
    ComposedChart,
    ReferenceLine,
} from "recharts";
import {
    TimeSeriesData,
    BloodPressureData,
    SleepData,
    DietaryData,
} from "../types";
import { formatNumber } from "../utils";

export const HeartRateChart: React.FC<{ data: TimeSeriesData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
                <defs>
                    <linearGradient id="colorHr" x1="0" y1="0" x2="0" y2="1">
                        <stop
                            offset="5%"
                            stopColor="#ef4444"
                            stopOpacity={0.3}
                        />
                        <stop
                            offset="95%"
                            stopColor="#ef4444"
                            stopOpacity={0}
                        />
                    </linearGradient>
                </defs>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="time"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    minTickGap={30}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    domain={["dataMin - 10", "dataMax + 10"]}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    itemStyle={{ color: "#f8fafc" }}
                    formatter={(value: number) => `${formatNumber(value)} bpm`}
                />
                <Area
                    type="monotone"
                    dataKey="value"
                    stroke="#ef4444"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#colorHr)"
                />
            </AreaChart>
        </ResponsiveContainer>
    </div>
);

export const StepBarChart: React.FC<{ data: TimeSeriesData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data}>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="time"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                />
                <Tooltip
                    cursor={{ fill: "#334155", opacity: 0.4 }}
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) =>
                        `${formatNumber(value)} steps`
                    }
                />
                <Bar dataKey="value" radius={[4, 4, 0, 0]}>
                    {data.map((entry, index) => (
                        <Cell
                            key={`cell-${index}`}
                            fill={entry.value > 8000 ? "#10b981" : "#3b82f6"}
                        />
                    ))}
                </Bar>
            </BarChart>
        </ResponsiveContainer>
    </div>
);

export const BloodPressureChart: React.FC<{ data: BloodPressureData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <LineChart data={data}>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="time"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    minTickGap={40}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    domain={[60, 160]}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) => `${formatNumber(value)} mmHg`}
                />
                <Legend verticalAlign="top" height={36} />
                <Line
                    type="monotone"
                    dataKey="systolic"
                    name="Systolic"
                    stroke="#f43f5e"
                    strokeWidth={2}
                    dot={{ r: 2, fill: "#f43f5e" }}
                    activeDot={{ r: 4 }}
                />
                <Line
                    type="monotone"
                    dataKey="diastolic"
                    name="Diastolic"
                    stroke="#3b82f6"
                    strokeWidth={2}
                    dot={{ r: 2, fill: "#3b82f6" }}
                    activeDot={{ r: 4 }}
                />
            </LineChart>
        </ResponsiveContainer>
    </div>
);

export const BloodGlucoseChart: React.FC<{ data: TimeSeriesData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
                <defs>
                    <linearGradient
                        id="colorGlucose"
                        x1="0"
                        y1="0"
                        x2="0"
                        y2="1"
                    >
                        <stop
                            offset="5%"
                            stopColor="#a855f7"
                            stopOpacity={0.3}
                        />
                        <stop
                            offset="95%"
                            stopColor="#a855f7"
                            stopOpacity={0}
                        />
                    </linearGradient>
                </defs>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="time"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    minTickGap={40}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    domain={[60, 180]}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    itemStyle={{ color: "#f8fafc" }}
                    formatter={(value: number) =>
                        `${formatNumber(value)} mg/dL`
                    }
                />
                <ReferenceLine
                    y={140}
                    label={{
                        position: "right",
                        value: "High",
                        fill: "#f43f5e",
                        fontSize: 10,
                    }}
                    stroke="#f43f5e"
                    strokeDasharray="3 3"
                />
                <ReferenceLine
                    y={70}
                    label={{
                        position: "right",
                        value: "Low",
                        fill: "#fbbf24",
                        fontSize: 10,
                    }}
                    stroke="#fbbf24"
                    strokeDasharray="3 3"
                />
                <Area
                    type="monotone"
                    dataKey="value"
                    name="Glucose (mg/dL)"
                    stroke="#a855f7"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#colorGlucose)"
                    dot={{ r: 2, fill: "#a855f7" }}
                />
            </AreaChart>
        </ResponsiveContainer>
    </div>
);

export const BodyWeightChart: React.FC<{ data: TimeSeriesData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <h3 className="text-lg font-semibold text-slate-200 mb-2">
            Body Weight
        </h3>
        <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
                <defs>
                    <linearGradient
                        id="colorWeight"
                        x1="0"
                        y1="0"
                        x2="0"
                        y2="1"
                    >
                        <stop
                            offset="5%"
                            stopColor="#06b6d4"
                            stopOpacity={0.3}
                        />
                        <stop
                            offset="95%"
                            stopColor="#06b6d4"
                            stopOpacity={0}
                        />
                    </linearGradient>
                </defs>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="time"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    unit=" lbs"
                    domain={["dataMin - 5", "dataMax + 5"]}
                    tickFormatter={(tick) => formatNumber(tick)}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) => `${formatNumber(value)} lbs`}
                />
                <Area
                    type="monotone"
                    dataKey="weight"
                    name="Weight"
                    stroke="#06b6d4"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#colorWeight)"
                />
            </AreaChart>
        </ResponsiveContainer>
    </div>
);

export const BodyFatPercentageChart: React.FC<{ data: TimeSeriesData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <h3 className="text-lg font-semibold text-slate-200 mb-2">
            Body Fat %
        </h3>
        <ResponsiveContainer width="100%" height="100%">
            <LineChart data={data}>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="time"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    unit="%"
                    domain={["dataMin - 2", "dataMax + 2"]}
                    tickFormatter={(tick) => formatNumber(tick)}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) => `${formatNumber(value)} %`}
                />
                <Line
                    type="monotone"
                    dataKey="body_fat"
                    name="Body Fat"
                    stroke="#8b5cf6"
                    strokeWidth={2}
                    dot={{ r: 2 }}
                />
            </LineChart>
        </ResponsiveContainer>
    </div>
);

export const SleepStagesChart: React.FC<{ data: SleepData[] }> = ({ data }) => (
    <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data}>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="date"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    label={{
                        value: "Hours",
                        angle: -90,
                        position: "insideLeft",
                        offset: 0,
                        fill: "#64748b",
                        fontSize: 10,
                    }}
                    tickFormatter={(tick) => formatNumber(tick)}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) => `${formatNumber(value)} hrs`}
                />
                <Legend verticalAlign="top" iconType="circle" />
                <Bar
                    dataKey="deepSleep"
                    name="Deep"
                    stackId="a"
                    fill="#4338ca"
                />
                <Bar dataKey="remSleep" name="REM" stackId="a" fill="#8b5cf6" />
                <Bar
                    dataKey="lightSleep"
                    name="Light"
                    stackId="a"
                    fill="#3b82f6"
                />
                <Bar dataKey="awake" name="Awake" stackId="a" fill="#64748b" />
            </BarChart>
        </ResponsiveContainer>
    </div>
);

export const MacroTrendsChart: React.FC<{ data: DietaryData[] }> = ({
    data,
}) => (
    <div className="h-64 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={data}>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="date"
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    tickFormatter={(tick) => formatNumber(tick)}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) => `${formatNumber(value)} g`}
                />
                <Legend verticalAlign="top" iconType="circle" />
                <Area
                    type="monotone"
                    dataKey="protein"
                    stackId="1"
                    stroke="#f43f5e"
                    fill="#f43f5e"
                    fillOpacity={0.6}
                />
                <Area
                    type="monotone"
                    dataKey="carbs"
                    stackId="1"
                    stroke="#fbbf24"
                    fill="#fbbf24"
                    fillOpacity={0.6}
                />
                <Area
                    type="monotone"
                    dataKey="fat"
                    stackId="1"
                    stroke="#10b981"
                    fill="#10b981"
                    fillOpacity={0.6}
                />
            </AreaChart>
        </ResponsiveContainer>
    </div>
);

export const MonthlyCalorieChart: React.FC<{
    data: (DietaryData & { trend?: number })[];
}> = ({ data }) => (
    <div className="h-72 w-full">
        <ResponsiveContainer width="100%" height="100%">
            <ComposedChart data={data}>
                <CartesianGrid
                    strokeDasharray="3 3"
                    vertical={false}
                    stroke="#1e293b"
                />
                <XAxis
                    dataKey="date"
                    stroke="#64748b"
                    fontSize={9}
                    tickLine={false}
                    axisLine={false}
                    minTickGap={10}
                />
                <YAxis
                    stroke="#64748b"
                    fontSize={10}
                    tickLine={false}
                    axisLine={false}
                    domain={["dataMin - 500", "dataMax + 200"]}
                />
                <Tooltip
                    contentStyle={{
                        backgroundColor: "#0f172a",
                        border: "1px solid #334155",
                        borderRadius: "8px",
                    }}
                    formatter={(value: number) => `${formatNumber(value)} kcal`}
                />
                <Legend verticalAlign="top" />
                <Bar
                    dataKey="calories"
                    name="Daily Intake"
                    fill="#eab308"
                    fillOpacity={0.4}
                    radius={[2, 2, 0, 0]}
                />
                <Line
                    type="monotone"
                    dataKey="trend"
                    name="7D Trendline"
                    stroke="#f59e0b"
                    strokeWidth={3}
                    dot={false}
                    animationDuration={2000}
                />
            </ComposedChart>
        </ResponsiveContainer>
    </div>
);
