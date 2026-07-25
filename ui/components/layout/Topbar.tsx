"use client";

import { Search, Bell, Command, ChevronDown } from "lucide-react";

export default function Topbar() {
  return (
    <header
      className="
h-16

border-b

border-white/[0.06]

bg-[#060914]/80

backdrop-blur-xl

flex

items-center

justify-between

px-6

sticky

top-0

z-30

"
    >
      {/* LEFT */}

      <div className="flex items-center gap-4">
        <div
          className="
text-sm
font-semibold
text-white
"
        >
          Operations Console
        </div>

        <div
          className="
hidden
md:flex

items-center

gap-2

px-3

py-1.5

rounded-lg

bg-white/[0.04]

border

border-white/10

text-xs

text-slate-400

"
        >
          Production
          <ChevronDown size={14} />
        </div>
      </div>

      {/* RIGHT */}

      <div
        className="
flex
items-center
gap-3
"
      >
        <button
          className="
flex
items-center
gap-2

px-3
py-2

rounded-lg

bg-white/[0.04]

border

border-white/10

text-xs

text-slate-400

"
        >
          <Search size={14} />
          Search
          <span className="text-slate-600">⌘K</span>
        </button>

        <button
          className="
relative

w-9
h-9

rounded-lg

border

border-white/10

flex

items-center

justify-center

hover:bg-white/5

"
        >
          <Bell size={16} />

          <span
            className="
absolute

top-1

right-1

w-2

h-2

rounded-full

bg-red-500

"
          />
        </button>
      </div>
    </header>
  );
}
