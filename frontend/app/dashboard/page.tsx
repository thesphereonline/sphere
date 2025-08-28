// app/test-tx/page.tsx
"use client";

import { useState } from "react";

export default function TestTxPage() {
  const [from, setFrom] = useState("alice");
  const [to, setTo] = useState("bob");
  const [amount, setAmount] = useState(100);
  const [fee, setFee] = useState(1);
  const [data, setData] = useState("Test transaction");
  const [sig, setSig] = useState("abc123"); // plain string now
  const [response, setResponse] = useState<any>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const res = await fetch("https://loving-light-production.up.railway.app/tx", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ From: from, To: to, Amount: amount, Fee: fee, Data: data, Sig: sig }),
    });

    const json = await res.json();
    setResponse(json);
  };

  return (
    <div className="flex flex-col items-center p-8">
      <h1 className="text-2xl font-bold text-orange-500">🚀 Add Test Transaction</h1>

      <form onSubmit={handleSubmit} className="mt-6 space-y-4 w-full max-w-md">
        <input
          type="text"
          placeholder="From"
          value={from}
          onChange={(e) => setFrom(e.target.value)}
          className="w-full p-2 border rounded"
        />
        <input
          type="text"
          placeholder="To"
          value={to}
          onChange={(e) => setTo(e.target.value)}
          className="w-full p-2 border rounded"
        />
        <input
          type="number"
          placeholder="Amount"
          value={amount}
          onChange={(e) => setAmount(Number(e.target.value))}
          className="w-full p-2 border rounded"
        />
        <input
          type="number"
          placeholder="Fee"
          value={fee}
          onChange={(e) => setFee(Number(e.target.value))}
          className="w-full p-2 border rounded"
        />
        <input
          type="text"
          placeholder="Data"
          value={data}
          onChange={(e) => setData(e.target.value)}
          className="w-full p-2 border rounded"
        />
        <input
          type="text"
          placeholder="Signature"
          value={sig}
          onChange={(e) => setSig(e.target.value)}
          className="w-full p-2 border rounded"
        />
        <button
          type="submit"
          className="w-full bg-orange-500 text-white p-2 rounded hover:bg-orange-600"
        >
          Submit Transaction
        </button>
      </form>

      {response && (
        <div className="mt-6 w-full max-w-md p-4 border rounded bg-black text-white">
          <h2 className="font-bold">✅ Block Created:</h2>
          <pre className="text-xs whitespace-pre-wrap">{JSON.stringify(response, null, 2)}</pre>
        </div>
      )}
    </div>
  );
}
