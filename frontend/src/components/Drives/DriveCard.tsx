import { Disc3, StepForward } from "lucide-react";
import { Drive } from "../../types";
import { Button } from "../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";

export default function DriveCard({
  drive,
  onEject,
  onLogs,
}: {
  drive: Drive;
  onEject: () => void;
  onLogs: () => void;
}) {
  return (
    <div>
      <Card className="min-w-50 h-full">
        <CardHeader className="w-full">
          <CardTitle>
            <div className="flex flex-row items-center justify-left gap-1">
              <Disc3 /> {drive.name}
            </div>
          </CardTitle>
          <CardDescription>{drive.device}</CardDescription>
        </CardHeader>
        <CardContent className="h-full">
          {drive.running ? (
            <div className="flex flex-row items-center justify-between h-full">
              Running
              <Button variant="link" className="font-bold" onClick={onLogs}>
                Logs
              </Button>
            </div>
          ) : (
            <div className="flex flex-row items-center justify-start h-full">
              Idle
            </div>
          )}
        </CardContent>
        <CardFooter>
          <Button variant="default" onClick={onEject}>
            <StepForward className="rotate-270" />
            Eject Disc
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
