import type { Metadata } from "next";
import { ConductorApp } from "@/components/ConductorApp";

export const metadata: Metadata = { title: "Agent House · Conductor" };

export default function ConductorPage() {
  return <ConductorApp />;
}
