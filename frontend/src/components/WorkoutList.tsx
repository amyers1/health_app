
import React from 'react';
import { Workout } from '../types';

interface WorkoutListProps {
  workouts: Workout[];
}

const WorkoutList: React.FC<WorkoutListProps> = ({ workouts }) => {
  return (
    <div className="bg-slate-900/60 border border-slate-800 rounded-2xl overflow-hidden">
      <div className="p-4 border-b border-slate-800 flex justify-between items-center">
        <h3 className="font-semibold text-slate-100">Recent Workouts</h3>
        <button className="text-blue-400 text-sm font-medium hover:text-blue-300">View All</button>
      </div>
      <div className="divide-y divide-slate-800">
        {workouts.map((workout) => (
          <div key={workout.id} className="p-4 flex items-center justify-between hover:bg-slate-800/40 transition-colors">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-full bg-blue-500/10 flex items-center justify-center text-blue-400">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M18 14V8a6 6 0 0 0-12 0v6"/><path d="M15 14v4a2 2 0 1 1-4 0v-4"/><path d="M10 9a3 3 0 0 1 6 0"/><path d="M6 14h12"/></svg>
              </div>
              <div>
                <h4 className="text-sm font-medium text-slate-100">{workout.name}</h4>
                <p className="text-xs text-slate-500">{workout.time}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-sm font-semibold text-slate-200">{workout.calories} kcal</p>
              <p className="text-xs text-slate-500">{workout.duration} min</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default WorkoutList;
