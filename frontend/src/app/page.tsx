"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { useCurrentUser } from "@/hooks/useAuth";
import { ActivityFeed } from "@/features/feed/ActivityFeed";
import { OnboardingPrompt } from "@/features/profile/OnboardingPrompt";
import { Avatar } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import api from "@/lib/api";
import {
  Code2,
  GitBranch,
  Users,
  FolderGit2,
  MessageSquare,
  Trophy,
  ArrowRight,
  Zap,
  Shield,
  Globe,
  TrendingUp,
  ExternalLink,
  Activity,
} from "lucide-react";

// --- Types ---
interface CommunityStats {
  total_members: number;
  total_repos: number;
  total_activities: number;
  active_today: number;
  top_languages: string[];
  commits_this_week: number;
}

interface PopularRepo {
  id: string;
  repo_name: string;
  repo_full_name: string;
  description: string;
  language: string;
  html_url: string;
  academic_tag: string;
  owner_alias: string;
  owner_avatar: string;
  activity_count: number;
}

// --- Landing Page (visitors) ---
function HeroSection() {
  return (
    <section className="relative overflow-hidden">
      <div className="absolute inset-0 gradient-primary-subtle" />
      <div className="absolute top-0 right-0 w-96 h-96 bg-secondary-100/30 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2" />
      <div className="absolute bottom-0 left-0 w-64 h-64 bg-primary-100/40 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2" />

      <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-20 sm:py-28">
        <div className="text-center max-w-3xl mx-auto">
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/80 border border-primary-200 shadow-sm mb-6 animate-fade-in">
            <div className="h-2 w-2 rounded-full bg-secondary-500 animate-pulse-soft" />
            <span className="text-xs font-medium text-primary-700">
              Open Source Platform for Students
            </span>
          </div>

          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-gray-900 leading-tight animate-slide-up">
            Showcase Your Code,{" "}
            <span className="text-gradient">Build Your Future</span>
          </h1>

          <p className="mt-6 text-lg sm:text-xl text-gray-600 max-w-2xl mx-auto animate-slide-up">
            Platform mahasiswa STMIK Widya Utama untuk menampilkan proyek open source,
            berkolaborasi dengan peers, dan membangun portofolio profesional.
          </p>

          <div className="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4 animate-slide-up">
            <Link href="/login">
              <Button
                size="lg"
                className="gradient-primary text-white border-0 shadow-lg hover:shadow-xl transition-all px-8"
              >
                Get Started
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
            <Link href="/members">
              <Button variant="outline" size="lg" className="px-8">
                <Users className="mr-2 h-4 w-4" />
                Browse Members
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

function FeaturesSection() {
  const features = [
    { icon: FolderGit2, title: "Project Showcase", description: "Tampilkan repository GitHub terbaik kamu dengan tag akademik.", color: "bg-primary-50 text-primary-600" },
    { icon: Zap, title: "Activity Tracking", description: "Webhook otomatis mencatat setiap push, PR, dan release.", color: "bg-secondary-50 text-secondary-600" },
    { icon: MessageSquare, title: "Discussions", description: "Forum per-repository untuk review kode dan kolaborasi.", color: "bg-purple-50 text-purple-600" },
    { icon: Shield, title: "Pseudonymous Identity", description: "Gunakan alias publik. Identitas hanya terlihat sesama pengguna.", color: "bg-orange-50 text-orange-600" },
    { icon: Trophy, title: "Badges & Streaks", description: "Dapatkan badges berdasarkan pencapaian dan konsistensi.", color: "bg-yellow-50 text-yellow-600" },
    { icon: Globe, title: "Public Portfolio", description: "Profil publik yang bisa dibagikan ke recruiter atau siapa saja.", color: "bg-green-50 text-green-600" },
  ];

  return (
    <section className="py-20 bg-white">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-gray-900">
            Everything You Need to <span className="text-gradient">Stand Out</span>
          </h2>
          <p className="mt-3 text-gray-600 max-w-xl mx-auto">
            Satu platform untuk semua kebutuhan portofolio open source mahasiswa.
          </p>
        </div>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => {
            const Icon = feature.icon;
            return (
              <div key={feature.title} className="group p-6 rounded-2xl border border-gray-100 bg-white hover:border-primary-200 hover:shadow-lg transition-all duration-300">
                <div className={`inline-flex p-3 rounded-xl ${feature.color} mb-4`}>
                  <Icon className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">{feature.title}</h3>
                <p className="text-sm text-gray-600 leading-relaxed">{feature.description}</p>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}

function CTASection() {
  return (
    <section className="py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="relative rounded-3xl overflow-hidden gradient-primary p-12 sm:p-16 text-center">
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/2" />
          <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-1/2 -translate-x-1/2" />
          <div className="relative">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">Ready to Build Your Portfolio?</h2>
            <p className="text-white/80 max-w-xl mx-auto mb-8">
              Mulai showcase karya open source kamu hari ini. Login dengan akun SIAKAD dan hubungkan GitHub.
            </p>
            <Link href="/login">
              <Button size="lg" className="bg-white text-primary-700 hover:bg-gray-50 shadow-lg px-8">
                Start Now — It&apos;s Free
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

// --- Community Dashboard (logged-in users) ---
function CommunityStatsBar({ stats }: { stats: CommunityStats | undefined }) {
  if (!stats) {
    return (
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-20 rounded-xl" />
        ))}
      </div>
    );
  }

  const items = [
    { icon: Users, label: "Members", value: stats.total_members, color: "text-primary-600 bg-primary-50" },
    { icon: FolderGit2, label: "Repositories", value: stats.total_repos, color: "text-secondary-600 bg-secondary-50" },
    { icon: Activity, label: "This Week", value: stats.commits_this_week, color: "text-green-600 bg-green-50" },
    { icon: Zap, label: "Active Today", value: stats.active_today, color: "text-orange-600 bg-orange-50" },
  ];

  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
      {items.map((item) => {
        const Icon = item.icon;
        return (
          <div key={item.label} className="glass-card rounded-xl p-4 hover-lift">
            <div className="flex items-center gap-3">
              <div className={`h-10 w-10 rounded-lg ${item.color} flex items-center justify-center`}>
                <Icon className="h-5 w-5" />
              </div>
              <div>
                <p className="text-2xl font-bold text-gray-900">{item.value}</p>
                <p className="text-xs text-gray-500">{item.label}</p>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

function PopularReposSection({ repos }: { repos: PopularRepo[] | undefined }) {
  if (!repos || repos.length === 0) {
    return null;
  }

  return (
    <section>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
          <TrendingUp className="h-5 w-5 text-secondary-600" />
          Popular Repositories
        </h2>
        <Link href="/showcase" className="text-xs text-primary-600 hover:underline">
          View all
        </Link>
      </div>
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {repos.map((repo) => (
          <a
            key={repo.id}
            href={repo.html_url}
            target="_blank"
            rel="noopener noreferrer"
            className="group"
          >
            <Card className="h-full hover:border-primary-200 hover:shadow-md transition-all duration-200">
              <CardContent className="p-4">
                <div className="flex items-start justify-between mb-2">
                  <div className="flex items-center gap-2 min-w-0">
                    <FolderGit2 className="h-4 w-4 text-primary-600 shrink-0" />
                    <span className="font-medium text-sm text-gray-900 truncate group-hover:text-primary-600 transition-colors">
                      {repo.repo_name}
                    </span>
                    <ExternalLink className="h-3 w-3 text-gray-300 opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
                  </div>
                  {repo.activity_count > 0 && (
                    <Badge variant="secondary" className="text-[10px] shrink-0 bg-green-50 text-green-700">
                      {repo.activity_count} activities
                    </Badge>
                  )}
                </div>
                <p className="text-xs text-gray-500 line-clamp-2 mb-3">
                  {repo.description || "No description"}
                </p>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    {repo.language && (
                      <span className="text-[10px] px-2 py-0.5 rounded-full bg-gray-100 text-gray-600">
                        {repo.language}
                      </span>
                    )}
                    <span className="text-[10px] px-2 py-0.5 rounded-full bg-primary-50 text-primary-700">
                      {repo.academic_tag.replace("_", " ")}
                    </span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Avatar
                      src={repo.owner_avatar}
                      alt={repo.owner_alias}
                      fallback={repo.owner_alias.charAt(0).toUpperCase()}
                      size="sm"
                      className="h-5 w-5"
                    />
                    <span className="text-[10px] text-gray-400">{repo.owner_alias}</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </a>
        ))}
      </div>
    </section>
  );
}

function TrendingLanguages({ languages }: { languages: string[] }) {
  if (!languages || languages.length === 0) return null;

  const colors = [
    "bg-blue-100 text-blue-700 border-blue-200",
    "bg-green-100 text-green-700 border-green-200",
    "bg-purple-100 text-purple-700 border-purple-200",
    "bg-orange-100 text-orange-700 border-orange-200",
    "bg-teal-100 text-teal-700 border-teal-200",
    "bg-pink-100 text-pink-700 border-pink-200",
    "bg-indigo-100 text-indigo-700 border-indigo-200",
    "bg-yellow-100 text-yellow-700 border-yellow-200",
  ];

  return (
    <div className="glass-card rounded-xl p-5">
      <h3 className="text-sm font-semibold text-gray-900 mb-3 flex items-center gap-2">
        <Code2 className="h-4 w-4 text-purple-600" />
        Trending Languages
      </h3>
      <div className="flex flex-wrap gap-2">
        {languages.map((lang, i) => (
          <span
            key={lang}
            className={`px-3 py-1 text-xs font-medium rounded-full border ${colors[i % colors.length]}`}
          >
            {lang}
          </span>
        ))}
      </div>
    </div>
  );
}

function ActiveMembersSection() {
  const { data: membersData } = useQuery<{ members: Array<{ id: string; alias: string; avatar_url: string; github_username: string }> }>({
    queryKey: ["membersPreview"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: { members: Array<{ id: string; alias: string; avatar_url: string; github_username: string }>; total: number } }>("/members");
      return data.data;
    },
  });

  const members = membersData?.members ?? [];
  if (members.length === 0) return null;

  return (
    <div className="glass-card rounded-xl p-5">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-gray-900 flex items-center gap-2">
          <Users className="h-4 w-4 text-indigo-600" />
          Community Members
        </h3>
        <Link href="/members" className="text-[10px] text-primary-600 hover:underline">
          View all
        </Link>
      </div>
      <div className="flex flex-wrap gap-2">
        {members.slice(0, 12).map((member) => (
          <Link key={member.id} href={`/profiles/${member.alias}`} title={member.alias}>
            <Avatar
              src={member.avatar_url}
              alt={member.alias}
              fallback={member.alias.charAt(0).toUpperCase()}
              size="sm"
              className="h-9 w-9 ring-2 ring-white hover:ring-primary-200 transition-all hover:scale-110"
            />
          </Link>
        ))}
        {members.length > 12 && (
          <Link href="/members" className="h-9 w-9 rounded-full bg-gray-100 flex items-center justify-center text-xs text-gray-500 hover:bg-primary-50 hover:text-primary-700 transition-colors">
            +{members.length - 12}
          </Link>
        )}
      </div>
    </div>
  );
}

// --- Main Page Component ---
export default function HomePage() {
  const { data: user } = useCurrentUser();

  const { data: stats } = useQuery<CommunityStats>({
    queryKey: ["communityStats"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: CommunityStats }>("/stats");
      return data.data;
    },
    enabled: !!user,
  });

  const { data: popularRepos } = useQuery<PopularRepo[]>({
    queryKey: ["popularRepos"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: PopularRepo[] }>("/repos/popular");
      return data.data;
    },
    enabled: !!user,
  });

  // Logged-in: Community Dashboard
  if (user) {
    return (
      <div className="mx-auto max-w-6xl px-4 py-8">
        {/* Welcome + Quick Actions */}
        <div className="mb-6 p-5 rounded-2xl gradient-primary-subtle border border-primary-100">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
            <div className="flex items-center gap-3 min-w-0">
              <div className="h-10 w-10 rounded-full gradient-primary flex items-center justify-center shrink-0">
                <Code2 className="h-5 w-5 text-white" />
              </div>
              <div className="min-w-0">
                <h1 className="text-lg font-bold text-gray-900 break-words">
                  Welcome back, {user.alias}!
                </h1>
                <p className="text-sm text-gray-600">
                  Here&apos;s what&apos;s happening in the community.
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Link href="/showcase">
                <Button size="sm" variant="outline" className="gap-1.5 text-xs">
                  <FolderGit2 className="h-3.5 w-3.5" />
                  My Showcase
                </Button>
              </Link>
              <Link href={`/profiles/${user.alias}`}>
                <Button size="sm" variant="outline" className="gap-1.5 text-xs">
                  <Trophy className="h-3.5 w-3.5" />
                  My Profile
                </Button>
              </Link>
            </div>
          </div>
        </div>

        <OnboardingPrompt />

        {/* Community Stats */}
        <div className="mb-8">
          <CommunityStatsBar stats={stats} />
        </div>

        {/* Main content grid: Feed + Sidebar */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left: Activity Feed (2/3) */}
          <div className="lg:col-span-2 space-y-6">
            {/* Popular Repos */}
            <PopularReposSection repos={popularRepos} />

            {/* Activity Feed */}
            <section>
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold text-gray-900 flex items-center gap-2">
                  <GitBranch className="h-5 w-5 text-primary-600" />
                  Recent Activity
                </h2>
              </div>
              <ActivityFeed />
            </section>
          </div>

          {/* Right sidebar (1/3) */}
          <div className="space-y-4">
            {/* Trending Languages */}
            {stats && stats.top_languages.length > 0 && (
              <TrendingLanguages languages={stats.top_languages} />
            )}

            {/* Active Members */}
            <ActiveMembersSection />

            {/* Quick Links */}
            <div className="glass-card rounded-xl p-5">
              <h3 className="text-sm font-semibold text-gray-900 mb-3">Quick Links</h3>
              <div className="space-y-2">
                <Link href="/showcase" className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors">
                  <FolderGit2 className="h-4 w-4" />
                  Manage Showcase
                </Link>
                <Link href="/members" className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors">
                  <Users className="h-4 w-4" />
                  Discover Members
                </Link>
                <Link href="/settings" className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors">
                  <Zap className="h-4 w-4" />
                  Settings
                </Link>
                <a
                  href="https://github.com/FallCatsinSeng/SWU_OSR"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-sm text-gray-600 hover:text-primary-600 transition-colors"
                >
                  <Globe className="h-4 w-4" />
                  Source Code
                  <ExternalLink className="h-3 w-3" />
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Not logged in: Marketing landing page
  return (
    <div>
      <HeroSection />
      <FeaturesSection />
      <CTASection />
    </div>
  );
}
