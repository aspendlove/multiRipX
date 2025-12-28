import { Drive } from "@/types";
import LogTerminal from "./LogTerminal";

export default function LogViewer({ drive }: { drive: Drive }) {
  return (
    <div className="size-full m-0 border-3 border-black rounded-md flex flex-col">
      <h1 className="w-full text-center font-bold text-2xl">Blah</h1>
      <div className="size-full border-3 border-black rounded-md">
        <LogTerminal drive={drive} />
      </div>
    </div>
  );
}
