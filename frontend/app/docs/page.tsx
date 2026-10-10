import type { Metadata } from "next";
import { DocsApp } from "@/components/DocsApp";

export const metadata: Metadata = { title: "Agent House · Documents" };

export default function DocsPage() {
  return <DocsApp />;
}
