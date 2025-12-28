import { Drive } from "../../types";
import { Button } from "../ui/button";
import DriveCard from "./DriveCard";

export default function DriveList({
  drives,
  onEject,
  onLogs,
  onRefresh,
  isRefreshing,
}: {
  drives: Array<Drive>;
  onEject: (drive: Drive) => void;
  onLogs: (drive: Drive) => void;
  onRefresh: () => void;
  isRefreshing: boolean;
}) {
  return (
    <div className="size-full flex flex-col items-center">
      {isRefreshing ? null : (
        <>
          <Button className="" onClick={onRefresh}>
            Refresh
          </Button>
          <div className="flex flex-row gap-5 flex-wrap size-fit">
            {drives
              .filter((drive) => {
                return drive.device !== "global";
              })
              .map((drive) => {
                return (
                  <DriveCard
                    key={drive.device}
                    onEject={() => {
                      onEject(drive);
                    }}
                    onLogs={() => {
                      onLogs(drive);
                    }}
                    drive={drive}
                  />
                );
              })}
          </div>
        </>
      )}
    </div>
  );
}
