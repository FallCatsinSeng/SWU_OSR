"use client";

interface ContributionHeatmapProps {
  activeDays: number;
  totalCommits: number;
}

// Generate a simulated heatmap based on user stats
function generateHeatmapData(activeDays: number, totalCommits: number) {
  const weeks = 20;
  const days = weeks * 7;
  const data: number[] = [];

  // Use total commits to seed randomness deterministically
  let seed = totalCommits + activeDays;
  const pseudoRandom = () => {
    seed = (seed * 1664525 + 1013904223) % 4294967296;
    return seed / 4294967296;
  };

  const avgPerDay = activeDays > 0 ? totalCommits / activeDays : 0;

  for (let i = 0; i < days; i++) {
    const rand = pseudoRandom();
    // Higher chance of activity for users with more active days
    const activityChance = Math.min(activeDays / days, 0.7);

    if (rand < activityChance) {
      const intensity = pseudoRandom();
      if (intensity < 0.4) data.push(1);
      else if (intensity < 0.7) data.push(2);
      else if (intensity < 0.9) data.push(3);
      else data.push(4);
    } else {
      data.push(0);
    }
  }

  return data;
}

export function ContributionHeatmap({
  activeDays,
  totalCommits,
}: ContributionHeatmapProps) {
  const data = generateHeatmapData(activeDays, totalCommits);
  const weeks = 20;

  return (
    <div className="overflow-x-auto">
      <div className="inline-flex gap-0.5">
        {Array.from({ length: weeks }).map((_, weekIdx) => (
          <div key={weekIdx} className="flex flex-col gap-0.5">
            {Array.from({ length: 7 }).map((_, dayIdx) => {
              const idx = weekIdx * 7 + dayIdx;
              const level = data[idx] ?? 0;
              return (
                <div
                  key={dayIdx}
                  className={`heatmap-cell w-3 h-3 heatmap-level-${level}`}
                  title={`${level} contribution${level !== 1 ? "s" : ""}`}
                />
              );
            })}
          </div>
        ))}
      </div>
      {/* Legend */}
      <div className="flex items-center gap-1 mt-2 justify-end">
        <span className="text-[10px] text-gray-400 mr-1">Less</span>
        {[0, 1, 2, 3, 4].map((level) => (
          <div
            key={level}
            className={`w-3 h-3 rounded-sm heatmap-level-${level}`}
          />
        ))}
        <span className="text-[10px] text-gray-400 ml-1">More</span>
      </div>
    </div>
  );
}
