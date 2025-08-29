"use client";
import { useEffect, useState } from "react";
export default function DexPage() {
  const API = "https://loving-light-production.up.railway.app";
  const [pools, setPools] = useState<any[]>([]);
  const [poolId, setPoolId] = useState<number | null>(null);
  const [fromToken, setFromToken] = useState("");
  const [amount, setAmount] = useState<string>("0");

  useEffect(()=>{ fetch(`${API}/dex/pools`).then(r=>r.json()).then(setPools) },[])

  const doSwap = async () => {
    if (!poolId) return alert("select pool");
    const res = await fetch(`${API}/dex/pools/${poolId}/swap`, {
      method:"POST", headers:{"Content-Type":"application/json"},
      body: JSON.stringify({ fromToken, amountIn: amount, minOut: "1", trader: "alice" })
    });
    if(!res.ok) { alert("swap failed"); return }
    const j = await res.json(); alert("amountOut: "+j.amountOut);
  }

  return (
    <div className="p-6 bg-black text-orange-400">
      <h1 className="text-2xl">Sphere DEX</h1>
      <select onChange={(e)=>setPoolId(Number(e.target.value))}>
        <option value="">Select pool</option>
        {pools.map(p=> <option key={p.id} value={p.id}>{p.token_a}/{p.token_b}</option>)}
      </select>
      <input placeholder="from token" value={fromToken} onChange={e=>setFromToken(e.target.value)} />
      <input placeholder="amount" value={amount} onChange={e=>setAmount(e.target.value)} />
      <button onClick={doSwap}>Swap</button>
    </div>
  )
}
