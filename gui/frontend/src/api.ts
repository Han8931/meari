// api.ts — a thin, hand-typed bridge to the Go bindings. Wails injects
// window.go.main.App and window.runtime at runtime; typing them here keeps the
// frontend build independent of the wails-generated wailsjs/ directory. Shapes
// mirror the DTOs in gui/app.go.

export interface TreeEntry {
  path: string;
  name: string;
  dir: boolean;
}

export interface NoteMeta {
  path: string;
  title: string;
  subject: string;
  tags: string[] | null;
}

export interface Note extends NoteMeta {
  body: string;
  source: string;
}

export interface CourseMeta {
  id: string;
  title: string;
  level: string;
  path: string;
}

export interface AppInfo {
  offline: boolean;
  vaultDir: string;
  version: string;
}

export interface ChatTurn {
  role: "user" | "assistant";
  content: string;
}

export interface EssayResult {
  score: number;
  feedback: string;
}

interface GoApp {
  Info(): Promise<AppInfo>;
  Tree(): Promise<TreeEntry[]>;
  OpenNote(path: string): Promise<Note>;
  Preview(body: string): Promise<string>;
  Backlinks(path: string): Promise<NoteMeta[]>;
  Search(query: string): Promise<NoteMeta[]>;
  Courses(): Promise<CourseMeta[]>;
  SaveNote(path: string, body: string): Promise<NoteMeta>;
  NewNote(path: string, title: string): Promise<NoteMeta>;
  Rename(oldPath: string, newPath: string): Promise<void>;
  Delete(path: string): Promise<void>;
  GradeEssay(prompt: string, answer: string): Promise<EssayResult>;
  ModelAnswer(prompt: string): Promise<string>;
  StartChat(path: string, history: ChatTurn[]): Promise<string>;
  StartExplain(selection: string): Promise<string>;
  StartPolish(body: string, instruction: string): Promise<string>;
  CancelStream(id: string): Promise<void>;
}

declare global {
  interface Window {
    go?: { main: { App: GoApp } };
    runtime?: {
      EventsOn(event: string, cb: (...data: unknown[]) => void): () => void;
      EventsOff(event: string): void;
    };
  }
}

// api proxies every call to the Wails-bound App. In a plain browser (no Wails
// runtime) calls reject, so the UI degrades instead of throwing at import.
export const api: GoApp = new Proxy({} as GoApp, {
  get(_t, prop: string) {
    return (...args: unknown[]) => {
      const app = window.go?.main?.App as
        | (Record<string, (...a: unknown[]) => unknown> | undefined);
      if (!app) return Promise.reject(new Error("Wails runtime unavailable"));
      return app[prop](...args);
    };
  },
});

export function onEvent(
  event: string,
  cb: (...data: unknown[]) => void,
): () => void {
  if (!window.runtime) return () => {};
  return window.runtime.EventsOn(event, cb);
}

// streamText subscribes to one AI stream: onDelta for each chunk, resolving
// with the full text on done and rejecting on error. Returns an unsubscribe.
export function streamText(
  id: string,
  onDelta: (chunk: string) => void,
): Promise<string> {
  return new Promise((resolve, reject) => {
    const offDelta = onEvent("stream:delta", (sid, chunk) => {
      if (sid === id) onDelta(String(chunk));
    });
    const offDone = onEvent("stream:done", (sid, full) => {
      if (sid !== id) return;
      cleanup();
      resolve(String(full ?? ""));
    });
    const offErr = onEvent("stream:error", (sid, msg) => {
      if (sid !== id) return;
      cleanup();
      reject(new Error(String(msg)));
    });
    function cleanup() {
      offDelta();
      offDone();
      offErr();
    }
  });
}
