// Editor.tsx — a CodeMirror 6 markdown editor with Vim keybindings. The
// @replit/codemirror-vim extension gives full Normal/Visual/Insert modes,
// registers, counts, and text objects; crucially, MOUSE selection integrates
// natively — dragging with the mouse enters Visual mode over the dragged span,
// so `y`/`d`/`c` act on a mouse selection exactly as in a real Vim + terminal.
//
// The Vim extension MUST be first in the extensions array (its own docs), so it
// wins key precedence over the default keymap.

import { useEffect, useRef } from "react";
import { EditorState, type Extension } from "@codemirror/state";
import { EditorView, keymap, lineNumbers } from "@codemirror/view";
import { defaultKeymap, history, historyKeymap } from "@codemirror/commands";
import { markdown } from "@codemirror/lang-markdown";
import { syntaxHighlighting, defaultHighlightStyle } from "@codemirror/language";
import { vim, Vim } from "@replit/codemirror-vim";

// Register :w / :wq to trigger a save through the host, once per module load.
let exCommandsRegistered = false;
function registerExCommands(save: () => void) {
  if (exCommandsRegistered) return;
  exCommandsRegistered = true;
  Vim.defineEx("write", "w", save);
  Vim.defineEx("wq", "wq", save);
  Vim.defineEx("x", "x", save);
}

export interface EditorHandle {
  getValue(): string;
  getSelection(): string;
}

interface EditorProps {
  doc: string;
  vimEnabled: boolean;
  onChange(value: string): void;
  onSave(): void;
  onReady(handle: EditorHandle): void;
}

export default function Editor({
  doc,
  vimEnabled,
  onChange,
  onSave,
  onReady,
}: EditorProps) {
  const host = useRef<HTMLDivElement>(null);
  const view = useRef<EditorView | null>(null);
  // Keep the latest callbacks reachable from CM without rebuilding the view.
  const cb = useRef({ onChange, onSave });
  cb.current = { onChange, onSave };

  // Rebuild the view when the doc identity or vim toggle changes. We key the
  // outer component by note path (see App), so `doc` here is the initial text.
  useEffect(() => {
    if (!host.current) return;
    registerExCommands(() => cb.current.onSave());

    const base: Extension[] = [
      history(),
      lineNumbers(),
      markdown(),
      syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
      EditorView.lineWrapping,
      keymap.of([...defaultKeymap, ...historyKeymap]),
      EditorView.updateListener.of((u) => {
        if (u.docChanged) cb.current.onChange(u.state.doc.toString());
      }),
      EditorView.theme({
        "&": { height: "100%", fontSize: "14px" },
        ".cm-scroller": { fontFamily: "var(--mono)", overflow: "auto" },
        ".cm-content": { caretColor: "var(--accent)" },
      }),
    ];
    // Vim first, so its keymap takes precedence (per the extension's docs).
    const extensions = vimEnabled ? [vim(), ...base] : base;

    const state = EditorState.create({ doc, extensions });
    const v = new EditorView({ state, parent: host.current });
    view.current = v;

    onReady({
      getValue: () => v.state.doc.toString(),
      getSelection: () => {
        const { from, to } = v.state.selection.main;
        return v.state.sliceDoc(from, to);
      },
    });

    return () => v.destroy();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [vimEnabled]);

  // Reflect external doc replacement (opening a different note reuses the view
  // only if same vim state; App keys by path so this mainly guards reloads).
  useEffect(() => {
    const v = view.current;
    if (!v) return;
    if (v.state.doc.toString() !== doc) {
      v.dispatch({
        changes: { from: 0, to: v.state.doc.length, insert: doc },
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [doc]);

  return <div className="editor-host" ref={host} />;
}
