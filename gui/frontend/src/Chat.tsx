// Chat.tsx — the tutor conversation pane. Streams replies token-by-token via
// the Wails event bridge, grounded on the open note.

import { useRef, useState } from "react";
import { api, streamText, type ChatTurn } from "./api";

interface ChatProps {
  notePath: string;
}

export default function Chat({ notePath }: ChatProps) {
  const [turns, setTurns] = useState<ChatTurn[]>([]);
  const [draft, setDraft] = useState("");
  const [streaming, setStreaming] = useState(false);
  const streamID = useRef<string>("");

  async function send(text: string) {
    const q = text.trim();
    if (!q || streaming) return;
    const next: ChatTurn[] = [...turns, { role: "user", content: q }];
    setTurns([...next, { role: "assistant", content: "" }]);
    setDraft("");
    setStreaming(true);
    try {
      const id = await api.StartChat(notePath, next);
      streamID.current = id;
      await streamText(id, (chunk) => {
        setTurns((cur) => {
          const copy = cur.slice();
          const last = copy[copy.length - 1];
          copy[copy.length - 1] = { ...last, content: last.content + chunk };
          return copy;
        });
      });
    } catch (e) {
      setTurns((cur) => {
        const copy = cur.slice();
        copy[copy.length - 1] = {
          role: "assistant",
          content: "⚠ " + (e as Error).message,
        };
        return copy;
      });
    } finally {
      setStreaming(false);
      streamID.current = "";
    }
  }

  function stop() {
    if (streamID.current) api.CancelStream(streamID.current);
  }

  return (
    <div className="chat">
      <div className="chat-log">
        {turns.length === 0 && (
          <div className="chat-empty">Ask the tutor about this note.</div>
        )}
        {turns.map((t, i) => (
          <div key={i} className={"turn turn-" + t.role}>
            <span className="badge">{t.role === "user" ? "you" : "tutor"}</span>
            <div className="turn-body">{t.content || (streaming ? "…" : "")}</div>
          </div>
        ))}
      </div>
      <div className="chat-input">
        <textarea
          value={draft}
          placeholder="ask the tutor…"
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && !e.shiftKey) {
              e.preventDefault();
              send(draft);
            }
          }}
        />
        <button onClick={() => (streaming ? stop() : send(draft))}>
          {streaming ? "Stop" : "Send"}
        </button>
      </div>
    </div>
  );
}
