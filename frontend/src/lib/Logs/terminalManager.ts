import { Terminal } from "xterm";
import { FitAddon } from "@xterm/addon-fit";
import { EventsOn } from "@wails/runtime/runtime";
import "@/../node_modules/xterm/css/xterm.css";

interface TermStore {
  term: Terminal;
  fit: FitAddon;
}

const instances: Record<string, TermStore> = {};

export const TerminalManager = {
  hasInstance(id: string) {
    return !!instances[id];
  },

  // Initializes a terminal and starts listening for events
  initTerminal(id: string) {
    if (instances[id]) return instances[id];

    const term = new Terminal({
      cursorBlink: true,
      theme: { background: "#09090b", foreground: "#fafafa" },
      fontSize: 13,
      fontFamily: "JetBrains Mono, Menlo, monospace",
      convertEol: true,
      scrollback: 10000,
      disableStdin: true,
    });

    const fit = new FitAddon();
    term.loadAddon(fit);

    instances[id] = { term, fit };

    // Unified listener for both global and drive logs
    EventsOn("log-update", (data: { driveId: string; message: string }) => {
      if (data.driveId === id) {
        term.write(data.message);
      }
    });

    return instances[id];
  },

  removeInstance(id: string) {
    if (instances[id] && id !== "global") {
      instances[id].term.dispose();
      delete instances[id];
    }
  },

  // Returns the instance, creating it if it doesn't exist
  getInstance(id: string) {
    return instances[id] || this.initTerminal(id);
  }
};
