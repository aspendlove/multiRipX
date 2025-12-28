import { useEffect, useRef } from "react";
import { Terminal } from "xterm";
import { EventsOn } from "../../../wailsjs/runtime";
import "@xterm/xterm/css/xterm.css";

export default function LogTerminal({ driveId }: { driveId: string }) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const xterm = useRef<Terminal | null>(null);

  useEffect(() => {
    if (!terminalRef.current) return;

    xterm.current = new Terminal({
      cursorBlink: true,
      theme: {
        background: "#09090b",
        foreground: "#fafafa",
      },
      fontSize: 13,
      fontFamily: "JetBrains Mono, Menlo, monospace",
      convertEol: true,
    });

    xterm.current.open(terminalRef.current);

    // 3. Listen for Wails events specific to this drive
    const quitListening = EventsOn("log-update", (data) => {
      if (data.driveId === driveId) {
        xterm.current?.write(data.message);
      }
    });

    return () => {
      quitListening();
      xterm.current?.dispose();
    };
  }, [driveId]);

  return (
    <div className="rounded-md border border-zinc-800 overflow-hidden">
      <div
        ref={terminalRef}
        className="p-2 bg-[#09090b]"
        style={{ minHeight: "300px" }}
      />
    </div>
  );
}
