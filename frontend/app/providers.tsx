"use client";

import { ThemeProvider } from "next-themes";
import type { ReactNode } from "react";

/**
 * next-themes drives the reference's [data-mode="light|dark"] attribute.
 * Using the exact attribute the ported CSS keys off of means the theme
 * toggle is a single attribute flip — no JS recompute, no flash.
 */
export function Providers({ children }: { children: ReactNode }) {
  return (
    <ThemeProvider
      attribute="data-mode"
      defaultTheme="light"
      enableSystem={false}
      themes={["light", "dark"]}
      disableTransitionOnChange
    >
      {children}
    </ThemeProvider>
  );
}
