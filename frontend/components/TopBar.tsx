"use client";

import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { Caret, Search, Bell, Sun } from "./icons";

const CRUMBS = [
  ["conductor", "/conductor"],
  ["mission", "/mission"],
  ["triggers", "/triggers"],
  ["agents", "/agents"],
];

export function TopBar({
  project,
  site,
  pack,
  packCount,
  pendingCount,
  user,
}: {
  project: string;
  site: string;
  pack: string;
  packCount: string;
  pendingCount: number;
  user: string;
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
          <span className="proj-name">{site}</span>
          <span className="caret"><Caret /></span>
        </div>
        <div style={{ marginLeft: 10 }} className="crumbs">
          {CRUMBS.map(([label, href], i) => (
            <span key={label} style={{ display: "contents" }}>
              <a href={href} className={label === "mission" ? "active" : ""}>
                {label}
              </a>
              {i < CRUMBS.length - 1 && <span className="sep">·</span>}
            </span>
          ))}
        </div>
      </div>
      <div className="top-right">
        <span className="pack-pill">
          <span className="glyph">{pack.slice(0, 2).toUpperCase()}</span>
          {pack}
          <span className="count">· {packCount}</span>
        </span>
        <button className="cmdk-trigger">
          <Search />
          <span className="ph">Find anything…</span>
          <span className="kbd">⌘K</span>
        </button>
        <button className="icon-btn inbox-bell" title="Pending checkpoints">
          <Bell />
          {pendingCount > 0 && <span className="badge">{pendingCount}</span>}
        </button>
        <button
          className="icon-btn"
          title="Toggle theme"
          onClick={() =>
            setTheme((mounted ? theme : "light") === "dark" ? "light" : "dark")
          }
        >
          <Sun />
        </button>
        <div className="user" title={user}>
          {user.slice(0, 2).toUpperCase()}
        </div>
      </div>
    </header>
  );
}
