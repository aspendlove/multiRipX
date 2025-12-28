import { Drive } from "../types";
import DriveList from "./Drives/DriveList";
import LogPanel from "./Logs/LogPanel";

const drives = [
  { device: "/dev/sr0", name: "Drive 1", running: true },
  { device: "/dev/sr1", name: "Drive 2", running: false },
  { device: "/dev/sr2", name: "Drive 3", running: false },
  { device: "/dev/sr3", name: "Drive 4", running: true },
  { device: "/dev/sr4", name: "Drive 5", running: true },
  { device: "/dev/sr5", name: "Drive 6", running: false },
  { device: "/dev/sr6", name: "Drive 7", running: false },
  { device: "/dev/sr7", name: "Drive 8", running: true },
];

export default function MainPage({}: {}) {
  return (
    <div className="size-full">
      <DriveList
        onEject={(drive) => {
          console.log(drive.device);
        }}
        onLogs={(drive) => {
          console.log(drive.device);
        }}
        drives={drives}
      />
      <LogPanel drives={drives} />
    </div>
  );
}
