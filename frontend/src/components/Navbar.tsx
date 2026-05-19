"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useCurrentUser, useLogout } from "@/hooks/useAuth";
import { NotificationBell } from "@/features/forum/NotificationBell";
import { Avatar } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuItem } from "@/components/ui/dropdown-menu";
import {
  LogOut,
  Settings,
  User,
  Menu,
  X,
  Home,
  FolderGit2,
  Users,
  Code2,
} from "lucide-react";

const NAV_LINKS = [
  { href: "/", label: "Feed", icon: Home },
  { href: "/showcase", label: "Showcase", icon: FolderGit2, auth: true },
  { href: "/members", label: "Members", icon: Users, auth: true },
];

export function Navbar() {
  const { data: user } = useCurrentUser();
  const logout = useLogout();
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);

  const isActive = (href: string) => {
    if (href === "/") return pathname === "/";
    return pathname.startsWith(href);
  };

  return (
    <>
      {/* Geist nav-bar: sticky, 64px, white bg, bottom hairline */}
      <nav className="sticky top-0 z-50 h-16 border-b border-geist-hairline bg-geist-canvas dark:border-neutral-800 dark:bg-black">
        <div className="mx-auto max-w-geist-page px-6 h-full">
          <div className="flex h-full items-center justify-between">
            {/* Left: Logo + Nav links */}
            <div className="flex items-center gap-8">
              {/* Logo */}
              <Link href="/" className="flex items-center gap-2">
                <div className="h-7 w-7 rounded-geist-sm bg-geist-primary dark:bg-white flex items-center justify-center">
                  <Code2 className="h-3.5 w-3.5 text-geist-on-primary dark:text-black" />
                </div>
                <span className="text-body-sm-strong text-geist-ink dark:text-neutral-50 hidden sm:inline">
                  SWU OSR
                </span>
              </Link>

              {/* Desktop Navigation — centre link row */}
              <div className="hidden md:flex items-center gap-1">
                {NAV_LINKS.map((link) => {
                  if (link.auth && !user) return null;
                  return (
                    <Link
                      key={link.href}
                      href={link.href}
                      className={`px-3 py-2 rounded-geist-full text-body-sm transition-colors ${
                        isActive(link.href)
                          ? "text-geist-ink bg-geist-canvas-soft-2 dark:text-neutral-50 dark:bg-neutral-800"
                          : "text-geist-body hover:text-geist-ink dark:text-neutral-400 dark:hover:text-neutral-50"
                      }`}
                    >
                      {link.label}
                    </Link>
                  );
                })}
              </div>
            </div>

            {/* Right: CTAs */}
            <div className="flex items-center gap-3">
              {user ? (
                <>
                  <NotificationBell />
                  <DropdownMenu
                    trigger={
                      <div className="flex items-center cursor-pointer p-1 rounded-geist-full hover:bg-geist-canvas-soft dark:hover:bg-neutral-800 transition-colors">
                        <Avatar
                          src={user.avatar_url}
                          alt={user.alias}
                          fallback={user.alias.charAt(0).toUpperCase()}
                          size="sm"
                        />
                      </div>
                    }
                  >
                    <div className="px-3 py-2 border-b border-geist-hairline dark:border-neutral-800 max-w-[200px]">
                      <p className="text-body-sm-strong text-geist-ink dark:text-neutral-50 break-words">
                        {user.alias}
                      </p>
                      <p className="text-caption text-geist-mute dark:text-neutral-500 truncate">
                        {user.nim}
                      </p>
                    </div>
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
                      className="text-geist-error dark:text-red-400"
                    >
                      <div className="flex items-center gap-2">
                        <LogOut className="h-4 w-4" />
                        Logout
                      </div>
                    </DropdownMenuItem>
                  </DropdownMenu>
                </>
              ) : (
                <div className="flex items-center gap-2">
                  <Link href="/login">
                    <Button variant="nav-secondary" size="nav">
                      Log In
                    </Button>
                  </Link>
                  <Link href="/login">
                    <Button variant="nav-primary" size="nav">
                      Sign Up
                    </Button>
                  </Link>
                </div>
              )}

              {/* Mobile hamburger */}
              <button
                onClick={() => setMobileOpen(!mobileOpen)}
                className="md:hidden p-2 rounded-geist-sm text-geist-body hover:bg-geist-canvas-soft-2 transition-colors dark:text-neutral-400 dark:hover:bg-neutral-800"
                aria-label="Toggle menu"
              >
                {mobileOpen ? (
                  <X className="h-5 w-5" />
                ) : (
                  <Menu className="h-5 w-5" />
                )}
              </button>
            </div>
          </div>
        </div>
      </nav>

      {/* Mobile menu — full overlay */}
      {mobileOpen && (
        <div className="md:hidden fixed inset-x-0 top-16 z-40 bg-geist-canvas dark:bg-black border-b border-geist-hairline dark:border-neutral-800 geist-level-4 animate-slide-down">
          <div className="px-4 py-3 space-y-1">
            {NAV_LINKS.map((link) => {
              if (link.auth && !user) return null;
              const Icon = link.icon;
              return (
                <Link
                  key={link.href}
                  href={link.href}
                  onClick={() => setMobileOpen(false)}
                  className={`flex items-center gap-3 px-3 py-3 rounded-geist-sm text-body-sm transition-colors ${
                    isActive(link.href)
                      ? "text-geist-ink bg-geist-canvas-soft-2 dark:text-neutral-50 dark:bg-neutral-800"
                      : "text-geist-body hover:text-geist-ink hover:bg-geist-canvas-soft dark:text-neutral-400 dark:hover:text-neutral-50 dark:hover:bg-neutral-900"
                  }`}
                >
                  <Icon className="h-4 w-4" />
                  {link.label}
                </Link>
              );
            })}
            {user && (
              <>
                <div className="border-t border-geist-hairline dark:border-neutral-800 my-2" />
                <Link
                  href={`/profiles/${user.alias}`}
                  onClick={() => setMobileOpen(false)}
                  className="flex items-center gap-3 px-3 py-3 rounded-geist-sm text-body-sm text-geist-body hover:text-geist-ink hover:bg-geist-canvas-soft dark:text-neutral-400 dark:hover:text-neutral-50 dark:hover:bg-neutral-900"
                >
                  <User className="h-4 w-4" />
                  My Profile
                </Link>
                <Link
                  href="/settings"
                  onClick={() => setMobileOpen(false)}
                  className="flex items-center gap-3 px-3 py-3 rounded-geist-sm text-body-sm text-geist-body hover:text-geist-ink hover:bg-geist-canvas-soft dark:text-neutral-400 dark:hover:text-neutral-50 dark:hover:bg-neutral-900"
                >
                  <Settings className="h-4 w-4" />
                  Settings
                </Link>
              </>
            )}
          </div>
        </div>
      )}

      {/* Overlay backdrop */}
      {mobileOpen && (
        <div
          className="md:hidden fixed inset-0 top-16 z-30 bg-geist-ink/20 dark:bg-black/60"
          onClick={() => setMobileOpen(false)}
        />
      )}
    </>
  );
}
