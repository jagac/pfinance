import AddAssetForm from "@/components/Form/AssetForm";
import { Sidebar } from "@/components/Sidebar/Sidebar";

export default function AssetFormPage() {
  return (
    <main className="grid gap-4 p-4 grid-cols-[220px_1fr] h-screen">
      <Sidebar />
      <AddAssetForm />
    </main>
  );
}
