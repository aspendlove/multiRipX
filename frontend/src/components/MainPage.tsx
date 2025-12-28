import { useEffect, useState, useCallback } from "react";
import { Drive } from "../types";
import DriveList from "./Drives/DriveList";
import LogPanel from "./Logs/LogPanel";
import { GetDrives, EjectDrive, RunRipper } from "@wails/go/main/App";
import { EventsOn } from "@wails/runtime/runtime";
import { Button } from "./ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import JobManager from "./Jobs/JobManager";
import { TerminalManager } from "@/lib/Logs/terminalManager";

export default function MainPage({}: {}) {
  const [drives, setDrives] = useState<Array<Drive>>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const [jobFile, setJobFile] = useState("");
  const [currentTab, setCurrentTab] = useState("status");
  const [isRefreshing, setIsRefreshing] = useState(false);

  const getDriveName = (device: string) => {
    const split = device.split("/");
    if (split.length === 0) return device;
    const driveNumber = split[split.length - 1].substring(2);
    return `Disk Drive ${driveNumber}`;
  };

  const refreshDrives = useCallback(async () => {
    setIsRefreshing(true);
    try {
      const detectedPaths = await GetDrives();
      setDrives((currentDrives) => {
        currentDrives.forEach((oldDrive) => {
          if (!detectedPaths.includes(oldDrive.device)) {
            console.log(
              `Device ${oldDrive.device} removed. Disposing terminal...`,
            );
            TerminalManager.removeInstance(oldDrive.device);
          }
        });

        return detectedPaths.map((path) => ({
          device: path,
          name: getDriveName(path),
          running: false,
        }));
      });

      detectedPaths.forEach((path) => {
        TerminalManager.initTerminal(path);
      });
    } catch (err) {
      console.error("Failed to refresh drives:", err);
    } finally {
      setIsRefreshing(false);
    }
  }, []);

  useEffect(() => {
    TerminalManager.initTerminal("global");
    refreshDrives();

    const off = EventsOn("jobs-complete", () => {
      setIsProcessing(false);
    });

    return off;
  }, [refreshDrives]);

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
        <TabsList className="w-full flex flex-row items-center justify-center">
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
          className="w-full flex flex-col items-center justify-start p-4"
        >
          <DriveList
            onEject={(drive) => {
              if (isProcessing) return;
              EjectDrive(drive.device);
            }}
            onRefresh={refreshDrives}
            isRefreshing={isRefreshing}
            drives={drives}
            onLogs={() => setCurrentTab("logs")}
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
          className="w-full flex flex-col items-center p-4"
        >
          <JobManager
            initialJobFile={jobFile}
            onJobFileChange={setJobFile}
            onJobStart={startJobs}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
