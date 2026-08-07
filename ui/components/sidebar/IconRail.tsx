"use client";

import { motion } from "framer-motion";
import type { LucideIcon } from "lucide-react";

interface Props {
  items: {
    title: string;
    icon: LucideIcon;
  }[];
}

export default function IconRail({ items }: Props) {
  return (
    <aside
      className="
flex
w-16
flex-col
items-center
gap-4
border-r
border-zinc-800
bg-zinc-950
py-6
"
    >
      {items.map((item) => {
        const Icon = item.icon;

        return (
          <motion.button
            key={item.title}
            whileHover={{
              scale: 1.1,
            }}
            className="
rounded-xl
p-3
text-zinc-400
hover:bg-blue-500/10
hover:text-blue-400
"
          >
            <Icon className="h-5 w-5" />
          </motion.button>
        );
      })}
    </aside>
  );
}
