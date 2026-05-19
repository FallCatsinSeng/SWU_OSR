import Link from "next/link";
import { Code2, Github, ExternalLink } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t border-gray-200 bg-white mt-auto">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="py-12 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="col-span-1 sm:col-span-2 lg:col-span-1">
            <Link href="/" className="flex items-center gap-2 group">
              <div className="h-8 w-8 rounded-lg gradient-primary flex items-center justify-center">
                <Code2 className="h-4 w-4 text-white" />
              </div>
              <span className="text-lg font-bold text-gradient">SWU OSR</span>
            </Link>
            <p className="mt-3 text-sm text-gray-500 max-w-xs">
              Open Source Repository — Platform mahasiswa untuk showcase karya,
              kolaborasi, dan membangun portofolio.
            </p>
          </div>

          {/* Platform */}
          <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-3">
              Platform
            </h3>
            <ul className="space-y-2">
              <li>
                <Link
                  href="/"
                  className="text-sm text-gray-500 hover:text-primary-600 transition-colors"
                >
                  Activity Feed
                </Link>
              </li>
              <li>
                <Link
                  href="/showcase"
                  className="text-sm text-gray-500 hover:text-primary-600 transition-colors"
                >
                  Showcase
                </Link>
              </li>
              <li>
                <Link
                  href="/members"
                  className="text-sm text-gray-500 hover:text-primary-600 transition-colors"
                >
                  Members
                </Link>
              </li>
            </ul>
          </div>

          {/* Organization */}
          <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-3">
              Organization
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="https://walisongo.ac.id"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-gray-500 hover:text-primary-600 transition-colors flex items-center gap-1"
                >
                  UIN Walisongo
                  <ExternalLink className="h-3 w-3" />
                </a>
              </li>
              <li>
                <span className="text-sm text-gray-500">
                  HMPSTI FST
                </span>
              </li>
              <li>
                <span className="text-sm text-gray-500">
                  Teknik Informatika
                </span>
              </li>
            </ul>
          </div>

          {/* Resources */}
          <div>
            <h3 className="text-sm font-semibold text-gray-900 mb-3">
              Resources
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="https://github.com/FallCatsinSeng/SWU_OSR"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-gray-500 hover:text-primary-600 transition-colors flex items-center gap-1"
                >
                  <Github className="h-3 w-3" />
                  Source Code
                </a>
              </li>
              <li>
                <Link
                  href="/login"
                  className="text-sm text-gray-500 hover:text-primary-600 transition-colors"
                >
                  Get Started
                </Link>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom bar */}
        <div className="border-t border-gray-100 py-6 flex flex-col sm:flex-row items-center justify-between gap-3">
          <p className="text-xs text-gray-400">
            &copy; {new Date().getFullYear()} HMPSTI — UIN Walisongo Semarang. Built with open source.
          </p>
          <div className="flex items-center gap-4">
            <a
              href="https://github.com/FallCatsinSeng/SWU_OSR"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-gray-600 transition-colors"
              aria-label="GitHub"
            >
              <Github className="h-4 w-4" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
