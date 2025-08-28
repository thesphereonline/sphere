// app/page.tsx
"use client";

import { motion } from "framer-motion";

export default function Home() {
  return (
    <main className="bg-black text-orange-500 min-h-screen flex flex-col items-center justify-center">
      <motion.h1
        className="text-6xl font-bold"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        The Sphere
      </motion.h1>
      <p className="mt-4 text-xl text-orange-300">
        A next-gen blockchain for entrepreneurs
      </p>
      <a
        href="/dashboard"
        className="mt-6 px-6 py-3 bg-orange-600 text-black rounded-2xl shadow-lg hover:bg-orange-400"
      >
        Launch App
      </a>
    </main>
  );
}
