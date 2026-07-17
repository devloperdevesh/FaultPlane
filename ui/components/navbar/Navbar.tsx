"use client";

import { Bell, Cpu, Globe, Search } from "lucide-react";

import ChaosButton from "./ChaosButton";
import StatusIndicator from "./StatusIndicator";

export default function Navbar() {
  return (
    <header
      className="
        flex
        h-16
        items-center
        justify-between
        border-b
        border-zinc-800
        bg-zinc-950
        px-8
        text-white
      "
    >
      {/* Product Identity + Search */}
      <div className="flex items-center gap-6">
        <div>
          <h1 className="text-lg font-semibold">FaultPlane</h1>

          <p className="text-xs text-zinc-500">Operations Control Plane</p>
        </div>

        <div
          className="
            flex
            items-center
            gap-2
            rounded-lg
            border
            border-zinc-800
            bg-zinc-900
            px-3
            py-2
          "
        >
          <Search className="h-4 w-4 text-zinc-400" />

          <input
            placeholder="Search runtime..."
            className="
              w-40
              bg-transparent
              text-sm
              outline-none
              placeholder:text-zinc-600
            "
          />
        </div>
      </div>

      {/* Runtime Controls */}
      <div className="flex items-center gap-3">
        <StatusIndicator />

        <div
          className="
            flex
            items-center
            gap-2
            rounded-lg
            border
            border-zinc-800
            bg-zinc-900
            px-3
            py-2
          "
        >
          <Globe className="h-4 w-4 text-blue-400" />

          <div>
            <p className="text-xs text-zinc-500">Region</p>

            <p className="text-sm">us-east-1</p>
          </div>
        </div>

        <div
          className="
            flex
            items-center
            gap-2
            rounded-lg
            border
            border-zinc-800
            bg-zinc-900
            px-3
            py-2
          "
        >
          <Cpu className="h-4 w-4 text-purple-400" />

          <div>
            <p className="text-xs text-zinc-500">Runtime</p>

            <p className="text-sm">v0.1.0</p>
          </div>
        </div>

        <Bell
          className="
            h-5
            w-5
            text-zinc-400
            hover:text-white
            cursor-pointer
          "
        />

        <ChaosButton />
      </div>
    </header>
  );
}
