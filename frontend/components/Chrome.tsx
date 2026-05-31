"use client";

import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { Caret, Search, Sun, Close } from "./icons";

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
  project = "Atlas Construction",
  pack = "EPC",
  packGlyph = "EP",
  packCount = "6 + 1",
}: {
  pageTitle: string;
  active: "conductor" | "mission" | "triggers" | "agents" | "lab";
  showCmdk?: boolean;
  /** Honest overrides — sample shells keep the EPC defaults. */
  project?: string;
  pack?: string;
  packGlyph?: string;
  packCount?: string;
}) {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);

  return (
    <header className="top">
      <div className="top-left">
        <div className="brand-mark" title="Agent House" />
        <div className="proj-switcher">
          <span className="proj-name">{project}</span>
          <span className="caret"><Caret /></span>
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
