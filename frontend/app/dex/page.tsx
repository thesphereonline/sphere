"use client";
import { useEffect, useState } from "react";

// Types
type Pool = {
  id: number;
  token_a: string;
  token_b: string;
};

export default function DexPage() {
  const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  // State hooks must be INSIDE the component
  const [pools, setPools] = useState<Pool[]>([]);
  const [poolId, setPoolId] = useState<number | null>(null);

  // Swap
  const [fromToken, setFromToken] = useState<string>("");
  const [amount, setAmount] = useState<string>("0");

  // Create Pool
  const [newTokenA, setNewTokenA] = useState("");
  const [newTokenB, setNewTokenB] = useState("");
  const [reserveA, setReserveA] = useState("1000");
  const [reserveB, setReserveB] = useState("1000");

  // Add Liquidity
  const [liquidityA, setLiquidityA] = useState("0");
  const [liquidityB, setLiquidityB] = useState("0");

  // Load pools
  const fetchPools = async () => {
    try {
      const res = await fetch(`${API}/dex/pools`);
      const data = await res.json();
      if (Array.isArray(data)) setPools(data);
    } catch (err) {
      console.error("Failed to fetch pools:", err);
      setPools([]);
    }
  };

  useEffect(() => {
    fetchPools();
  }, []);

  // Swap tokens
  const doSwap = async () => {
    if (!poolId) return alert("Please select a pool");
    try {
      const res = await fetch(`${API}/dex/pools/${poolId}/swap`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          fromToken,
          amountIn: parseFloat(amount),
          minOut: 1,
          trader: "alice",
        }),
      });
      if (!res.ok) return alert("Swap failed");
      const j = await res.json();
      alert("✅ amountOut: " + j.amountOut);
      fetchPools();
    } catch (err) {
      console.error("Swap error:", err);
      alert("❌ Swap failed");
    }
  };

  // Create a new pool
  const createPool = async () => {
    try {
      const res = await fetch(`${API}/dex/pools`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          tokenA: newTokenA,
          tokenB: newTokenB,
          reserveA: parseFloat(reserveA),
          reserveB: parseFloat(reserveB),
        }),
      });
      if (!res.ok) throw new Error("Failed to create pool");
      alert(`✅ Pool ${newTokenA}/${newTokenB} created!`);
      setNewTokenA("");
      setNewTokenB("");
      setReserveA("1000");
      setReserveB("1000");
      fetchPools();
    } catch (err) {
      console.error(err);
      alert("❌ Error creating pool");
    }
  };

  // Add liquidity
  const addLiquidity = async () => {
    if (!poolId) return alert("Select a pool first");
    try {
      const res = await fetch(`${API}/dex/pools/${poolId}/add`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          owner: "alice",
          amountA: liquidityA,
          amountB: liquidityB,
        }),
      });
      if (!res.ok) throw new Error("Add liquidity failed");
      alert("✅ Liquidity added!");
      fetchPools();
    } catch (err) {
      console.error(err);
      alert("❌ Error adding liquidity");
    }
  };

  return (
    <div className="p-6 bg-black text-orange-400 min-h-screen space-y-6">
      <h1 className="text-3xl font-bold">Sphere DEX</h1>

      {/* Pool Selector & Swap */}
      <div className="p-4 bg-orange-900/20 rounded-2xl space-y-3">
        <h2 className="text-xl">Swap Tokens</h2>
        <select
          className="w-full p-2 rounded bg-black border border-orange-400"
          onChange={(e) =>
            setPoolId(e.target.value ? Number(e.target.value) : null)
          }
        >
          <option value="">-- Select a pool --</option>
          {pools.map((p) => (
            <option key={p.id} value={p.id}>
              {p.token_a}/{p.token_b}
            </option>
          ))}
        </select>
        <input
          className="w-full p-2 rounded bg-black border border-orange-400"
          placeholder="From Token"
          value={fromToken}
          onChange={(e) => setFromToken(e.target.value)}
        />
        <input
          className="w-full p-2 rounded bg-black border border-orange-400"
          placeholder="Amount"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
        />
        <button
          onClick={doSwap}
          className="w-full bg-orange-400 hover:bg-orange-500 text-black p-2 rounded"
        >
          Swap
        </button>
      </div>

      {/* Create Pool */}
      <div className="p-4 bg-orange-900/20 rounded-2xl space-y-2">
        <h2 className="text-xl">Create New Pool</h2>
        <input
          className="w-full p-2 rounded bg-orange-800 text-orange-100"
          placeholder="Token A"
          value={newTokenA}
          onChange={(e) => setNewTokenA(e.target.value)}
        />
        <input
          className="w-full p-2 rounded bg-orange-800 text-orange-100"
          placeholder="Token B"
          value={newTokenB}
          onChange={(e) => setNewTokenB(e.target.value)}
        />
        <input
          className="w-full p-2 rounded bg-orange-800 text-orange-100"
          placeholder="Reserve A"
          value={reserveA}
          onChange={(e) => setReserveA(e.target.value)}
        />
        <input
          className="w-full p-2 rounded bg-orange-800 text-orange-100"
          placeholder="Reserve B"
          value={reserveB}
          onChange={(e) => setReserveB(e.target.value)}
        />
        <button
          onClick={createPool}
          className="w-full bg-orange-400 hover:bg-orange-500 text-black p-2 rounded"
        >
          Create Pool
        </button>
      </div>

      {/* Add Liquidity */}
      <div className="p-4 bg-orange-900/20 rounded-2xl space-y-2">
        <h2 className="text-xl">Add Liquidity</h2>
        <input
          className="w-full p-2 rounded bg-orange-800 text-orange-100"
          placeholder="Amount Token A"
          value={liquidityA}
          onChange={(e) => setLiquidityA(e.target.value)}
        />
        <input
          className="w-full p-2 rounded bg-orange-800 text-orange-100"
          placeholder="Amount Token B"
          value={liquidityB}
          onChange={(e) => setLiquidityB(e.target.value)}
        />
        <button
          onClick={addLiquidity}
          className="w-full bg-orange-400 hover:bg-orange-500 text-black p-2 rounded"
        >
          Add Liquidity
        </button>
      </div>
    </div>
  );
}
