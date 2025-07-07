import { WalletCard } from "./components/WalletCard";
import { TxForm } from "./components/TxForm";

export default function Home() {
  return (
    <main className="max-w-xl mx-auto py-10">
      <h1 className="text-4xl font-bold text-orange-500 mb-6">🪐 Sphere Protocol Wallet</h1>
      <WalletCard />
      <TxForm />
    </main>
  );
}
