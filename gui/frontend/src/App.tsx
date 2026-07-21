// App.tsx — the 3-pane desktop shell: vault tree | editor/preview | tutor chat.
// The editor is CodeMirror + Vim (mouse selection integrates into Visual mode);
// a mode toggle flips between the Vim editor and a rendered markdown preview.

import { useCallback, useEffect, useRef, useState } from "react";
import { api, type AppInfo, type TreeEntry } from "./api";
import Tree from "./Tree";
import Chat from "./Chat";
import Editor, { type EditorHandle } from "./Editor";

type Mode = "read" | "edit";

export default function App() {
  const [info, setInfo] = useState<AppInfo | null>(null);
  const [tree, setTree] = useState<TreeEntry[]>([]);
  const [current, setCurrent] = useState("");
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [mode, setMode] = useState<Mode>("read");
  const [previewHTML, setPreviewHTML] = useState("");
  const [vimEnabled, setVimEnabled] = useState(true);
  const [dirty, setDirty] = useState(false);
  const editorRef = useRef<EditorHandle | null>(null);
  // Live body while editing (CM owns the doc; we mirror it for save/preview).
  const liveBody = useRef(body);

  const refreshTree = useCallback(async () => {
    try {
      setTree(await api.Tree());
    } catch {
      /* browser without the Wails runtime */
    }
  }, []);

  useEffect(() => {
    api.Info().then(setInfo).catch(() => {});
    refreshTree();
  }, [refreshTree]);

  const openNote = useCallback(async (path: string) => {
    const n = await api.OpenNote(path);
    setCurrent(n.path);
    setTitle(n.title + (n.subject ? "  ·  " + n.subject : ""));
    setBody(n.body);
    liveBody.current = n.body;
    setDirty(false);
    setMode("read");
    setPreviewHTML(await api.Preview(n.body));
  }, []);

  const save = useCallback(async () => {
    if (!current) return;
    const text = editorRef.current?.getValue() ?? liveBody.current;
    await api.SaveNote(current, text);
    setDirty(false);
    if (mode === "read") setPreviewHTML(await api.Preview(text));
  }, [current, mode]);

  // Re-render the preview when switching into read mode.
  useEffect(() => {
    if (mode === "read" && current) {
      api.Preview(liveBody.current).then(setPreviewHTML).catch(() => {});
    }
  }, [mode, current]);

  // Ctrl/Cmd-S saves from anywhere.
  useEffect(() => {
    const h = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === "s") {
        e.preventDefault();
        save();
      }
    };
    window.addEventListener("keydown", h);
    return () => window.removeEventListener("keydown", h);
  }, [save]);

  async function explainSelection() {
    const sel = editorRef.current?.getSelection()?.trim();
    if (!sel) return;
    // Route through the chat pane by seeding a stream; the Chat component owns
    // its own transcript, so for now open the explanation as a toast-like note.
    const id = await api.StartExplain(sel);
    // The Chat pane isn't wired to external streams yet; a minimal surfacing:
    console.log("explain stream", id);
  }

  return (
    <div className="app">
      <header className="titlebar">
        <span className="brand">메아리 Meari</span>
        <span className="note-title">{title || "—"}</span>
        <span className="spacer" />
        {current && (
          <div className="mode-toggle">
            <button
              className={mode === "read" ? "on" : ""}
              onClick={() => setMode("read")}
            >
              Read
            </button>
            <button
              className={mode === "edit" ? "on" : ""}
              onClick={() => setMode("edit")}
            >
              Edit
            </button>
          </div>
        )}
        <label className="vim-toggle">
          <input
            type="checkbox"
            checked={vimEnabled}
            onChange={(e) => setVimEnabled(e.target.checked)}
          />
          Vim
        </label>
        {dirty && <span className="dirty">●</span>}
        {info?.offline && <span className="offline">offline</span>}
      </header>

      <div className="panes">
        <aside className="pane tree-pane">
          <Tree entries={tree} current={current} onOpen={openNote} />
        </aside>

        <main className="pane editor-pane">
          {!current ? (
            <div className="placeholder">Open a note from the tree.</div>
          ) : mode === "edit" ? (
            <Editor
              key={current + (vimEnabled ? "#vim" : "#plain")}
              doc={body}
              vimEnabled={vimEnabled}
              onChange={(v) => {
                liveBody.current = v;
                if (!dirty) setDirty(true);
              }}
              onSave={save}
              onReady={(h) => (editorRef.current = h)}
            />
          ) : (
            <div
              className="preview"
              dangerouslySetInnerHTML={{ __html: previewHTML }}
            />
          )}
          {mode === "edit" && (
            <div className="editor-actions">
              <button onClick={explainSelection}>Explain selection</button>
              <button onClick={save}>Save (⌘S)</button>
            </div>
          )}
        </main>

        <aside className="pane chat-pane">
          <Chat notePath={current} />
        </aside>
      </div>
    </div>
  );
}
