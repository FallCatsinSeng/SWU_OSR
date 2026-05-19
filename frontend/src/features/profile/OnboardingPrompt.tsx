"use client";

import Link from "next/link";
import { useCurrentUser } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";

export function OnboardingPrompt() {
  const { data: user } = useCurrentUser();

  if (!user) return null;
  if (user.alias !== user.github_username && user.bio !== "") return null;

  return (
    <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 mb-6">
      <p className="text-sm text-blue-800">
        Welcome! Customize your profile with a unique alias and bio.
      </p>
      <Link href="/settings" className="inline-block mt-2">
        <Button variant="outline" size="sm">
          Set up your profile
        </Button>
      </Link>
    </div>
  );
}
