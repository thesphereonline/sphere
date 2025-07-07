"use client"

import { useEffect, useState } from "react"

export default function NFTsPage() {
  const [nfts, setNFTs] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    const stored = localStorage.getItem("wallet")
    if (!stored) {
      setError("No wallet found")
      setLoading(false)
      return
    }

    const { publicKey } = JSON.parse(stored)

    fetch(`http://localhost:8080/nfts?owner=${publicKey}`)
      .then((res) => {
        if (!res.ok) throw new Error("Failed to load NFTs")
        return res.json()
      })
      .then((data) => setNFTs(data ?? []))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="p-4 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-4 text-white">🎨 My NFTs</h1>

      {loading && <p className="text-gray-300">Loading...</p>}
      {error && <p className="text-red-500">❌ {error}</p>}
      {!loading && nfts.length === 0 && !error && (
        <p className="text-gray-400">You don't own any NFTs yet.</p>
      )}

      <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 mt-4">
        {nfts.map((nft, idx) => (
          <div
            key={idx}
            className="bg-zinc-900 border border-orange-500 rounded-xl p-4 shadow-lg text-white"
          >
            <img
              src={nft.image_url}
              alt={nft.name}
              className="w-full h-48 object-cover rounded mb-2 border border-zinc-700"
            />
            <h3 className="text-lg font-semibold">{nft.name}</h3>
            <p className="text-sm text-zinc-400">
              {nft.metadata?.description || "No description"}
            </p>
          </div>
        ))}
      </div>
    </div>
  )
}
