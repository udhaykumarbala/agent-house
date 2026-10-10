/** Shared on-demand mermaid loader (pinned cdnjs build). Used by the
 *  Documents page and the Conductor chat to render ```mermaid fences as
 *  live diagrams. The Next bundle stays lean — the script loads only when
 *  a diagram is actually on screen. */

const MERMAID_SRC =
  "https://cdnjs.cloudflare.com/ajax/libs/mermaid/10.9.1/mermaid.min.js";

export interface MermaidAPI {
  initialize: (c: object) => void;
  run: (c: object) => Promise<void>;
}

let loading: Promise<MermaidAPI | null> | null = null;

export function ensureMermaid(): Promise<MermaidAPI | null> {
  if (typeof window === "undefined") return Promise.resolve(null);
  const w = window as unknown as { mermaid?: MermaidAPI };
  if (w.mermaid) return Promise.resolve(w.mermaid);
  if (!loading) {
    loading = new Promise((resolve) => {
      const s = document.createElement("script");
      s.src = MERMAID_SRC;
      s.onload = () => {
        w.mermaid?.initialize({ startOnLoad: false, theme: "neutral", securityLevel: "strict" });
        resolve(w.mermaid ?? null);
      };
      s.onerror = () => resolve(null);
      document.head.appendChild(s);
    });
  }
  return loading;
}

/** Render all still-unprocessed .mermaid nodes under root. Safe to call
 *  repeatedly — mermaid marks processed nodes with data-processed. */
export function renderMermaidIn(root: ParentNode | null): void {
  if (!root) return;
  const nodes = root.querySelectorAll(".mermaid:not([data-processed])");
  if (nodes.length === 0) return;
  ensureMermaid().then((m) => {
    m?.run({ nodes }).catch(() => {});
  });
}
