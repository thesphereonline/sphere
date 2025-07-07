"use client";

import { useState } from "react";

export function MintNFT() {
  const [name, setName] = useState("");
  const [imageURL, setImageURL] = useState("");
  const [status, setStatus] = useState("");

  async function handleMint() {
    const stored = localStorage.getItem("wallet");
    if (!stored) return setStatus("No wallet found");

    const { publicKey } = JSON.parse(stored);

    const res = await fetch("http://localhost:8080/mint-nft", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        owner: publicKey,
        name,
        image_url: imageURL,
        metadata: { description: "Minted via Sphere" },
      }),
    });

    if (res.ok) {
      setStatus("✅ NFT minted");
      setName("");
      setImageURL("");
    } else {
      setStatus("❌ Mint failed");
    }
  }

  return (
    <div className="p-4 bg-black text-white border border-orange-400 rounded-2xl shadow-xl mt-6">
      <h2 className="text-xl font-bold mb-2">🎨 Mint NFT</h2>
      <input
        className="mb-2 p-2 bg-gray-800 rounded w-full"
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Name"
      />
      <input
        className="mb-2 p-2 bg-gray-800 rounded w-full"
        value={imageURL}
        onChange={(e) => setImageURL(e.target.value)}
        placeholder="Image URL"
      />
      <button
        onClick={handleMint}
        className="bg-orange-500 hover:bg-orange-600 p-2 rounded text-black font-bold w-full"
      >
        Mint NFT
      </button>
      <div className="mt-2">{status}</div>
    </div>
  );
}
