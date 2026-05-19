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
    <section className="relative overflow-hidden bg-geist-canvas dark:bg-black">
      {/* Mesh gradient — hero scale atmospheric backdrop */}
      <div className="absolute inset-0 geist-mesh-gradient" />

      <div className="relative mx-auto max-w-geist-page px-6 py-24 sm:py-32 lg:py-40">
        <div className="text-center max-w-3xl mx-auto">
          {/* Mono eyebrow badge */}
          <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-geist-full bg-geist-canvas-soft border border-geist-hairline dark:bg-neutral-900 dark:border-neutral-800 mb-8 animate-fade-in">
            <span className="text-caption-mono text-geist-body dark:text-neutral-300">
              Open source platform for students
            </span>
          </div>

          {/* Hero headline — display-xl, sentence-case, period-terminated */}
          <h1 className="text-display-xl text-geist-ink dark:text-white animate-slide-up">
            Showcase your code, build your future.
          </h1>

          {/* Lead paragraph — body-lg */}
          <p className="mt-6 text-body-lg text-geist-body dark:text-neutral-300 max-w-2xl mx-auto animate-slide-up">
            Platform mahasiswa STMIK Widya Utama untuk menampilkan proyek open
            source, berkolaborasi dengan peers, dan membangun portofolio
            profesional.
          </p>

          {/* CTA row — pill primary + pill secondary */}
          <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4 animate-slide-up">
            <Link href="/login">
              <Button variant="default" size="default">
                Get Started
                <ArrowRight className="ml-2 h-4 w-4" />
              </Button>
            </Link>
            <Link href="/members">
              <Button variant="secondary" size="default">
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
    {
      icon: FolderGit2,
      title: "Project showcase.",
      description:
        "Tampilkan repository GitHub terbaik kamu dengan tag akademik.",
    },
    {
      icon: Activity,
      title: "Activity tracking.",
      description:
        "Webhook otomatis mencatat setiap push, PR, dan release.",
    },
    {
      icon: MessageSquare,
      title: "Discussions.",
      description:
        "Forum per-repository untuk review kode dan kolaborasi.",
    },
    {
      icon: Shield,
      title: "Pseudonymous identity.",
      description:
        "Gunakan alias publik. Identitas hanya terlihat sesama pengguna.",
    },
    {
      icon: Trophy,
      title: "Badges & streaks.",
      description:
        "Dapatkan badges berdasarkan pencapaian dan konsistensi.",
    },
    {
      icon: Globe,
      title: "Public portfolio.",
      description:
        "Profil publik yang bisa dibagikan ke recruiter atau siapa saja.",
    },
  ];

  return (
    <section className="py-24 bg-geist-canvas-soft dark:bg-[#0A0A0A]">
      <div className="mx-auto max-w-geist-page px-6">
        <div className="text-center mb-16">
          {/* Mono eyebrow */}
          <p className="text-caption-mono text-geist-mute dark:text-white0 uppercase mb-4">
            Features
          </p>
          {/* Section headline — display-lg, sentence-case, period-terminated */}
          <h2 className="text-display-lg text-geist-ink dark:text-white">
            Everything you need to stand out.
          </h2>
          <p className="mt-4 text-body-lg text-geist-body dark:text-neutral-300 max-w-xl mx-auto">
            Satu platform untuk semua kebutuhan portofolio open source
            mahasiswa.
          </p>
        </div>

        {/* 3-up card grid */}
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => {
            const Icon = feature.icon;
            return (
              <div
                key={feature.title}
                className="p-6 rounded-geist-md bg-geist-canvas dark:bg-neutral-900 geist-level-3 transition-shadow hover:shadow-geist-4"
              >
                <div className="h-10 w-10 rounded-geist-sm bg-geist-canvas-soft-2 dark:bg-neutral-800 flex items-center justify-center mb-4">
                  <Icon className="h-5 w-5 text-geist-ink dark:text-white" />
                </div>
                {/* Card headline — display-sm, sentence-case */}
                <h3 className="text-display-sm text-geist-ink dark:text-white mb-2">
                  {feature.title}
                </h3>
                <p className="text-body-sm text-geist-body dark:text-neutral-300">
                  {feature.description}
                </p>
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
    <section className="py-24 bg-geist-canvas-soft dark:bg-[#0A0A0A]">
      <div className="mx-auto max-w-geist-page px-6">
        {/* Polarity-flipped dark band — in dark mode it flips to white */}
        <div className="rounded-geist-lg bg-geist-primary dark:bg-white p-12 sm:p-16 text-center">
          <h2 className="text-display-lg text-geist-on-primary dark:text-black mb-4">
            Ready to build your portfolio?
          </h2>
          <p className="text-body-lg text-geist-on-primary/70 dark:text-black/60 max-w-xl mx-auto mb-10">
            Mulai showcase karya open source kamu hari ini. Login dengan akun
            SIAKAD dan hubungkan GitHub.
          </p>
          <Link href="/login">
            {/* Inverted CTA — white pill on dark surface / black pill on white surface in dark mode */}
            <Button variant="secondary" size="default" className="dark:bg-black dark:text-white dark:hover:bg-neutral-900">
              Start now — it&apos;s free
              <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </Link>
        </div>
      </div>
    </section>
  );
}

// --- Community Dashboard (logged-in users) ---
function CommunityStatsBar({
  stats,
}: {
  stats: CommunityStats | undefined;
}) {
  if (!stats) {
    return (
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-20 rounded-geist-md" />
        ))}
      </div>
    );
  }

  const items = [
    { icon: Users, label: "Members", value: stats.total_members },
    { icon: FolderGit2, label: "Repositories", value: stats.total_repos },
    { icon: Activity, label: "This week", value: stats.commits_this_week },
    { icon: TrendingUp, label: "Active today", value: stats.active_today },
  ];

  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
      {items.map((item) => {
        const Icon = item.icon;
        return (
          <div
            key={item.label}
            className="rounded-geist-md bg-geist-canvas dark:bg-neutral-900 geist-level-2 p-4"
          >
            <div className="flex items-center gap-3">
              <div className="h-9 w-9 rounded-geist-sm bg-geist-canvas-soft-2 dark:bg-neutral-800 flex items-center justify-center">
                <Icon className="h-4 w-4 text-geist-ink dark:text-white" />
              </div>
              <div>
                <p className="text-display-sm text-geist-ink dark:text-white">{item.value}</p>
                <p className="text-caption text-geist-mute dark:text-white0">{item.label}</p>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

function PopularReposSection({
  repos,
}: {
  repos: PopularRepo[] | undefined;
}) {
  if (!repos || repos.length === 0) {
    return null;
  }

  return (
    <section>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-body-sm-strong text-geist-ink dark:text-white flex items-center gap-2">
          <TrendingUp className="h-4 w-4 text-geist-mute dark:text-white0" />
          Popular repositories
        </h2>
        <Link
          href="/showcase"
          className="text-caption text-geist-link hover:text-geist-link-deep transition-colors dark:text-neutral-300 dark:hover:text-white"
        >
          View all
        </Link>
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {repos.map((repo) => (
          <a
            key={repo.id}
            href={repo.html_url}
            target="_blank"
            rel="noopener noreferrer"
            className="group"
          >
            <Card className="h-full transition-shadow hover:shadow-geist-3">
              <CardContent className="p-4">
                <div className="flex items-start justify-between mb-2">
                  <div className="flex items-center gap-2 min-w-0">
                    <FolderGit2 className="h-4 w-4 text-geist-mute dark:text-white0 shrink-0" />
                    <span className="text-body-sm-strong text-geist-ink dark:text-white truncate group-hover:text-geist-link dark:group-hover:text-white transition-colors">
                      {repo.repo_name}
                    </span>
                    <ExternalLink className="h-3 w-3 text-geist-hairline-strong dark:text-neutral-300 opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
                  </div>
                  {repo.activity_count > 0 && (
                    <Badge variant="success" className="text-[10px] shrink-0">
                      {repo.activity_count} activities
                    </Badge>
                  )}
                </div>
                <p className="text-caption text-geist-mute dark:text-white0 line-clamp-2 mb-3">
                  {repo.description || "No description"}
                </p>
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    {repo.language && (
                      <Badge variant="secondary" className="text-[10px]">
                        {repo.language}
                      </Badge>
                    )}
                    <Badge variant="default" className="text-[10px]">
                      {repo.academic_tag.replace("_", " ")}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-1">
                    <Avatar
                      src={repo.owner_avatar}
                      alt={repo.owner_alias}
                      fallback={repo.owner_alias.charAt(0).toUpperCase()}
                      size="sm"
                      className="h-5 w-5"
                    />
                    <span className="text-caption text-geist-mute dark:text-white0">
                      {repo.owner_alias}
                    </span>
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

  return (
    <div className="rounded-geist-md bg-geist-canvas dark:bg-neutral-900 geist-level-2 p-5">
      <h3 className="text-body-sm-strong text-geist-ink dark:text-white mb-3 flex items-center gap-2">
        <Code2 className="h-4 w-4 text-geist-mute dark:text-white0" />
        Trending languages
      </h3>
      <div className="flex flex-wrap gap-2">
        {languages.map((lang) => (
          <Badge key={lang} variant="outline" className="text-caption">
            {lang}
          </Badge>
        ))}
      </div>
    </div>
  );
}

function ActiveMembersSection() {
  const { data: membersData } = useQuery<{
    members: Array<{
      id: string;
      alias: string;
      avatar_url: string;
      github_username: string;
    }>;
  }>({
    queryKey: ["membersPreview"],
    queryFn: async () => {
      const { data } = await api.get<{
        ok: boolean;
        data: {
          members: Array<{
            id: string;
            alias: string;
            avatar_url: string;
            github_username: string;
          }>;
          total: number;
        };
      }>("/members");
      return data.data;
    },
  });

  const members = membersData?.members ?? [];
  if (members.length === 0) return null;

  return (
    <div className="rounded-geist-md bg-geist-canvas dark:bg-neutral-900 geist-level-2 p-5">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-body-sm-strong text-geist-ink dark:text-white flex items-center gap-2">
          <Users className="h-4 w-4 text-geist-mute dark:text-white0" />
          Community members
        </h3>
        <Link
          href="/members"
          className="text-caption text-geist-link hover:text-geist-link-deep transition-colors dark:text-neutral-300 dark:hover:text-white"
        >
          View all
        </Link>
      </div>
      <div className="flex flex-wrap gap-2">
        {members.slice(0, 12).map((member) => (
          <Link
            key={member.id}
            href={`/profiles/${member.alias}`}
            title={member.alias}
          >
            <Avatar
              src={member.avatar_url}
              alt={member.alias}
              fallback={member.alias.charAt(0).toUpperCase()}
              size="sm"
              className="h-8 w-8 ring-2 ring-geist-canvas dark:ring-black hover:ring-geist-hairline dark:hover:ring-neutral-700 transition-all"
            />
          </Link>
        ))}
        {members.length > 12 && (
          <Link
            href="/members"
            className="h-8 w-8 rounded-full bg-geist-canvas-soft-2 dark:bg-neutral-800 flex items-center justify-center text-caption text-geist-mute dark:text-white0 hover:bg-geist-canvas-soft dark:hover:bg-neutral-700 hover:text-geist-ink dark:hover:text-white transition-colors"
          >
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
      const { data } = await api.get<{ ok: boolean; data: CommunityStats }>(
        "/stats"
      );
      return data.data;
    },
    enabled: !!user,
  });

  const { data: popularRepos } = useQuery<PopularRepo[]>({
    queryKey: ["popularRepos"],
    queryFn: async () => {
      const { data } = await api.get<{ ok: boolean; data: PopularRepo[] }>(
        "/repos/popular"
      );
      return data.data;
    },
    enabled: !!user,
  });

  // Logged-in: Community Dashboard
  if (user) {
    return (
      <div className="mx-auto max-w-geist-page px-6 py-8">
        {/* Welcome header */}
        <div className="mb-8 p-6 rounded-geist-md bg-geist-canvas dark:bg-neutral-900 geist-level-2">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div className="flex items-center gap-3 min-w-0">
              <div className="h-10 w-10 rounded-geist-sm bg-geist-primary dark:bg-white flex items-center justify-center shrink-0">
                <Code2 className="h-5 w-5 text-geist-on-primary dark:text-black" />
              </div>
              <div className="min-w-0">
                <h1 className="text-display-sm text-geist-ink dark:text-white break-words">
                  Welcome back, {user.alias}.
                </h1>
                <p className="text-body-sm text-geist-body dark:text-neutral-300">
                  Here&apos;s what&apos;s happening in the community.
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Link href="/showcase">
                <Button variant="outline" size="sm">
                  <FolderGit2 className="mr-1.5 h-3.5 w-3.5" />
                  My Showcase
                </Button>
              </Link>
              <Link href={`/profiles/${user.alias}`}>
                <Button variant="outline" size="sm">
                  <Trophy className="mr-1.5 h-3.5 w-3.5" />
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
            <PopularReposSection repos={popularRepos} />

            <section>
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-body-sm-strong text-geist-ink dark:text-white flex items-center gap-2">
                  <GitBranch className="h-4 w-4 text-geist-mute dark:text-white0" />
                  Recent activity
                </h2>
              </div>
              <ActivityFeed />
            </section>
          </div>

          {/* Right sidebar (1/3) */}
          <div className="space-y-4">
            {stats && stats.top_languages.length > 0 && (
              <TrendingLanguages languages={stats.top_languages} />
            )}

            <ActiveMembersSection />

            {/* Quick Links */}
            <div className="rounded-geist-md bg-geist-canvas dark:bg-neutral-900 geist-level-2 p-5">
              <h3 className="text-body-sm-strong text-geist-ink dark:text-white mb-3">
                Quick links
              </h3>
              <div className="space-y-2">
                <Link
                  href="/showcase"
                  className="flex items-center gap-2 text-body-sm text-geist-body hover:text-geist-ink transition-colors dark:text-neutral-300 dark:hover:text-white"
                >
                  <FolderGit2 className="h-4 w-4" />
                  Manage Showcase
                </Link>
                <Link
                  href="/members"
                  className="flex items-center gap-2 text-body-sm text-geist-body hover:text-geist-ink transition-colors dark:text-neutral-300 dark:hover:text-white"
                >
                  <Users className="h-4 w-4" />
                  Discover Members
                </Link>
                <Link
                  href="/settings"
                  className="flex items-center gap-2 text-body-sm text-geist-body hover:text-geist-ink transition-colors dark:text-neutral-300 dark:hover:text-white"
                >
                  <Activity className="h-4 w-4" />
                  Settings
                </Link>
                <a
                  href="https://github.com/FallCatsinSeng/SWU_OSR"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-body-sm text-geist-body hover:text-geist-ink transition-colors dark:text-neutral-300 dark:hover:text-white"
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
