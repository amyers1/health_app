import React from "react";
import { formatNumber } from "../utils";

interface MetricCardProps {
    title: string;
    value: string | number;
    unit: string;
    icon: React.ReactNode;
    colorClass: string;
    subtitle?: string;
    trend?: {
        value: number;
        isUp: boolean;
    };
}

const MetricCard: React.FC<MetricCardProps> = ({
    title,
    value,
    unit,
    icon,
    colorClass,
    subtitle,
    trend,
}) => {
    return (
        <div className="bg-slate-900/60 border border-slate-800 rounded-2xl p-4 flex flex-col justify-between hover:border-slate-700 transition-colors">
            <div className="flex justify-between items-start mb-2">
                <div className={`p-2 rounded-xl bg-opacity-10 ${colorClass}`}>
                    {icon}
                </div>
                {trend && (
                    <span
                        className={`text-xs font-medium ${trend.isUp ? "text-emerald-400" : "text-rose-400"}`}
                    >
                        {trend.isUp ? "↑" : "↓"} {formatNumber(trend.value)}%
                    </span>
                )}
            </div>
            <div>
                <h3 className="text-slate-400 text-xs font-medium uppercase tracking-wider mb-1">
                    {title}
                </h3>
                <div className="flex items-baseline gap-1">
                    <span className="text-2xl font-bold text-slate-50">
                        {formatNumber(value)}
                    </span>
                    <span className="text-slate-500 text-sm">{unit}</span>
                </div>
                {subtitle && (
                    <p className="text-slate-500 text-[10px] mt-1 truncate">
                        {subtitle}
                    </p>
                )}
            </div>
        </div>
    );
};

export default MetricCard;
