"use client";

import { motion } from "framer-motion";

interface Props {
  collapsed: boolean;
}

export default function SidebarCollapse({ collapsed }: Props) {
  return (
    <motion.div
      animate={{
        width: collapsed ? 64 : 280,
      }}
      transition={{
        type: "spring",
        stiffness: 200,
        damping: 25,
      }}
      className="
h-full
border-r
border-zinc-800
bg-zinc-950
"
    />
  );
}
