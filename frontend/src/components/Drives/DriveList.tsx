import { Drive } from "../../types";
import DriveCard from "./DriveCard";

export default function DriveList({
  drives,
  onEject,
  onLogs,
}: {
  drives: Array<Drive>;
  onEject: (drive: Drive) => void;
  onLogs: (drive: Drive) => void;
}) {
  return (
    <div className="flex flex-row gap-5 flex-wrap">
      {drives.map((drive) => {
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
  );
}
