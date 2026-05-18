"use client";

import { LoginForm } from "@/features/auth/LoginForm";

export default function LoginPage() {
  return (
    <div className="flex min-h-[calc(100vh-4rem)] items-center justify-center px-4">
      <LoginForm />
    </div>
  );
}
