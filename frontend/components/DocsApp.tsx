"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { Chrome } from "./Chrome";
import { API_BASE } from "@/lib/api";
import { renderMarkdown } from "@/lib/md";
import { getCompany, companyConfig } from "@/lib/company";

/** Documents — rich agent-authored artifacts (PRDs, research, designs)
 *  rendered as real documents: GFM tables, code blocks, and live mermaid
 *  flowcharts. Sources: projects/<id>/discovery/ (client-discovery
 *  pipeline) and projects/<id>/.plans/ (the SDLC orchestrator). */

interface DocRow {
  path: string;
  name: string;
  title: string;
  size: number;
  modified: number;
}
interface ProjectRow {
  id: string;
  docs: number;
  status?: string;
  modified: number;
}
interface ManifestStep {
  agent: string;
  action: string;
  note?: string;
  ok: boolean;
}
interface Manifest {
  status: string;
  current?: string;
  steps: ManifestStep[];
}

// Pinned mermaid from cdnjs, loaded once on demand (same precedent as the
// legacy dashboard's CDN marked.js). The Next bundle stays lean.
const MERMAID_SRC =
  "https://cdnjs.cloudflare.com/ajax/libs/mermaid/10.9.1/mermaid.min.js";
let mermaidLoading: Promise<unknown> | null = null;
function ensureMermaid(): Promise<unknown> {
  if (typeof window === "undefined") return Promise.resolve(null);
  const w = window as unknown as { mermaid?: { initialize: (c: object) => void; run: (c: object) => Promise<void> } };
  if (w.mermaid) return Promise.resolve(w.mermaid);
  if (!mermaidLoading) {
    mermaidLoading = new Promise((resolve, reject) => {
      const s = document.createElement("script");
      s.src = MERMAID_SRC;
      s.onload = () => {
        w.mermaid?.initialize({ startOnLoad: false, theme: "neutral", securityLevel: "strict" });
        resolve(w.mermaid);
      };
      s.onerror = reject;
      document.head.appendChild(s);
    });
  }
  return mermaidLoading;
}

