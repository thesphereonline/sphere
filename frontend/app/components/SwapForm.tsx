"use client";

import { useState, useEffect } from "react";

export function SwapForm() {
  const [tokenIn, setTokenIn] = useState("TOKENA");
  const [tokenOut, setTokenOut] = useState("TOKENB");
  const [amountIn, setAmountIn] = useState("");
  const [quote, setQuote] = useState<number | null>(null);
  const [status, setStatus] = useState("");

  useEffect(() => {
    if (!amountIn || isNaN(+amountIn)) return setQuote(null);
    const fetchQuote = async () => {
      const res = await fetch(`http://localhost:8080/quote?token_in=${tokenIn}&token_out=${tokenOut}&amount_in=${amountIn}`);
      const data = await res.json();
      setQuote(data.amount_out);
    };
    fetchQuote();
  }, [tokenIn, tokenOut, amountIn]);

  const handleSwap = async () => {
    const stored = localStorage.getItem("wallet");
    if (!stored) return setStatus("❌ No wallet found");
    const { publicKey, nonce } = JSON.parse(stored);

    const res = await fetch("http://localhost:8080/swap", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        from: publicKey,
        token_in: tokenIn,
        token_out: tokenOut,
        amount_in: Number(amountIn),
        nonce,
      }),
    });

    if (res.ok) {
      setStatus("✅ Swap successful");
    } else {
      const err = await res.text();
      setStatus(`❌ Swap failed: ${err}`);
    }
  };

  return (
    <div className="p-4 bg-black text-white border border-orange-400 rounded-2xl shadow-xl max-w-md w-full mx-auto mt-8">
      <h2 className="text-xl font-bold mb-4">🔁 Swap Tokens</h2>
      
      <label className="block mb-1">From Token</label>
      <select value={tokenIn} onChange={(e) => setTokenIn(e.target.value)} className="mb-2 p-2 w-full bg-zinc-800 rounded">
        <option value="TOKENA">TOKENA</option>
        <option value="TOKENB">TOKENB</option>
      </select>

      <label className="block mb-1">To Token</label>
      <select value={tokenOut} onChange={(e) => setTokenOut(e.target.value)} className="mb-2 p-2 w-full bg-zinc-800 rounded">
        <option value="TOKENB">TOKENB</option>
        <option value="TOKENA">TOKENA</option>
      </select>

      <label className="block mb-1">Amount In</label>
      <input
        type="number"
        className="mb-3 p-2 w-full bg-zinc-800 rounded"
        placeholder="Amount"
        value={amountIn}
        onChange={(e) => setAmountIn(e.target.value)}
      />

      {quote !== null && (
        <div className="text-sm text-gray-300 mb-3">
          Estimated Output: <span className="font-bold">{quote}</span> {tokenOut}
        </div>
      )}

      <button
        onClick={handleSwap}
        className="w-full bg-orange-500 hover:bg-orange-600 text-black font-bold py-2 px-4 rounded"
      >
        Swap
      </button>

      {status && <p className="mt-3 text-sm text-orange-400">{status}</p>}
    </div>
  );
}
