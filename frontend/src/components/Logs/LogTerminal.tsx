import { useEffect, useRef } from "react";
import { Terminal } from "xterm";
import { FitAddon } from "@xterm/addon-fit";
import { EventsOn } from "@wails/runtime/runtime";
import "@/../node_modules/xterm/css/xterm.css";

const terminalInstances: Record<string, { term: Terminal; fit: FitAddon }> = {};

export default function LogTerminal({ driveId }: { driveId: string }) {
  const terminalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!terminalRef.current) return;
    if (!terminalInstances[driveId]) {
      const term = new Terminal({
        cursorBlink: false,
        theme: { background: "#09090b", foreground: "#fafafa" },
        fontSize: 13,
        fontFamily: "JetBrains Mono, Menlo, monospace",
        convertEol: true,
        scrollback: 5000,
        disableStdin: true,
      });

      const fitAddon = new FitAddon();
      term.loadAddon(fitAddon);
      terminalInstances[driveId] = { term, fit: fitAddon };

      EventsOn("log-update", (data: any) => {
        if (data.driveId === driveId) {
          term.write(data.message);
        }
      });

      if (driveId === "global") {
        EventsOn("global-log", (msg: string) => {
          term.write(msg);
        });
      }
    }

    const { term, fit } = terminalInstances[driveId];
    if (terminalRef.current) {
      terminalRef.current.innerHTML = "";
      term.open(terminalRef.current);
      setTimeout(() => {
        if (terminalRef.current) fit.fit();
      }, 10);
    }

    const resizeObserver = new ResizeObserver(() => {
      window.requestAnimationFrame(() => {
        try {
          fit.fit();
        } catch (e) {
          // Prevent crash if resize happens during unmount
        }
      });
    });

    resizeObserver.observe(terminalRef.current);

    return () => {
      resizeObserver.disconnect();
    };
  }, [driveId]);

  return (
    <div className="size-full bg-[#09090b]">
      <div ref={terminalRef} className="size-full" />
    </div>
  );
}
