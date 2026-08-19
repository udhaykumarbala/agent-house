"use client";

import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { Caret, Search, Sun, Close } from "./icons";
import { COMPANIES, setCompany, type Company } from "@/lib/company";

const CRUMBS: [string, string][] = [
  ["conductor", "/conductor"],
  ["mission", "/mission"],
  ["triggers", "/triggers"],
  ["agents", "/agents"],
  ["lab", "/lab"],
];

/** Shared topbar for Conductor / Triggers / AgentBuilder (no inbox bell). */
export function Chrome({
  pageTitle,
  active,
  showCmdk = false,
  project = "Alredaa",
  company,
  pack = "EPC",
  packGlyph = "EP",
  packCount = "6 + 1",
}: {
  pageTitle: string;
  active: "conductor" | "mission" | "triggers" | "agents" | "lab";
  showCmdk?: boolean;
  /** Honest overrides — sample shells keep the EPC defaults. */
  project?: string;
  /** When set, the project name becomes a live company switcher. */
  company?: Company;
  pack?: string;
  packGlyph?: string;
  packCount?: string;
}) {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const [switcherOpen, setSwitcherOpen] = useState(false);
  useEffect(() => setMounted(true), []);

  // Switching company reloads with ?company= so every component re-reads it.
  const pickCompany = (c: Company) => {
    setCompany(c);
    const url = new URL(window.location.href);
    url.searchParams.set("company", c);
    window.location.href = url.toString();
  };

  return (
    <header className="top">
      <div className="top-left">
        <div className="brand-mark" title="Agent House" />
        <div
          className={`proj-switcher${company ? " switchable" : ""}`}
          onClick={() => company && setSwitcherOpen((o) => !o)}
        >
          <span className="proj-name">{project}</span>
          <span className="caret"><Caret /></span>
          {company && switcherOpen && (
            <div className="company-menu" onClick={(e) => e.stopPropagation()}>
              <div className="company-menu-label">Switch company</div>
              {Object.values(COMPANIES).map((c) => (
                <button
                  key={c.id}
                  className={`company-opt${c.id === company ? " active" : ""}`}
                  onClick={() => pickCompany(c.id)}
                >
                  <span className="company-glyph">{c.glyph}</span>
                  <span className="company-meta">
                    <span className="company-label">{c.label}</span>
                    <span className="company-tag">{c.tagline}</span>
                  </span>
                  {c.id === company && <span className="company-dot" />}
                </button>
              ))}
            </div>
          )}
        </div>
        <span className="proj-sep">/</span>
        <div className="proj-switcher">
          <span className="proj-name">{pageTitle}</span>
          <span className="caret"><Caret /></span>
        </div>
        <div className="crumbs">
          {CRUMBS.map(([label, href], i) => (
            <span key={label} style={{ display: "contents" }}>
              <a href={href} className={label === active ? "active" : ""}>
                {label}
              </a>
              {i < CRUMBS.length - 1 && <span className="sep">·</span>}
            </span>
          ))}
        </div>
      </div>
      <div className="top-right">
        <span className="pack-pill">
          <span className="glyph">{packGlyph}</span>
          {pack}
          <span className="count">· {packCount}</span>
        </span>
        {showCmdk && (
          <button className="cmdk-trigger">
            <Search />
            <span className="ph">Find anything…</span>
            <span className="kbd">⌘K</span>
          </button>
        )}
        <button
          className="icon-btn"
          title="Toggle theme"
          onClick={() =>
            setTheme((mounted ? theme : "light") === "dark" ? "light" : "dark")
          }
        >
          <Sun />
        </button>
        <div className="user" title="Marcus L.">ML</div>
      </div>
    </header>
  );
}

export { Close };
