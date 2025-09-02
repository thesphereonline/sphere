// Add at the top with other state variables
const [newTokenA, setNewTokenA] = useState("");
const [newTokenB, setNewTokenB] = useState("");
const [reserveA, setReserveA] = useState("1000");
const [reserveB, setReserveB] = useState("1000");

// Add function to create a pool
const createPool = async () => {
  try {
    const res = await fetch(`${API}/dex/pools`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        tokenA: newTokenA,
        tokenB: newTokenB,
        reserveA: parseFloat(reserveA),
        reserveB: parseFloat(reserveB),
      }),
    });

    if (!res.ok) throw new Error("Failed to create pool");

    alert(`✅ Pool ${newTokenA}/${newTokenB} created!`);
    setNewTokenA("");
    setNewTokenB("");
    setReserveA("1000");
    setReserveB("1000");
    // Refresh pools
    const updatedPools = await fetch(`${API}/dex/pools`).then(r => r.json());
    setPools(updatedPools);
  } catch (err) {
    console.error(err);
    alert("❌ Error creating pool");
  }
};
