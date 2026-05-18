"use client";

import Link from "next/link";
import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { NotificationBell } from "@/features/forum/NotificationBell";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { LogOut, Settings, User } from "lucide-react";

export function Navbar() {
  const { data: user } = useCurrentUser();
  const logout = useLogout();

  return (
    <nav className="border-b border-gray-200 bg-white">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-6">
            <Link href="/" className="text-xl font-bold text-blue-600">
              SWU OSR
            </Link>
            <div className="hidden sm:flex items-center gap-4">
              <Link
                href="/"
                className="text-sm text-gray-600 hover:text-gray-900"
              >
                Feed
              </Link>
              {user && (
                <Link
                  href="/showcase"
                  className="text-sm text-gray-600 hover:text-gray-900"
                >
                  Showcase
                </Link>
              )}
            </div>
          </div>

          <div className="flex items-center gap-3">
            {user ? (
              <>
                <NotificationBell />
                <DropdownMenu
                  trigger={
                    <Avatar
                      src={user.avatar_url}
                      alt={user.alias}
                      fallback={user.alias.charAt(0).toUpperCase()}
                      size="sm"
                      className="cursor-pointer"
                    />
                  }
                >
                  <DropdownMenuItem>
                    <Link
                      href={`/profiles/${user.alias}`}
                      className="flex items-center gap-2"
                    >
                      <User className="h-4 w-4" />
                      Profile
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Link
                      href="/settings"
                      className="flex items-center gap-2"
                    >
                      <Settings className="h-4 w-4" />
                      Settings
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => logout.mutate()}
                    className="text-red-600"
                  >
                    <div className="flex items-center gap-2">
                      <LogOut className="h-4 w-4" />
                      Logout
                    </div>
                  </DropdownMenuItem>
                </DropdownMenu>
              </>
            ) : (
              <Link href="/login">
                <Button size="sm">Sign In</Button>
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
