"use client";

import AnimatedPacket from "./AnimatedPacket";
import Connection from "./Connection";

export default function SocketMigration() {
  return (
    <div
      className="
rounded-xl
border
border-zinc-800
bg-zinc-950
p-6
"
    >
      <h2
        className="
mb-6
text-sm
font-semibold
text-white
"
      >
        TCP Socket Migration
      </h2>

      <div
        className="
flex
items-center
gap-4
"
      >
        <div
          className="
rounded-lg
bg-zinc-900
px-5
py-4
text-sm
text-white
"
        >
          Worker A
        </div>

        <div
          className="
relative
flex
flex-1
items-center
"
        >
          <Connection label="TCP State Transfer" />

          <AnimatedPacket />
        </div>

        <div
          className="
rounded-lg
bg-emerald-950
border
border-emerald-800
px-5
py-4
text-sm
text-emerald-400
"
        >
          Worker B
        </div>
      </div>

      <div
        className="
mt-6
grid
grid-cols-3
gap-4
text-xs
"
      >
        <div
          className="
rounded-lg
bg-zinc-900
p-3
"
        >
          Socket:
          <br />
          TCP-1024
        </div>

        <div
          className="
rounded-lg
bg-zinc-900
p-3
"
        >
          State:
          <br />
          Migrating
        </div>

        <div
          className="
rounded-lg
bg-zinc-900
p-3
"
        >
          Recovery:
          <br />
          18ms
        </div>
      </div>
    </div>
  );
}
