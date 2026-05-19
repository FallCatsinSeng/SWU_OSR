"use client";

import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { ProfileEditForm } from "@/features/profile/ProfileEditForm";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { ExternalLink, CheckCircle } from "lucide-react";

export default function SettingsPage() {
  const { data: user } = useCurrentUser();
  const logout = useLogout();

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Settings</h1>

      <div className="space-y-6">
        <section>
          <h2 className="text-lg font-semibold text-gray-900 mb-3">Profile</h2>
          <ProfileEditForm />
        </section>

        {user && (
          <section>
            <h2 className="text-lg font-semibold text-gray-900 mb-3">
              GitHub Connection
            </h2>
            <Card>
              <CardContent className="p-4">
                <div className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <span className="text-sm text-gray-700">Connected as</span>
                  <a
                    href={`https://github.com/${user.github_username}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-sm text-blue-600 hover:underline flex items-center gap-1"
                  >
                    @{user.github_username}
                    <ExternalLink className="h-3 w-3" />
                  </a>
                </div>
              </CardContent>
            </Card>
          </section>
        )}

        {user && (
          <section>
            <h2 className="text-lg font-semibold text-gray-900 mb-3">
              Linked Identity
            </h2>
            <Card>
              <CardContent className="p-4">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-gray-700">NIM:</span>
                  <Badge variant="secondary">{user.nim}</Badge>
                </div>
              </CardContent>
            </Card>
          </section>
        )}

        <section>
          <h2 className="text-lg font-semibold text-gray-900 mb-3">Account</h2>
          <Button
            variant="destructive"
            onClick={() => logout.mutate()}
            disabled={logout.isPending}
          >
            {logout.isPending ? "Logging out..." : "Logout"}
          </Button>
        </section>
      </div>
    </div>
  );
}
