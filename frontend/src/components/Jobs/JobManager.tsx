import { GetCurrentJobs, OpenFile } from "@wails/go/main/App";
import { Button } from "../ui/button";
import { config } from "@wails/go/models";
import { useEffect, useState } from "react";
import JobVisualizer from "./JobVisualizer";
import { Play } from "lucide-react";

export default function JobManager({
  initialJobFile,
  onJobFileChange,
  onJobStart,
}: {
  initialJobFile: string;
  onJobFileChange: (jobFile: string) => void;
  onJobStart: (jobFile: string) => void;
}) {
  const [jobs, setJobs] = useState<config.JobsConfig>();
  const [jobFile, setJobFile] = useState(initialJobFile);
  useEffect(() => {
    GetCurrentJobs(initialJobFile).then((jobs) => {
      setJobs(jobs);
    });
  }, []);

  const fileSelectorButton = (
    <Button
      className="flex-1 mx-0.5"
      onClick={() => {
        OpenFile().then((newJobFile) => {
          setJobFile(newJobFile);
          onJobFileChange(newJobFile);
          GetCurrentJobs(newJobFile).then((jobs) => {
            setJobs(jobs);
          });
        });
      }}
    >
      {jobFile ? jobFile : "Job File"}
    </Button>
  );

  return (
    <div className="size-full flex flex-col">
      {jobFile ? (
        <div className="flex flex-row items-stretch justify-center w-full">
          {fileSelectorButton}
          <Button
            className="flex-1 mx-0.5"
            onClick={() => {
              if (jobs) {
                onJobStart(jobFile);
              }
            }}
          >
            <Play />
            Run
          </Button>
        </div>
      ) : (
        fileSelectorButton
      )}
      {jobFile && jobs ? <JobVisualizer data={jobs} /> : null}
    </div>
  );
}
