import { Drive } from "@/types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import { Button } from "../ui/button";
import LogViewer from "./LogViewer";
import { useState } from "react";

export default function LogPanel({ drives }: { drives: Array<Drive> }) {
  const [currentTab, setCurrentTab] = useState("status");
  if (drives.length > 0 && drives[drives.length - 1].device !== "global") {
    drives.push({
      device: "global",
      name: "Global",
      running: true,
    });
  }
  return (
    <div className="size-full">
      <Tabs className="size-full" onValueChange={setCurrentTab}>
        <TabsList className="w-full flex flex-row items-center justify-center">
          <div className="bg-white border-solid border-black border-3 rounded-md">
            {drives.map((drive) => {
              return (
                <TabsTrigger value={drive.device}>
                  <Button
                    disabled={drive.device === currentTab}
                    className="m-0.5"
                  >
                    {drive.device}
                  </Button>
                </TabsTrigger>
              );
            })}
          </div>
        </TabsList>
        {drives.map((drive) => {
          return (
            <TabsContent value={drive.device} className="size-full">
              <LogViewer drive={drive} />
            </TabsContent>
          );
        })}
      </Tabs>
    </div>
  );
}
