"use client";

import { motion } from "framer-motion";

import TenantStatusBadge from "./TenantStatusBadge";
import QuotaProgress from "./QuotaProgress";

const tenants = [
  {
    name: "Acme AI",
    cpu: 72,
    memory: 64,
    requests: "120K",
    status: "healthy",
  },

  {
    name: "Nova Labs",
    cpu: 91,
    memory: 86,
    requests: "240K",
    status: "warning",
  },

  {
    name: "Edge Systems",
    cpu: 35,
    memory: 42,
    requests: "50K",
    status: "healthy",
  },
] as const;

export default function TenantTable() {
  return (
    <section
      className="
rounded-2xl
border
border-white/10
bg-zinc-950/80
backdrop-blur-xl
p-6
"
    >
      <div className="mb-6">
        <h2
          className="
text-lg
font-semibold
text-white
"
        >
          Tenant Resource Overview
        </h2>

        <p
          className="
mt-1
text-sm
text-zinc-500
"
        >
          Multi tenant isolation and quota governance
        </p>
      </div>

      <div className="space-y-4">
        {tenants.map((tenant, index) => (
          <motion.div
            key={tenant.name}
            initial={{
              opacity: 0,
              y: 20,
            }}
            animate={{
              opacity: 1,
              y: 0,
            }}
            transition={{
              delay: index * 0.1,
            }}
            whileHover={{
              scale: 1.01,
            }}
            className="
rounded-xl
border
border-white/10
bg-zinc-900/60
p-5
"
          >
            <div
              className="
flex
justify-between
items-center
mb-5
"
            >
              <h3
                className="
text-white
font-medium
"
              >
                {tenant.name}
              </h3>

              <TenantStatusBadge status={tenant.status} />
            </div>

            <div
              className="
grid
gap-5
md:grid-cols-3
"
            >
              <QuotaProgress label="CPU" value={tenant.cpu} />

              <QuotaProgress label="Memory" value={tenant.memory} />

              <div>
                <p
                  className="
text-xs
text-zinc-500
"
                >
                  Requests
                </p>

                <p
                  className="
mt-2
font-mono
text-white
"
                >
                  {tenant.requests}
                </p>
              </div>
            </div>
          </motion.div>
        ))}
      </div>
    </section>
  );
}
