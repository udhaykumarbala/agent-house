import type { Metadata } from "next";
import { LabApp } from "@/components/LabApp";

export const metadata: Metadata = { title: "Agent House · Lab" };

export default function LabPage() {
  return <LabApp />;
}
