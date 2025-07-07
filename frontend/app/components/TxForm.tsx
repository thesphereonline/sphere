"use client";

import { useState } from "react";
import { signTx } from "@/lib/wallet";
import { submitTx } from "@/app/actions/tx";

export function TxForm() {
  const [to, setTo] = useState("");
  const [amount, setAmount] = useState(1);
  const [nonce, setNonce] = useState(0);
  const [status, setStatus] = useState("");

  async function handleSubmit() {
    const stored = localStorage.getItem("wallet");
    if (!stored) {
      setStatus("❌ No wallet found");
      return;
    }

    const { publicKey, privateKey } = JSON.parse(stored);
    const tx = {
      from: publicKey,
      to,
      amount,
      nonce,
      timestamp: Date.now(),
    };

    // Hash transaction (excluding signature)
    const encoder = new TextEncoder();
    const json = JSON.stringify({ ...tx, signature: "" });
    const hashBuffer = await crypto.subtle.digest("SHA-256", encoder.encode(json));
    const hashHex = Array.from(new Uint8Array(hashBuffer))
      .map(b => b.toString(16).padStart(2, "0"))
      .join("");

    // Sign hash
    const signature = signTx(privateKey, hashHex);

    const txSigned = { ...tx, signature };

    try {
      const res = await submitTx(txSigned);
      setStatus("✅ " + res);
    } catch (e: unknown) {
      if (e instanceof Error) {
        setStatus("❌ " + e.message);
      } else {
        setStatus("❌ Unknown error");
      }
    }
  }

  return (
    <div className="p-4 bg-black text-white border border-orange-400 rounded-2xl shadow-xl mt-4">
      <h2 className="text-xl font-bold mb-2">🔁 Send Transaction</h2>
      <input
        className="mb-2 p-2 bg-gray-800 rounded w-full"
        value={to}
        onChange={(e) => setTo(e.target.value)}
        placeholder="Recipient address"
      />
      <input
        className="mb-2 p-2 bg-gray-800 rounded w-full"
        type="number"
        value={amount}
        onChange={(e) => setAmount(+e.target.value)}
        placeholder="Amount"
      />
      <input
        className="mb-2 p-2 bg-gray-800 rounded w-full"
        type="number"
        value={nonce}
        onChange={(e) => setNonce(+e.target.value)}
        placeholder="Nonce (e.g. 0)"
      />
      <button
        onClick={handleSubmit}
        className="bg-orange-500 hover:bg-orange-600 p-2 rounded text-black font-bold w-full"
      >
        Send
      </button>
      <div className="mt-2">{status}</div>
    </div>
  );
}
