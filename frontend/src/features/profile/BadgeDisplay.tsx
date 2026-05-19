"use client";

import type { UserStats } from "@/types/user";
import {
  GitBranch,
  Flame,
  Star,
  Rocket,
  Trophy,
  Zap,
  Target,
  Award,
} from "lucide-react";

interface BadgeConfig {
  id: string;
  label: string;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
  color: string;
  bgColor: string;
  check: (stats: UserStats) => boolean;
}

const BADGES: BadgeConfig[] = [
  {
    id: "first-commit",
    label: "First Commit",
    description: "Made your first commit",
    icon: GitBranch,
    color: "text-green-600",
    bgColor: "bg-green-50 border-green-200",
    check: (stats) => stats.total_commits >= 1,
  },
  {
    id: "on-fire",
    label: "On Fire",
    description: "7-day contribution streak",
    icon: Flame,
    color: "text-orange-600",
    bgColor: "bg-orange-50 border-orange-200",
    check: (stats) => stats.current_streak >= 7,
  },
  {
    id: "prolific",
    label: "Prolific",
    description: "Reached 50 commits",
    icon: Star,
    color: "text-yellow-600",
    bgColor: "bg-yellow-50 border-yellow-200",
    check: (stats) => stats.total_commits >= 50,
  },
  {
    id: "centurion",
    label: "Centurion",
    description: "Reached 100 commits",
    icon: Trophy,
    color: "text-purple-600",
    bgColor: "bg-purple-50 border-purple-200",
    check: (stats) => stats.total_commits >= 100,
  },
  {
    id: "polyglot",
    label: "Polyglot",
    description: "Used 3+ programming languages",
    icon: Zap,
    color: "text-blue-600",
    bgColor: "bg-blue-50 border-blue-200",
    check: (stats) => stats.languages.length >= 3,
  },
  {
    id: "repo-master",
    label: "Repo Master",
    description: "Showcase 5+ repositories",
    icon: Rocket,
    color: "text-indigo-600",
    bgColor: "bg-indigo-50 border-indigo-200",
    check: (stats) => stats.total_repos >= 5,
  },
  {
    id: "dedicated",
    label: "Dedicated",
    description: "30+ active days",
    icon: Target,
    color: "text-teal-600",
    bgColor: "bg-teal-50 border-teal-200",
    check: (stats) => stats.active_days >= 30,
  },
  {
    id: "legend",
    label: "Legend",
    description: "500+ commits",
    icon: Award,
    color: "text-rose-600",
    bgColor: "bg-rose-50 border-rose-200",
    check: (stats) => stats.total_commits >= 500,
  },
];

interface BadgeDisplayProps {
  stats: UserStats | null;
}

export function BadgeDisplay({ stats }: BadgeDisplayProps) {
  if (!stats) {
    return (
      <div className="text-center py-6">
        <div className="h-12 w-12 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-2">
          <Award className="h-6 w-6 text-gray-300" />
        </div>
        <p className="text-sm text-gray-400">
          Start contributing to earn badges!
        </p>
      </div>
    );
  }

  const earned = BADGES.filter((badge) => badge.check(stats));
  const locked = BADGES.filter((badge) => !badge.check(stats));

  if (earned.length === 0) {
    return (
      <div className="text-center py-6">
        <div className="h-12 w-12 rounded-full bg-gray-50 flex items-center justify-center mx-auto mb-2">
          <Award className="h-6 w-6 text-gray-300" />
        </div>
        <p className="text-sm text-gray-400">
          Keep contributing to unlock your first badge!
        </p>
        {/* Show locked badges */}
        <div className="flex flex-wrap justify-center gap-2 mt-4">
          {locked.slice(0, 4).map((badge) => {
            const Icon = badge.icon;
            return (
              <div
                key={badge.id}
                className="flex items-center gap-1.5 px-2 py-1 rounded-lg bg-gray-50 border border-gray-100 opacity-40"
                title={badge.description}
              >
                <Icon className="h-3.5 w-3.5 text-gray-400" />
                <span className="text-xs text-gray-400">{badge.label}</span>
              </div>
            );
          })}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {/* Earned badges */}
      <div className="flex flex-wrap gap-2">
        {earned.map((badge) => {
          const Icon = badge.icon;
          return (
            <div
              key={badge.id}
              className={`flex items-center gap-2 px-3 py-2 rounded-xl border ${badge.bgColor} transition-transform hover:scale-105`}
              title={badge.description}
            >
              <Icon className={`h-4 w-4 ${badge.color}`} />
              <span className="text-xs font-medium text-gray-700">
                {badge.label}
              </span>
            </div>
          );
        })}
      </div>

      {/* Locked badges hint */}
      {locked.length > 0 && (
        <div className="flex items-center gap-2 pt-2">
          <span className="text-[10px] text-gray-400">
            {locked.length} more to unlock
          </span>
          <div className="flex gap-1">
            {locked.slice(0, 3).map((badge) => {
              const Icon = badge.icon;
              return (
                <div
                  key={badge.id}
                  className="h-5 w-5 rounded bg-gray-50 flex items-center justify-center opacity-40"
                  title={`${badge.label}: ${badge.description}`}
                >
                  <Icon className="h-3 w-3 text-gray-400" />
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
