"use client";

import { useEffect } from "react";

// Static-export safe redirect (next/navigation redirect() is disallowed in export).
export default function Home() {
  useEffect(() => {
    window.location.replace("/mission/");
  }, []);
  return null;
}
