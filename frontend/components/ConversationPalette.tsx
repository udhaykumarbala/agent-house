"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import {
  fetchConversations,
  searchConversations,
  type ConversationMeta,
  type ConversationHit,
} from "@/lib/api";

/**
 * Command-palette style history browser. Two modes driven by the query
 * string:
 *   - empty  → list mode (recent conversations, sorted by last activity)
 *   - filled → search mode (text hits across every conversation, with
 *              snippets and a debounced /api/conversations/search call)
 */
export function ConversationPalette({
  open,
  currentId,
  onSelect,
  onClose,
}: {
  open: boolean;
  currentId?: string;
  onSelect: (id: string) => void;
  onClose: () => void;
}) {
  const [q, setQ] = useState("");
  const [list, setList] = useState<ConversationMeta[]>([]);
  const [hits, setHits] = useState<ConversationHit[]>([]);
  const [mode, setMode] = useState<"semantic" | "lexical">("lexical");
  const [loading, setLoading] = useState(false);
  const [activeIdx, setActiveIdx] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  // Reset + focus on open.
  useEffect(() => {
    if (open) {
      setQ("");
      setActiveIdx(0);
      fetchConversations().then(setList);
      setTimeout(() => inputRef.current?.focus(), 30);
    }
  }, [open]);

  // Debounced search.
  useEffect(() => {
    if (!open) return;
    if (!q.trim()) {
      setHits([]);
      return;
    }
    setLoading(true);
    const t = setTimeout(async () => {
      const r = await searchConversations(q);
      setHits(r.hits);
      setMode(r.mode);
      setLoading(false);
      setActiveIdx(0);
    }, 180);
    return () => clearTimeout(t);
  }, [q, open]);

  const rows: { id: string; primary: string; secondary: string }[] = useMemo(() => {
    if (q.trim()) {
      return hits.map((h) => {
        const n = h.match_count ?? 1;
        const tag = n > 1 ? ` · ${n} matches` : "";
        return {
          id: h.conversation_id,
          primary: (h.name || h.conversation_id) + tag,
          secondary: `${h.role} · ${h.snippet}`,
        };
      });
    }
    return list.map((c) => ({
      id: c.id,
      primary: c.name,
      secondary: `${c.message_count} message${c.message_count === 1 ? "" : "s"} · ${
        c.last_at ? new Date(c.last_at).toLocaleString() : "never"
      }`,
    }));
  }, [q, list, hits]);

  // Keyboard navigation.
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      } else if (e.key === "ArrowDown") {
        e.preventDefault();
        setActiveIdx((i) => Math.min(rows.length - 1, i + 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setActiveIdx((i) => Math.max(0, i - 1));
      } else if (e.key === "Enter") {
        e.preventDefault();
        const row = rows[activeIdx];
        if (row) onSelect(row.id);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, rows, activeIdx, onSelect, onClose]);

  if (!open) return null;

  return (
    <div className="palette-backdrop" onClick={onClose}>
      <div className="palette" onClick={(e) => e.stopPropagation()}>
        <div className="palette-input">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <circle cx="11" cy="11" r="7" /><path d="M21 21l-4-4" />
          </svg>
          <input
            ref={inputRef}
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="Search history… (type to filter, ↑↓ to navigate, ⏎ to open, esc to close)"
          />
          {loading && <span className="palette-loading">searching…</span>}
        </div>
        <div className="palette-list">
          {rows.length === 0 && (
            <div className="palette-empty">
              {q.trim() ? "No matches." : "No conversations yet."}
            </div>
          )}
          {rows.map((row, i) => (
            <div
              key={`${row.id}-${i}`}
              className={`palette-row${i === activeIdx ? " active" : ""}${
                row.id === currentId ? " current" : ""
              }`}
              onClick={() => onSelect(row.id)}
              onMouseEnter={() => setActiveIdx(i)}
            >
              <div className="palette-row-primary">{row.primary}</div>
              <div className="palette-row-secondary">{row.secondary}</div>
              {row.id === currentId && (
                <span className="palette-current-tag">current</span>
              )}
            </div>
          ))}
        </div>
        <div className="palette-foot">
          {q.trim() ? (
            <>
              {hits.length} hit{hits.length === 1 ? "" : "s"} for &quot;{q}&quot;
              <span
                className={`palette-mode mode-${mode}`}
                title={
                  mode === "semantic"
                    ? "Cosine-ranked over message embeddings"
                    : "Substring matching (set OPENAI_API_KEY for semantic search)"
                }
              >
                {mode === "semantic" ? "semantic · RAG" : "lexical"}
              </span>
            </>
          ) : (
            `${list.length} conversation${list.length === 1 ? "" : "s"}`
          )}
          <span style={{ marginLeft: "auto" }}>
            <kbd>↑↓</kbd> nav · <kbd>⏎</kbd> open · <kbd>esc</kbd> close
          </span>
        </div>
      </div>
    </div>
  );
}
