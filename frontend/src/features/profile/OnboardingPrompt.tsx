"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useCurrentUser } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { X } from "lucide-react";

const DISMISS_KEY = "swu_osr_onboarding_dismissed";

export function OnboardingPrompt() {
  const { data: user } = useCurrentUser();
  const [dismissed, setDismissed] = useState(true);

  useEffect(() => {
    setDismissed(localStorage.getItem(DISMISS_KEY) === "true");
  }, []);

  if (!user) return null;
  if (dismissed) return null;
  if (user.alias !== user.github_username && user.bio !== "") return null;

  const handleDismiss = () => {
    localStorage.setItem(DISMISS_KEY, "true");
    setDismissed(true);
  };

  return (
    <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 mb-6 relative">
      <button
        onClick={handleDismiss}
        className="absolute top-2 right-2 text-blue-400 hover:text-blue-600"
        aria-label="Dismiss"
      >
        <X className="h-4 w-4" />
      </button>
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
