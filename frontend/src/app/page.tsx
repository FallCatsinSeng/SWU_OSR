import { ActivityFeed } from "@/features/feed/ActivityFeed";

export default function HomePage() {
  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <div className="mb-8 text-center">
        <h1 className="text-4xl font-bold text-gray-900 mb-3">
          SWU Open Source Repository
        </h1>
        <p className="text-lg text-gray-600 max-w-2xl mx-auto">
          A pseudonymous academic open-source showcase platform for
          Srinakharinwirot University students and faculty.
        </p>
      </div>

      <section>
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">
          Recent Activity
        </h2>
        <ActivityFeed />
      </section>
    </div>
  );
}
