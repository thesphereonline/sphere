"use client";
import { useState, useEffect } from "react";

type Block = {
  Height: number;
  Hash: string;
  PrevHash: string;
  Timestamp: number;
  Validator: string;
};

export default function Dashboard() {
  const [blocks, setBlocks] = useState<Block[]>([]);

  useEffect(() => {
    fetch("https://loving-light-production.up.railway.app/blocks")
      .then((res) => res.json())
      .then((data: Block[]) => setBlocks(data));
  }, []);

  return (
    <div className="bg-black text-orange-500 min-h-screen p-6">
      <h1 className="text-4xl font-bold">Sphere Dashboard</h1>
      <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="bg-orange-900/20 rounded-2xl p-4">
          <h2 className="text-2xl">Latest Blocks</h2>
          <ul className="mt-2">
            {blocks.slice(-5).map((b) => (
              <li key={b.Hash} className="text-sm">
                #{b.Height} – {b.Hash.slice(0, 10)}...
              </li>
            ))}
          </ul>
        </div>
        <div className="bg-orange-900/20 rounded-2xl p-4">
          <h2 className="text-2xl">Your Wallet</h2>
          <p>Balance: 1000 $SPHERE</p>
          <button className="mt-2 px-4 py-2 bg-orange-600 text-black rounded-xl">
            Stake Now
          </button>
        </div>
      </div>
    </div>
  );
}
