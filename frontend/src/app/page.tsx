"use client";

import Link from "next/link";
import { useCurrentUser } from "@/hooks/useAuth";
import { ActivityFeed } from "@/features/feed/ActivityFeed";
import { OnboardingPrompt } from "@/features/profile/OnboardingPrompt";
import { Button } from "@/components/ui/button";
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
} from "lucide-react";

function HeroSection() {
  return (
    <section className="relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 gradient-primary-subtle" />
      <div className="absolute top-0 right-0 w-96 h-96 bg-secondary-100/30 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2" />
      <div className="absolute bottom-0 left-0 w-64 h-64 bg-primary-100/40 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2" />

      <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-20 sm:py-28">
        <div className="text-center max-w-3xl mx-auto">
          {/* Badge */}
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
            Platform mahasiswa UIN Walisongo untuk menampilkan proyek open source,
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

          {/* Stats row */}
          <div className="mt-16 grid grid-cols-2 sm:grid-cols-4 gap-4 max-w-2xl mx-auto animate-fade-in">
            <StatCard icon={Users} label="Members" value="50+" />
            <StatCard icon={FolderGit2} label="Repositories" value="120+" />
            <StatCard icon={GitBranch} label="Commits" value="2.5K+" />
            <StatCard icon={MessageSquare} label="Discussions" value="300+" />
          </div>
        </div>
      </div>
    </section>
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  value: string;
}) {
  return (
    <div className="glass-card rounded-xl p-4 text-center hover-lift">
      <Icon className="h-5 w-5 text-secondary-600 mx-auto mb-1" />
      <p className="text-2xl font-bold text-gray-900">{value}</p>
      <p className="text-xs text-gray-500">{label}</p>
    </div>
  );
}

function FeaturesSection() {
  const features = [
    {
      icon: FolderGit2,
      title: "Project Showcase",
      description:
        "Tampilkan repository GitHub terbaik kamu dengan tag akademik — coursework, thesis, hackathon.",
      color: "bg-primary-50 text-primary-600",
    },
    {
      icon: Zap,
      title: "Activity Tracking",
      description:
        "Webhook otomatis mencatat setiap push, PR, dan release. Lihat kontribusimu tumbuh.",
      color: "bg-secondary-50 text-secondary-600",
    },
    {
      icon: MessageSquare,
      title: "Discussions",
      description:
        "Forum per-repository untuk review kode, feedback, dan kolaborasi sesama mahasiswa.",
      color: "bg-purple-50 text-purple-600",
    },
    {
      icon: Shield,
      title: "Pseudonymous Identity",
      description:
        "Gunakan alias publik. Identitas akademik hanya terlihat oleh sesama pengguna terdaftar.",
      color: "bg-orange-50 text-orange-600",
    },
    {
      icon: Trophy,
      title: "Badges & Streaks",
      description:
        "Dapatkan badges berdasarkan pencapaian — first commit, 7-day streak, 100 commits.",
      color: "bg-yellow-50 text-yellow-600",
    },
    {
      icon: Globe,
      title: "Public Portfolio",
      description:
        "Profil publik yang bisa dibagikan ke recruiter, dosen, atau siapa saja.",
      color: "bg-green-50 text-green-600",
    },
  ];

  return (
    <section className="py-20 bg-white">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-gray-900">
            Everything You Need to{" "}
            <span className="text-gradient">Stand Out</span>
          </h2>
          <p className="mt-3 text-gray-600 max-w-xl mx-auto">
            Satu platform untuk semua kebutuhan portofolio open source mahasiswa.
          </p>
        </div>

        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => {
            const Icon = feature.icon;
            return (
              <div
                key={feature.title}
                className="group p-6 rounded-2xl border border-gray-100 bg-white hover:border-primary-200 hover:shadow-lg transition-all duration-300"
              >
                <div
                  className={`inline-flex p-3 rounded-xl ${feature.color} mb-4`}
                >
                  <Icon className="h-6 w-6" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  {feature.title}
                </h3>
                <p className="text-sm text-gray-600 leading-relaxed">
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
    <section className="py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="relative rounded-3xl overflow-hidden gradient-primary p-12 sm:p-16 text-center">
          {/* Decorative elements */}
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -translate-y-1/2 translate-x-1/2" />
          <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-1/2 -translate-x-1/2" />

          <div className="relative">
            <h2 className="text-3xl sm:text-4xl font-bold text-white mb-4">
              Ready to Build Your Portfolio?
            </h2>
            <p className="text-white/80 max-w-xl mx-auto mb-8">
              Mulai showcase karya open source kamu hari ini. Login dengan akun
              SIAKAD dan hubungkan GitHub.
            </p>
            <Link href="/login">
              <Button
                size="lg"
                className="bg-white text-primary-700 hover:bg-gray-50 shadow-lg px-8"
              >
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

export default function HomePage() {
  const { data: user } = useCurrentUser();

  // If user is logged in, show the activity feed dashboard
  if (user) {
    return (
      <div className="mx-auto max-w-4xl px-4 py-8">
        {/* Welcome banner */}
        <div className="mb-8 p-6 rounded-2xl gradient-primary-subtle border border-primary-100">
          <div className="flex items-center gap-3 mb-2">
            <div className="h-10 w-10 rounded-full gradient-primary flex items-center justify-center">
              <Code2 className="h-5 w-5 text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold text-gray-900">
                Welcome back, {user.alias}!
              </h1>
              <p className="text-sm text-gray-600">
                Here&apos;s what&apos;s happening in the community.
              </p>
            </div>
          </div>
        </div>

        <OnboardingPrompt />

        {/* Quick actions */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-8">
          <Link
            href="/showcase"
            className="glass-card rounded-xl p-4 text-center hover-lift group"
          >
            <FolderGit2 className="h-5 w-5 text-primary-600 mx-auto mb-2 group-hover:scale-110 transition-transform" />
            <span className="text-xs font-medium text-gray-700">
              My Showcase
            </span>
          </Link>
          <Link
            href={`/profiles/${user.alias}`}
            className="glass-card rounded-xl p-4 text-center hover-lift group"
          >
            <Trophy className="h-5 w-5 text-secondary-600 mx-auto mb-2 group-hover:scale-110 transition-transform" />
            <span className="text-xs font-medium text-gray-700">
              My Profile
            </span>
          </Link>
          <Link
            href="/members"
            className="glass-card rounded-xl p-4 text-center hover-lift group"
          >
            <Users className="h-5 w-5 text-purple-600 mx-auto mb-2 group-hover:scale-110 transition-transform" />
            <span className="text-xs font-medium text-gray-700">
              Members
            </span>
          </Link>
          <Link
            href="/settings"
            className="glass-card rounded-xl p-4 text-center hover-lift group"
          >
            <Zap className="h-5 w-5 text-orange-600 mx-auto mb-2 group-hover:scale-110 transition-transform" />
            <span className="text-xs font-medium text-gray-700">
              Settings
            </span>
          </Link>
        </div>

        {/* Activity Feed */}
        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900 flex items-center gap-2">
              <GitBranch className="h-5 w-5 text-primary-600" />
              Recent Activity
            </h2>
          </div>
          <ActivityFeed />
        </section>
      </div>
    );
  }

  // If not logged in, show the marketing/landing page
  return (
    <div>
      <HeroSection />
      <FeaturesSection />
      <CTASection />
    </div>
  );
}
