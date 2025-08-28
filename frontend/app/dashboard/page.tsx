"use client";

import { useState, useEffect } from "react";

type Block = {
  Height: number;
  Hash: string;
  PrevHash: string;
  Timestamp: number;
  Validator: string;
  Transactions: {
    From: string;
    To: string;
    Amount: number;
    Fee: number;
    Data: string;
    Sig: string;
  }[];
};

export default function Dashboard() {
  const [blocks, setBlocks] = useState<Block[]>([]);
  const [from, setFrom] = useState("alice");
  const [to, setTo] = useState("bob");
  const [amount, setAmount] = useState(100);
  const [data, setData] = useState("Test transaction");
  const [sig, setSig] = useState("test-signature"); // plain string now
  const API_URL = "https://loving-light-production.up.railway.app";

  // Fetch blocks
  const fetchBlocks = async () => {
    try {
      const res = await fetch(`${API_URL}/blocks`);
      if (!res.ok) throw new Error("Failed to fetch blocks");
      const data: Block[] = await res.json();
      setBlocks(data);
    } catch (err) {
      console.error("Error fetching blocks:", err);
    }
  };

  useEffect(() => {
    fetchBlocks();
  }, []);

  // Submit a new tx -> block
  const handleAddBlock = async () => {
    try {
      const res = await fetch(`${API_URL}/tx`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          From: from,
          To: to,
          Amount: amount,
          Fee: 1,
          Data: data,
          Sig: sig, // send plain string
        }),
      });

      if (!res.ok) throw new Error("Failed to add block");
      await fetchBlocks(); // refresh after adding
      alert("✅ Block added!");
    } catch (err) {
      console.error(err);
      alert("❌ Error adding block. Check backend logs.");
    }
  };

  return (
    <div className="bg-black text-orange-500 min-h-screen p-6">
      <h1 className="text-4xl font-bold">Sphere Dashboard</h1>

      {/* Grid layout */}
      <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
        
        {/* Block List */}
        <div className="bg-orange-900/20 rounded-2xl p-4">
          <h2 className="text-2xl">Latest Blocks</h2>
          <ul className="mt-2 space-y-1">
            {blocks.slice(-5).map((b) => (
              <li key={b.Hash} className="text-sm">
                <span className="font-semibold">#{b.Height}</span> –{" "}
                {b.Hash.slice(0, 12)}...
              </li>
            ))}
          </ul>
        </div>

        {/* Test Block Form */}
        <div className="bg-orange-900/20 rounded-2xl p-4 flex flex-col gap-2">
          <h2 className="text-2xl">Add Test Transaction</h2>
          <input
            type="text"
            value={from}
            onChange={(e) => setFrom(e.target.value)}
            placeholder="From"
            className="p-2 rounded bg-orange-800 text-orange-100"
          />
          <input
            type="text"
            value={to}
            onChange={(e) => setTo(e.target.value)}
            placeholder="To"
            className="p-2 rounded bg-orange-800 text-orange-100"
          />
          <input
            type="number"
            value={amount}
            onChange={(e) => setAmount(Number(e.target.value))}
            placeholder="Amount"
            className="p-2 rounded bg-orange-800 text-orange-100"
          />
          <input
            type="text"
            value={data}
            onChange={(e) => setData(e.target.value)}
            placeholder="Data"
            className="p-2 rounded bg-orange-800 text-orange-100"
          />
          <input
            type="text"
            value={sig}
            onChange={(e) => setSig(e.target.value)}
            placeholder="Signature"
            className="p-2 rounded bg-orange-800 text-orange-100"
          />
          <button
            onClick={handleAddBlock}
            className="mt-2 px-4 py-2 bg-orange-600 text-black rounded-xl hover:bg-orange-400"
          >
            Add Test Block
          </button>
        </div>
      </div>
    </div>
  );
}
