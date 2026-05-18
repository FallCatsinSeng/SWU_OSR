"use client";

import { RepoSelector } from "@/features/showcase/RepoSelector";
import { ShowcaseGrid } from "@/features/showcase/ShowcaseGrid";

export default function ShowcasePage() {
  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">
        Manage Your Showcase
      </h1>

      <section className="mb-8">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">
          Your Showcase Repositories
        </h2>
        <ShowcaseGrid />
      </section>

      <section>
        <h2 className="text-xl font-semibold text-gray-900 mb-4">
          Add Repositories
        </h2>
        <RepoSelector />
      </section>
    </div>
  );
}
