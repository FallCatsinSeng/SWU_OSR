"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useCurrentUser } from "@/hooks/useAuth";
import { LeaderboardTable } from "@/features/leaderboard/LeaderboardTable";
import { PeriodToggle } from "@/features/leaderboard/PeriodToggle";
import { UserPointsCard } from "@/features/leaderboard/UserPointsCard";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import api from "@/lib/api";
import type {
  LeaderboardPeriod,
  LeaderboardResult,
  UserPointsSummary,
} from "@/types/leaderboard";
import { Trophy, Info, GitBranch, GitPullRequest, MessageSquare, FolderGit2, Flame } from "lucide-react";

export default function LeaderboardPage() {
  const [period, setPeriod] = useState<LeaderboardPeriod>("weekly");
  const { data: user } = useCurrentUser();

  const {
    data: leaderboard,
    isLoading,
    isError,
  } = useQuery<LeaderboardResult>({
    queryKey: ["leaderboard", period],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: LeaderboardResult }>(
        "/leaderboard",
        { params: { period, limit: 50, offset: 0 } }
      );
      return data.data;
    },
  });

  const { data: myPoints } = useQuery<UserPointsSummary>({
    queryKey: ["leaderboard", "me", period],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: UserPointsSummary }>(
        "/leaderboard/me",
        { params: { period } }
      );
      return data.data;
    },
    enabled: !!user,
  });

  return (
    <div className="mx-auto max-w-geist-page px-6 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
          <div>
            <h1 className="text-display-lg text-geist-ink dark:text-white flex items-center gap-3">
              <Trophy className="h-7 w-7 text-amber-500" />
              Leaderboard
            </h1>
            <p className="text-body-sm text-geist-body dark:text-neutral-400 mt-1">
              Top contributors ranked by activity points.
            </p>
          </div>
          <PeriodToggle period={period} onChange={setPeriod} />
        </div>

        {/* My points card (authenticated users only) */}
        {user && myPoints && <UserPointsCard summary={myPoints} />}
      </div>

      {/* Point system info */}
      <Card className="mb-6">
        <CardContent className="p-4">
          <div className="flex items-start gap-3">
            <Info className="h-4 w-4 text-geist-mute dark:text-neutral-500 mt-0.5 shrink-0" />
            <div>
              <p className="text-caption text-geist-body dark:text-neutral-400 mb-2">
                <span className="font-medium text-geist-ink dark:text-white">How points work:</span>
              </p>
              <div className="flex flex-wrap gap-3 text-caption text-geist-mute dark:text-neutral-500">
                <span className="inline-flex items-center gap-1">
                  <GitBranch className="h-3 w-3" /> Push: 3 pts
                </span>
                <span className="inline-flex items-center gap-1">
                  <GitPullRequest className="h-3 w-3" /> PR: 5 pts
                </span>
                <span className="inline-flex items-center gap-1">
                  <FolderGit2 className="h-3 w-3" /> Showcase: 10 pts
                </span>
                <span className="inline-flex items-center gap-1">
                  <MessageSquare className="h-3 w-3" /> Thread: 2 pts
                </span>
                <span className="inline-flex items-center gap-1">
                  <MessageSquare className="h-3 w-3" /> Comment: 1 pt
                </span>
                <span className="inline-flex items-center gap-1">
                  <Flame className="h-3 w-3" /> 7-day streak: +15 pts
                </span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Leaderboard table */}
      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 10 }).map((_, i) => (
            <Skeleton key={i} className="h-16 w-full rounded-geist-md" />
          ))}
        </div>
      ) : isError ? (
        <Card>
          <CardContent className="p-8 text-center">
            <p className="text-body-sm text-geist-body dark:text-neutral-400">
              Failed to load leaderboard. Please try again later.
            </p>
          </CardContent>
        </Card>
      ) : (
        <LeaderboardTable
          entries={leaderboard?.entries ?? []}
          currentUserId={user?.id}
        />
      )}
    </div>
  );
}
