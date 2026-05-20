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
  textColor: string;
  check: (stats: UserStats) => boolean;
}

const BADGES: BadgeConfig[] = [
  {
    id: "first-commit",
    label: "First Commit",
    description: "Made your first commit on a showcase repository",
    icon: GitBranch,
    color: "text-green-600 dark:text-green-400",
    bgColor: "bg-green-50 border-green-200 dark:bg-green-950/50 dark:border-green-800",
    textColor: "text-green-800 dark:text-green-200",
    check: (stats) => stats.total_commits >= 1,
  },
  {
    id: "on-fire",
    label: "On Fire",
    description: "Maintained a 7-day contribution streak",
    icon: Flame,
    color: "text-orange-600 dark:text-orange-400",
    bgColor: "bg-orange-50 border-orange-200 dark:bg-orange-950/50 dark:border-orange-800",
    textColor: "text-orange-800 dark:text-orange-200",
    check: (stats) => stats.current_streak >= 7,
  },
  {
    id: "prolific",
    label: "Prolific",
    description: "Reached 50 total commits across all repositories",
    icon: Star,
    color: "text-yellow-600 dark:text-yellow-400",
    bgColor: "bg-yellow-50 border-yellow-200 dark:bg-yellow-950/50 dark:border-yellow-800",
    textColor: "text-yellow-800 dark:text-yellow-200",
    check: (stats) => stats.total_commits >= 50,
  },
  {
    id: "centurion",
    label: "Centurion",
    description: "Achieved 100 total commits — true dedication",
    icon: Trophy,
    color: "text-purple-600 dark:text-purple-400",
    bgColor: "bg-purple-50 border-purple-200 dark:bg-purple-950/50 dark:border-purple-800",
    textColor: "text-purple-800 dark:text-purple-200",
    check: (stats) => stats.total_commits >= 100,
  },
  {
    id: "polyglot",
    label: "Polyglot",
    description: "Used 3 or more programming languages in your projects",
    icon: Zap,
    color: "text-blue-600 dark:text-blue-400",
    bgColor: "bg-blue-50 border-blue-200 dark:bg-blue-950/50 dark:border-blue-800",
    textColor: "text-blue-800 dark:text-blue-200",
    check: (stats) => stats.languages.length >= 3,
  },
  {
    id: "repo-master",
    label: "Repo Master",
    description: "Showcasing 5 or more repositories on your profile",
    icon: Rocket,
    color: "text-indigo-600 dark:text-indigo-400",
    bgColor: "bg-indigo-50 border-indigo-200 dark:bg-indigo-950/50 dark:border-indigo-800",
    textColor: "text-indigo-800 dark:text-indigo-200",
    check: (stats) => stats.total_repos >= 5,
  },
  {
    id: "dedicated",
    label: "Dedicated",
    description: "Active for 30 or more days — consistency is key",
    icon: Target,
    color: "text-teal-600 dark:text-teal-400",
    bgColor: "bg-teal-50 border-teal-200 dark:bg-teal-950/50 dark:border-teal-800",
    textColor: "text-teal-800 dark:text-teal-200",
    check: (stats) => stats.active_days >= 30,
  },
  {
    id: "legend",
    label: "Legend",
    description: "Reached 500+ commits — a true open source legend",
    icon: Award,
    color: "text-rose-600 dark:text-rose-400",
    bgColor: "bg-rose-50 border-rose-200 dark:bg-rose-950/50 dark:border-rose-800",
    textColor: "text-rose-800 dark:text-rose-200",
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
        <div className="h-12 w-12 rounded-full bg-gray-50 dark:bg-neutral-800 flex items-center justify-center mx-auto mb-2">
          <Award className="h-6 w-6 text-gray-300 dark:text-gray-500" />
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400">
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
        <div className="h-12 w-12 rounded-full bg-gray-50 dark:bg-neutral-800 flex items-center justify-center mx-auto mb-2">
          <Award className="h-6 w-6 text-gray-300 dark:text-gray-500" />
        </div>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Keep contributing to unlock your first badge!
        </p>
        {/* Show locked badges as preview */}
        <div className="flex flex-wrap justify-center gap-2 mt-4">
          {locked.slice(0, 4).map((badge) => {
            const Icon = badge.icon;
            return (
              <div
                key={badge.id}
                className="flex items-center gap-1.5 px-2 py-1 rounded-lg bg-gray-100 dark:bg-neutral-800 border border-gray-200 dark:border-neutral-700 opacity-50"
                title={badge.description}
              >
                <Icon className="h-3.5 w-3.5 text-gray-400 dark:text-gray-500" />
                <span className="text-xs text-gray-500 dark:text-gray-400">{badge.label}</span>
              </div>
            );
          })}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {/* Earned badges — with visible description */}
      <div className="flex flex-col gap-2">
        {earned.map((badge) => {
          const Icon = badge.icon;
          return (
            <div
              key={badge.id}
              className={`flex items-start gap-3 px-3 py-2.5 rounded-xl border ${badge.bgColor} transition-transform hover:scale-[1.02]`}
            >
              <div className="flex-shrink-0 mt-0.5">
                <Icon className={`h-4 w-4 ${badge.color}`} />
              </div>
              <div className="flex flex-col min-w-0">
                <span className={`text-xs font-semibold ${badge.textColor}`}>
                  {badge.label}
                </span>
                <span className="text-[11px] text-gray-600 dark:text-gray-400 leading-tight mt-0.5">
                  {badge.description}
                </span>
              </div>
            </div>
          );
        })}
      </div>

      {/* Locked badges hint */}
      {locked.length > 0 && (
        <div className="flex items-center gap-2 pt-2 border-t border-gray-100 dark:border-neutral-800">
          <span className="text-[11px] text-gray-500 dark:text-gray-400">
            {locked.length} more to unlock
          </span>
          <div className="flex gap-1">
            {locked.slice(0, 3).map((badge) => {
              const Icon = badge.icon;
              return (
                <div
                  key={badge.id}
                  className="h-5 w-5 rounded bg-gray-100 dark:bg-neutral-800 flex items-center justify-center opacity-50"
                  title={`${badge.label}: ${badge.description}`}
                >
                  <Icon className="h-3 w-3 text-gray-400 dark:text-gray-500" />
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
