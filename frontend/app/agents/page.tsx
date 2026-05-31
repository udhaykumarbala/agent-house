import type { Metadata } from "next";
import { AgentBuilderApp } from "@/components/AgentBuilderApp";

export const metadata: Metadata = { title: "Agent House · Agents" };

export default function AgentsPage() {
  return <AgentBuilderApp />;
}
