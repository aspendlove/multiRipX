import { useEffect, useRef } from "react";
import { TerminalManager } from "@/lib/Logs/terminalManager";

export default function LogTerminal({ driveId }: { driveId: string }) {
  const terminalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!terminalRef.current) return;

    const { term, fit } = TerminalManager.getInstance(driveId);

    // Prepare container and attach
    terminalRef.current.innerHTML = "";
    term.open(terminalRef.current);

    // Small delay for WebKitGTK layout calculation
    const timeout = setTimeout(() => fit.fit(), 20);

    const resizeObserver = new ResizeObserver(() => {
      window.requestAnimationFrame(() => fit.fit());
    });
    resizeObserver.observe(terminalRef.current);

    return () => {
      clearTimeout(timeout);
      resizeObserver.disconnect();
    };
  }, [driveId]);

  return (
    <div className="size-full overflow-hidden bg-[#09090b]">
      <div ref={terminalRef} className="size-full" />
    </div>
  );
}
