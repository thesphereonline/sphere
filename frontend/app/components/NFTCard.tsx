export function NFTCard({ nft }: { nft: any }) {
  return (
    <div className="bg-zinc-900 p-3 rounded-xl border border-orange-400 shadow">
      <img
        src={nft.image_url}
        alt={nft.name}
        className="rounded mb-2 object-cover h-48 w-full"
      />
      <div className="text-white font-bold">{nft.name}</div>
      <div className="text-xs text-zinc-400 break-all">{nft.owner}</div>
    </div>
  );
}
