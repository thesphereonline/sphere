"use client";

import { useEffect, useState } from "react";
import { generateWallet } from "@/lib/wallet";

export function WalletCard() {
  const [wallet, setWallet] = useState<{ address: string } | null>(null);

  useEffect(() => {
    let stored = localStorage.getItem("wallet");
    if (!stored) {
      const w = generateWallet();
      localStorage.setItem("wallet", JSON.stringify(w));
      stored = JSON.stringify(w);
    }
    const parsed = JSON.parse(stored);
    setWallet({ address: parsed.publicKey });
  }, []);

  return (
    <div className="rounded-2xl shadow-xl p-4 bg-black text-white border border-orange-400">
      <h2 className="text-xl font-bold mb-2">💼 Wallet</h2>
      {wallet ? (
        <div>
          <div><b>Address:</b> {wallet.address.slice(0, 32)}...</div>
        </div>
      ) : (
        <div>Loading...</div>
      )}
    </div>
  );
}
