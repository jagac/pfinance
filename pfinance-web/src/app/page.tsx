import { Dashboard } from "@/components/Dashboard/Dashboard";
import { Sidebar } from "@/components/Sidebar/Sidebar";

export default function Home() {
  return (
    <main className="grid gap-4 p-4 grid-cols-[220px_1fr] h-screen">
      <Sidebar />
      <Dashboard />
    </main>
  );
}
