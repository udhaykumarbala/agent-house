import type { Metadata } from "next";
import { TriggersApp } from "@/components/TriggersApp";

export const metadata: Metadata = { title: "Agent House · Triggers" };

export default function TriggersPage() {
  return <TriggersApp />;
}