export function DocsApp() {
  const company = getCompany();
  const cfg = companyConfig(company);

  const [projects, setProjects] = useState<ProjectRow[]>([]);
  const [project, setProject] = useState<string>("");
  const [docs, setDocs] = useState<DocRow[]>([]);
  const [status, setStatus] = useState<string>("");
  const [manifest, setManifest] = useState<Manifest | null>(null);
  const [docPath, setDocPath] = useState<string>("");
  const [html, setHtml] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const bodyRef = useRef<HTMLDivElement>(null);

  const loadProjects = useCallback(async () => {
    try {
      const r = await fetch(`${API_BASE}/api/docs`);
      const d = await r.json();
      setProjects(d.projects ?? []);
      return (d.projects ?? []) as ProjectRow[];
    } catch {
      return [];
    }
  }, []);

  const loadDocs = useCallback(async (p: string) => {
    if (!p) return;
    try {
      const r = await fetch(`${API_BASE}/api/docs?project=${encodeURIComponent(p)}`);
      const d = await r.json();
      setDocs(d.docs ?? []);
      setStatus(d.status ?? "");
    } catch {
      /* server unreachable — keep whatever we have */
    }
    try {
      const r = await fetch(`${API_BASE}/api/discovery/status?project=${encodeURIComponent(p)}`);
      setManifest(r.ok ? await r.json() : null);
    } catch {
      setManifest(null);
    }
  }, []);

  const openDoc = useCallback(async (p: string, path: string) => {
    setDocPath(path);
    setLoading(true);
    try {
      const r = await fetch(
        `${API_BASE}/api/file-content?path=${encodeURIComponent(p + "/" + path)}`
      );
      const text = await r.text();
      setHtml(renderMarkdown(text));
    } catch {
      setHtml("<p>Could not load the document.</p>");
    }
    setLoading(false);
  }, []);

  // Boot: pick project from ?project= or the most recent one.
  useEffect(() => {
    (async () => {
      const rows = await loadProjects();
      const fromURL = new URL(window.location.href).searchParams.get("project");
      const pick = fromURL && rows.some((r) => r.id === fromURL) ? fromURL : rows[0]?.id;
      if (pick) {
        setProject(pick);
      }
    })();
  }, [loadProjects]);

  useEffect(() => {
    if (project) loadDocs(project);
  }, [project, loadDocs]);

  // Open the first doc automatically once the list arrives.
  useEffect(() => {
    if (!docPath && docs.length > 0 && project) openDoc(project, docs[0].path);
  }, [docs, docPath, project, openDoc]);

  // While a pipeline is running, poll so documents appear as agents finish.
  useEffect(() => {
    if (status !== "running") return;
    const t = setInterval(() => {
      loadDocs(project);
      loadProjects();
    }, 4000);
    return () => clearInterval(t);
  }, [status, project, loadDocs, loadProjects]);

  // Hand freshly rendered ```mermaid blocks to mermaid.js.
  useEffect(() => {
    if (!html.includes('class="mermaid"')) return;
    ensureMermaid().then((m) => {
      const mm = m as { run: (c: object) => Promise<void> } | null;
      mm?.run({ nodes: bodyRef.current?.querySelectorAll(".mermaid") ?? [] }).catch(() => {});
    });
  }, [html]);

  const stepIcon = (s: ManifestStep) => (s.ok ? "✓" : "✕");

  return (
    <div className="shell docs-app">
      <Chrome
        pageTitle="Documents"
        active="docs"
        project={cfg.label}
        company={company}
        pack={cfg.short}
        packGlyph={cfg.glyph}
        packCount={String(cfg.roles.length)}
      />
      <div className="docs-layout">
        <aside className="docs-side">
          <div className="docs-side-label">Projects</div>
          <div className="docs-projects">
            {projects.length === 0 && (
              <div className="docs-empty">No documents yet — run a discovery or a build.</div>
            )}
            {projects.map((p) => (
              <button
                key={p.id}
                className={`docs-proj${p.id === project ? " active" : ""}`}
                onClick={() => {
                  setProject(p.id);
                  setDocPath("");
                  setHtml("");
                }}
              >
                <span className="docs-proj-name">{p.id}</span>
                <span className="docs-proj-meta">
                  {p.docs} doc{p.docs === 1 ? "" : "s"}
                  {p.status === "running" && <span className="docs-live"> · live</span>}
                </span>
              </button>
            ))}
          </div>

          {project && (
            <>
              <div className="docs-side-label">Documents</div>
              <div className="docs-list">
                {docs.map((d) => (
                  <button
                    key={d.path}
                    className={`docs-item${d.path === docPath ? " active" : ""}`}
                    onClick={() => openDoc(project, d.path)}
                  >
                    <span className="docs-item-title">{d.title}</span>
                    <span className="docs-item-path">{d.path}</span>
                  </button>
                ))}
                {docs.length === 0 && <div className="docs-empty">No documents in this project yet.</div>}
              </div>
            </>
          )}

          {manifest && (
            <>
              <div className="docs-side-label">
                Pipeline {status === "running" ? "· running…" : `· ${status}`}
              </div>
              <div className="docs-steps">
                {manifest.steps?.map((s, i) => (
                  <div key={i} className={`docs-step${s.ok ? "" : " failed"}`}>
                    <span className="docs-step-ico">{stepIcon(s)}</span>
                    <span className="docs-step-agent">{s.agent}</span>
                    <span className="docs-step-action">{s.action}</span>
                  </div>
                ))}
                {status === "running" && (
                  <div className="docs-step pending">⏳ {manifest.current || "next agent working…"}</div>
                )}
              </div>
            </>
          )}
        </aside>

        <main className="docs-main">
          {loading && <div className="docs-empty">Loading…</div>}
          {!loading && !html && (
            <div className="docs-empty">
              Select a document. Discovery runs land here live — PRD, competitor research,
              compliance research, architecture, and the client reply draft.
            </div>
          )}
          {!loading && html && (
            <article
              ref={bodyRef}
              className="doc-body"
              dangerouslySetInnerHTML={{ __html: html }}
            />
          )}
        </main>
      </div>

      <style>{`
        .docs-layout { display: flex; gap: 0; height: calc(100vh - 56px); }
        .docs-side { width: 300px; min-width: 240px; border-right: 1px solid var(--line, rgba(128,128,128,.25)); overflow-y: auto; padding: 14px; }
        .docs-side-label { font-size: 11px; letter-spacing: .08em; text-transform: uppercase; opacity: .55; margin: 14px 0 6px; }
        .docs-proj, .docs-item { display: block; width: 100%; text-align: left; background: none; border: 1px solid transparent; border-radius: 8px; padding: 8px 10px; cursor: pointer; color: inherit; }
        .docs-proj:hover, .docs-item:hover { border-color: rgba(128,128,128,.3); }
        .docs-proj.active, .docs-item.active { border-color: rgba(128,128,128,.5); background: rgba(128,128,128,.08); }
        .docs-proj-name { font-weight: 600; display: block; }
        .docs-proj-meta, .docs-item-path { font-size: 11px; opacity: .55; display: block; }
        .docs-live { color: #2bb673; font-weight: 600; }
        .docs-item-title { display: block; font-size: 13px; }
        .docs-steps { font-size: 12px; }
        .docs-step { display: flex; gap: 6px; align-items: baseline; padding: 3px 2px; }
        .docs-step-ico { color: #2bb673; }
        .docs-step.failed .docs-step-ico { color: #d64545; }
        .docs-step-agent { font-weight: 600; }
        .docs-step-action { opacity: .6; }
        .docs-step.pending { opacity: .6; }
        .docs-empty { opacity: .55; font-size: 13px; padding: 12px 6px; }
        .docs-main { flex: 1; overflow-y: auto; padding: 28px 40px 80px; }
        .doc-body { max-width: 860px; margin: 0 auto; line-height: 1.65; font-size: 15px; }
        .doc-body h3 { font-size: 26px; margin: 8px 0 14px; }
        .doc-body h4 { font-size: 20px; margin: 22px 0 10px; }
        .doc-body h5 { font-size: 16px; margin: 18px 0 8px; }
        .doc-body table { border-collapse: collapse; margin: 14px 0; width: 100%; font-size: 13.5px; }
        .doc-body th, .doc-body td { border: 1px solid rgba(128,128,128,.35); padding: 6px 10px; text-align: left; vertical-align: top; }
        .doc-body th { background: rgba(128,128,128,.1); }
        .doc-body pre.codeblock { background: rgba(128,128,128,.1); border-radius: 8px; padding: 12px 14px; overflow-x: auto; font-size: 12.5px; }
        .doc-body blockquote { border-left: 3px solid rgba(128,128,128,.4); margin: 10px 0; padding: 2px 14px; opacity: .85; }
        .doc-body .mermaid { background: #fff; border: 1px solid rgba(128,128,128,.25); border-radius: 10px; padding: 16px; margin: 16px 0; display: flex; justify-content: center; overflow-x: auto; }
        .doc-body a { text-decoration: underline; }
        @media (max-width: 760px) { .docs-layout { flex-direction: column; height: auto; } .docs-side { width: auto; border-right: none; border-bottom: 1px solid rgba(128,128,128,.25); } }
      `}</style>
    </div>
  );
}
