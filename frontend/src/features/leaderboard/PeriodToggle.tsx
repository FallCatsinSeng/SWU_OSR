"use client";

import { cn } from "@/lib/utils";
import type { LeaderboardPeriod } from "@/types/leaderboard";

interface PeriodToggleProps {
  period: LeaderboardPeriod;
  onChange: (period: LeaderboardPeriod) => void;
}

const periods: { value: LeaderboardPeriod; label: string }[] = [
  { value: "weekly", label: "This Week" },
  { value: "monthly", label: "This Month" },
  { value: "semester", label: "Semester" },
  { value: "all_time", label: "All Time" },
];

export function PeriodToggle({ period, onChange }: PeriodToggleProps) {
  return (
    <div className="inline-flex items-center gap-1 p-1 rounded-geist-md bg-geist-canvas-soft dark:bg-neutral-900 border border-geist-hairline dark:border-neutral-800">
      {periods.map((p) => (
        <button
          key={p.value}
          onClick={() => onChange(p.value)}
          className={cn(
            "px-3 py-1.5 rounded-geist-sm text-caption font-medium transition-all",
            period === p.value
              ? "bg-geist-canvas dark:bg-neutral-800 text-geist-ink dark:text-white shadow-geist-1"
              : "text-geist-mute dark:text-neutral-500 hover:text-geist-body dark:hover:text-neutral-300"
          )}
        >
          {p.label}
        </button>
      ))}
    </div>
  );
}
