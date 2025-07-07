"use server";

export async function submitTx(tx: any) {
  const res = await fetch("http://localhost:8080/submit-tx", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(tx),
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return await res.text();
}
