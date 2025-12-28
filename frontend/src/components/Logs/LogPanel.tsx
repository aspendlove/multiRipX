import { Drive } from "@/types";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@radix-ui/react-tabs";
import { Button } from "../ui/button";
import LogViewer from "./LogViewer";

export default function LogPanel({ drives }: { drives: Array<Drive> }) {
  return (
    <div className="size-full">
      <Tabs defaultValue={drives[0].device} className="size-full">
        <TabsList className="w-full flex flex-row items-center justify-center">
          <div className="bg-white border-solid border-black border-3 rounded-md">
            {drives.map((drive) => {
              return (
                <TabsTrigger value={drive.device}>
                  <Button className="m-0.5">{drive.device}</Button>
                </TabsTrigger>
              );
            })}
          </div>
        </TabsList>
        {drives.map((drive) => {
          return (
            <TabsContent value={drive.device} className="size-full">
              <LogViewer />
            </TabsContent>
          );
        })}
      </Tabs>
    </div>
  );
}
