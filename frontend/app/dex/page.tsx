"use client";
import { useEffect, useState } from "react";

// Define a type for a pool
type Pool = {
  id: number;
  token_a: string;
  token_b: string;
};

export default function DexPage() {
  const API = "https://loving-light-production.up.railway.app";
  const [pools, setPools] = useState<Pool[]>([]);
  const [poolId, setPoolId] = useState<number | null>(null);
  const [fromToken, setFromToken] = useState<string>("");
  const [amount, setAmount] = useState<string>("0");

  // Load pools on mount
  useEffect(() => {
    fetch(`${API}/dex/pools`)
      .then(r => r.json())
      .then(data => {
        if (Array.isArray(data)) {
          setPools(data);
        } else {
          console.error("Unexpected pools response:", data);
          setPools([]);
        }
      })
      .catch(err => {
        console.error("Failed to fetch pools:", err);
        setPools([]);
      });
  }, []);

  // Perform a swap
  const doSwap = async () => {
    if (!poolId) {
      alert("Please select a pool");
      return;
    }
    try {
      const res = await fetch(`${API}/dex/pools/${poolId}/swap`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          fromToken,
          amountIn: amount,
          minOut: "1",
          trader: "alice",
        }),
      });
      if (!res.ok) {
        alert("Swap failed");
        return;
      }
      const j = await res.json();
      alert("amountOut: " + j.amountOut);
    } catch (err) {
      console.error("Swap error:", err);
      alert("Swap failed due to network error");
    }
  };

  return (
    <div className="p-6 bg-black text-orange-400 min-h-screen">
      <h1 className="text-2xl mb-4">Sphere DEX</h1>

      <div className="space-y-4">
        {/* Pool Selector */}
        <div>
          <label className="block mb-1">Select Pool</label>
          <select
            className="bg-black border border-orange-400 p-2 rounded w-full"
            onChange={(e) =>
              setPoolId(e.target.value ? Number(e.target.value) : null)
            }
          >
            <option value="">-- Choose a pool --</option>
            {pools.length === 0 && <option disabled>No pools found</option>}
            {Array.isArray(pools) &&
              pools.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.token_a}/{p.token_b}
                </option>
              ))}
          </select>
        </div>

        {/* From Token Input */}
        <div>
          <label className="block mb-1">From Token</label>
          <input
            className="bg-black border border-orange-400 p-2 rounded w-full"
            placeholder="e.g. TOKENA"
            value={fromToken}
            onChange={(e) => setFromToken(e.target.value)}
          />
        </div>

        {/* Amount Input */}
        <div>
          <label className="block mb-1">Amount</label>
          <input
            className="bg-black border border-orange-400 p-2 rounded w-full"
            placeholder="0"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
          />
        </div>

        {/* Swap Button */}
        <button
          onClick={doSwap}
          className="bg-orange-400 text-black px-4 py-2 rounded hover:bg-orange-500"
        >
          Swap
        </button>
      </div>
    </div>
  );
}
