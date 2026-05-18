"use client";

import { ProfileEditForm } from "@/features/profile/ProfileEditForm";

export default function SettingsPage() {
  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Settings</h1>
      <ProfileEditForm />
    </div>
  );
}
