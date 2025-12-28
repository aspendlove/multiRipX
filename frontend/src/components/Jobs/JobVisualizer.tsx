import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Disc, Tv, Film, HardDrive, FolderOpen } from "lucide-react";
import { config } from "@wails/go/models";

export default function JobVisualizer({ data }: { data: config.JobsConfig }) {
  return (
    /* transform-gpu forces WebKitGTK to use the compositor layer */
    <div className="flex flex-col h-full space-y-4 transform-gpu">
      {/* Header Area */}
      <div className="flex items-center justify-between px-1">
        <div>
          <h2 className="text-xl font-bold tracking-tight">Jobs</h2>
          <div className="flex items-center text-muted-foreground text-xs mt-1">
            <FolderOpen className="mr-2 h-3 w-3" />
            <span className="font-mono">{data.outputDir}</span>
          </div>
        </div>
        <Badge variant="outline">{data.jobs.length} Devices</Badge>
      </div>

      <Separator />

      {/* Vertical List instead of Grid - much easier on WebKit layout engine */}
      <div className="flex flex-col gap-3 overflow-y-auto pr-2 pb-10">
        {data.jobs.map((job, idx) => (
          <Card
            key={`${job.drive}-${idx}`}
            className="border border-border bg-card shadow-none rounded-sm"
          >
            <CardContent className="p-0">
              {/* Device Header - Use solid backgrounds, no transparency */}
              <div className="flex items-center justify-between bg-secondary p-3 border-b border-border">
                <div className="flex items-center space-x-3">
                  <HardDrive className="h-4 w-4 text-primary" />
                  <span className="font-bold text-sm">{job.drive}</span>
                  <span className="text-[10px] font-mono text-muted-foreground bg-background px-2 py-0.5 rounded border border-border">
                    {job.outputDir || "root"}
                  </span>
                </div>

                <div className="flex items-center space-x-2">
                  <Badge
                    className={
                      job.discType === "bluray"
                        ? "bg-blue-700 hover:bg-blue-700"
                        : "bg-orange-700 hover:bg-orange-700"
                    }
                  >
                    <Disc className="mr-1 h-3 w-3" />
                    {job.discType?.toUpperCase()}
                  </Badge>
                </div>
              </div>

              {/* Content Area */}
              <div className="p-3 space-y-3">
                {job.shows && job.shows.length > 0 && (
                  <div className="space-y-1">
                    <div className="flex items-center text-[10px] font-bold uppercase text-muted-foreground mb-1">
                      <Tv className="mr-2 h-3 w-3" /> TV
                    </div>
                    {job.shows.map((show, sIdx) => (
                      <div
                        key={sIdx}
                        className="flex items-center justify-between text-xs bg-background p-2 border border-border rounded-sm"
                      >
                        <span className="font-medium">{show.name}</span>
                        <span className="text-muted-foreground font-mono">
                          S{show.season}E{show.episode} [T{show.title}]
                        </span>
                      </div>
                    ))}
                  </div>
                )}

                {job.movies && job.movies.length > 0 && (
                  <div className="space-y-1">
                    <div className="flex items-center text-[10px] font-bold uppercase text-muted-foreground mb-1">
                      <Film className="mr-2 h-3 w-3" /> Movies
                    </div>
                    {job.movies.map((movie, mIdx) => (
                      <div
                        key={mIdx}
                        className="flex justify-between items-center text-xs bg-background p-2 border border-border rounded-sm"
                      >
                        <span className="font-medium">{movie.name}</span>
                        <Badge variant="outline" className="text-[10px] h-4">
                          Track {movie.title}
                        </Badge>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
