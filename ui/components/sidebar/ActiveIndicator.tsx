"use client";

import { motion } from "framer-motion";

export default function ActiveIndicator() {
  return (
    <motion.div
      layoutId="active"
      className="
absolute
left-0
h-8
w-1
rounded-r
bg-blue-500
"
    />
  );
}
