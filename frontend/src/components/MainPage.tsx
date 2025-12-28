import { useEffect, useState } from "react";
import { Drive } from "../types";
import DriveList from "./Drives/DriveList";
import LogPanel from "./Logs/LogPanel";
import { GetDrives, EjectDrive, RunRipper } from "@wails/go/main/App";
import { EventsOn } from "@wails/runtime/runtime";
import { Button } from "./ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import JobManager from "./Jobs/JobManager";
import { start } from "repl";

export default function MainPage({}: {}) {
  const [drives, setDrives] = useState<Array<Drive>>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const [jobFile, setJobFile] = useState("");
  const [currentTab, setCurrentTab] = useState("status");

  const getDriveName = (device: string) => {
    const split = device.split("/");
    if (split.length === 0) {
      return device;
    }
    const driveNumber = split[split.length - 1].substring(2);
    return `Disk Drive ${driveNumber}`;
  };

  useEffect(() => {
    GetDrives().then((drives) => {
      setDrives(
        drives.map((drive) => {
          return {
            device: drive,
            name: getDriveName(drive),
            running: false,
          };
        }),
      );
    });
    const off = EventsOn("jobs-complete", () => {
      setIsProcessing(false);
    });

    return off;
  }, []);

  const startJobs = () => {
    setIsProcessing(true);
    RunRipper(jobFile).catch((err) => {
      setIsProcessing(false);
      console.error(err);
    });
  };

  return (
    <div className="size-full flex flex-col">
      <Tabs
        className="size-full"
        defaultValue="status"
        onValueChange={setCurrentTab}
      >
        <TabsList className="w-full flex flex-row items-center justify-center sticky top-0 z-50">
          <div className="bg-white border-solid border-black border-3 rounded-md">
            <TabsTrigger value="status">
              <Button className="m-0.5" disabled={currentTab === "status"}>
                Status
              </Button>
            </TabsTrigger>
            <TabsTrigger value="logs">
              <Button className="m-0.5" disabled={currentTab === "logs"}>
                Logs
              </Button>
            </TabsTrigger>
            <TabsTrigger value="jobs">
              <Button className="m-0.5" disabled={currentTab === "jobs"}>
                Jobs
              </Button>
            </TabsTrigger>
          </div>
        </TabsList>
        <TabsContent
          value="status"
          className="size-full flex flex-col items-center p-4"
        >
          <DriveList
            onEject={(drive) => {
              EjectDrive(drive.device);
            }}
            onLogs={(drive) => {
              console.log(drive.device);
            }}
            drives={drives}
          />
        </TabsContent>
        <TabsContent
          value="logs"
          className="size-full flex flex-col items-center p-4"
        >
          <LogPanel drives={drives} />
        </TabsContent>
        <TabsContent
          value="jobs"
          className="size-full flex flex-col items-center p-4"
        >
          <JobManager
            initialJobFile={jobFile}
            onJobFileChange={setJobFile}
            onJobStart={() => {
              startJobs();
            }}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
