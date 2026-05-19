import { Card, CardContent } from "@/components/ui/card";

export default function MembersPage() {
  return (
    <div className="mx-auto max-w-4xl px-4 py-8">
      <h1 className="text-2xl font-bold text-gray-900 mb-6">
        Discover Members
      </h1>
      <Card>
        <CardContent className="p-6 text-center">
          <p className="text-gray-600">
            Member discovery is coming soon. In the meantime, you can find
            profiles by visiting their activity feed links.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
