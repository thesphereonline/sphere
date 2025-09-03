"use client";
import { useEffect, useState } from "react";

type Block = {
  id: number;
  hash: string;
  previous_hash: string;
  timestamp: string;
};

type Validator = {
  id: number;
  name: string;
  stake: number;
};

export default function DashboardPage() {
  const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const [blocks, setBlocks] = useState<Block[]>([]);
  const [validators, setValidators] = useState<Validator[]>([]);

  // Fetch blockchain data
  const fetchBlocks = async () => {
    try {
      const res = await fetch(`${API}/blocks`);
      if (!res.ok) throw new Error("Failed to load blocks");
      const data = await res.json();
      setBlocks(data);
    } catch (err) {
      console.error(err);
      setBlocks([]);
    }
  };

  const fetchValidators = async () => {
    try {
      const res = await fetch(`${API}/validators`);
      if (!res.ok) throw new Error("Failed to load validators");
      const data = await res.json();
      setValidators(data);
    } catch (err) {
      console.error(err);
      setValidators([]);
    }
  };

  useEffect(() => {
    fetchBlocks();
    fetchValidators();
  }, []);

  return (
    <div className="p-6 bg-black text-orange-400 min-h-screen space-y-6">
      <h1 className="text-3xl font-bold">Sphere Dashboard</h1>

      {/* Latest Blocks */}
      <div className="p-4 bg-orange-900/20 rounded-2xl space-y-3">
        <h2 className="text-xl">Latest Blocks</h2>
        <ul className="space-y-2">
          {blocks.length > 0 ? (
            blocks.map((b) => (
              <li
                key={b.id}
                className="p-2 bg-black border border-orange-500 rounded-lg"
              >
                <p className="font-mono text-sm">
                  <span className="font-bold">Hash:</span> {b.hash}
                </p>
                <p className="font-mono text-xs text-orange-300">
                  Prev: {b.previous_hash}
                </p>
                <p className="text-xs text-orange-200">
                  {new Date(b.timestamp).toLocaleString()}
                </p>
              </li>
            ))
          ) : (
            <p className="text-orange-300">No blocks found</p>
          )}
        </ul>
      </div>

      {/* Validators */}
      <div className="p-4 bg-orange-900/20 rounded-2xl space-y-3">
        <h2 className="text-xl">Validators</h2>
        <ul className="space-y-2">
          {validators.length > 0 ? (
            validators.map((v) => (
              <li
                key={v.id}
                className="p-2 bg-black border border-orange-500 rounded-lg"
              >
                <p>
                  <span className="font-bold">{v.name}</span> —{" "}
                  {v.stake.toLocaleString()} tokens staked
                </p>
              </li>
            ))
          ) : (
            <p className="text-orange-300">No validators found</p>
          )}
        </ul>
      </div>
    </div>
  );
}
