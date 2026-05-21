"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthContext } from "@/components/AuthProvider";
import { hasLoggedInHint } from "@/lib/auth";
import { Skeleton } from "@/components/ui/skeleton";

/**
 * Root page — acts purely as a router.
 * Instantly redirects to /dashboard (authenticated) or /welcome (visitor).
 * Never renders the landing page or dashboard itself, eliminating flash.
 */
export default function RootPage() {
  const { isReady, isAuthenticated } = useAuthContext();
  const router = useRouter();

  useEffect(() => {
    if (!isReady) return;
    if (isAuthenticated) {
      router.replace("/dashboard");
    } else {
      router.replace("/welcome");
    }
  }, [isReady, isAuthenticated, router]);

  // While determining auth state, show a minimal loading indicator.
  // If we have a hint the user was logged in, show dashboard-like skeleton.
  // If no hint, show nothing (instant redirect to /welcome will happen).
  if (!isReady && !hasLoggedInHint()) {
    return null;
  }

  return (
    <div className="mx-auto max-w-geist-page px-6 py-8">
      <div className="mb-8 p-6 rounded-geist-md bg-geist-canvas dark:bg-neutral-900">
        <div className="flex items-center gap-3">
          <Skeleton className="h-10 w-10 rounded-geist-sm" />
          <div className="space-y-2">
            <Skeleton className="h-5 w-48" />
            <Skeleton className="h-4 w-64" />
          </div>
        </div>
      </div>
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mb-8">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-20 rounded-geist-md" />
        ))}
      </div>
    </div>
  );
}
