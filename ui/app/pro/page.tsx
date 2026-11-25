"use client";

import Link from "next/link";
import { useState } from "react";
import { ArrowRight, BadgeCheck, Check, ChevronRight, LockKeyhole, Shield, Sparkles, Zap } from "lucide-react";

const proFeatures = [
  ["Advanced recovery policy orchestration", "/recovery-rules"],
  ["Failure Lab with controlled scenario workflows", "/failure-lab"],
  ["Checkpoint analysis and state-recovery tooling", "/checkpoints"],
  ["Time Travel and execution-history inspection", "/experimental/time-travel"],
  ["Speculative execution-path analysis", "/experimental/speculative"],
  ["Advanced telemetry and incident correlation", "/telemetry"],
  ["FinOps visibility for workload interruption costs", "/experimental/finops"],
  ["Team controls, audit capabilities and priority support", "/team"],
];

export default function ProPage() {
  const [cycle, setCycle] = useState<"monthly" | "annual">("monthly");
  return <main className="fp-pro-page">
    <div className="fp-pro-hero">
      <div className="fp-pro-eyebrow"><Sparkles size={14}/> FAULTPLANE PRO <span>PRIVATE · PROPRIETARY</span></div>
      <h1>Resilience for workloads<br/><em>that cannot afford to stop.</em></h1>
      <p>Go beyond visibility. Add advanced recovery workflows, deeper execution insight and operational controls for long-running AI workloads.</p>
      <div className="fp-pro-hero-actions"><a href="#plans" className="fp-pro-primary">Explore Pro plans <ArrowRight size={16}/></a><Link href="/failure-lab" className="fp-pro-secondary">Try interactive demo</Link><Link href="/dashboard" className="fp-pro-secondary">Back to dashboard</Link></div>
      <div className="fp-pro-proof-row"><span><Shield size={14}/> Private Pro distribution</span><span><LockKeyhole size={14}/> Server-verified access required</span><span><BadgeCheck size={14}/> Feature availability depends on plan</span></div>
    </div>

    <section className="fp-pro-feature-layout">
      <div className="fp-pro-feature-copy"><div className="fp-pro-kicker">BUILT FOR THE RUNTIME BOUNDARY</div><h2>Turn infrastructure signals into a controlled recovery workflow.</h2><p>FaultPlane Pro is designed for teams operating long-running agents, inference services and distributed workloads. It builds on the FaultPlane runtime concept: signal, state, policy, enforcement, recovery and telemetry.</p>
        <div className="fp-pro-loop">{["Signal","State","Policy","Enforce","Recover","Telemetry"].map((x,i)=><div key={x}><b>0{i+1}</b><span>{x}</span>{i!==5&&<ChevronRight size={14}/>}</div>)}</div>
      </div>
      <div className="fp-pro-feature-card"><div className="fp-pro-card-top"><span><Zap size={15}/> PRO CAPABILITIES</span><span className="fp-pro-lock"><LockKeyhole size={12}/> Proprietary</span></div><h3>Advanced resilience toolkit</h3><ul>{proFeatures.map(([f,href])=><li key={f}><Check size={15}/><Link href={href} className="fp-pro-feature-link">{f}<ArrowRight size={12}/></Link></li>)}</ul><p className="fp-pro-small-note">Feature availability is subject to implementation, plan entitlements and connected runtime capabilities. Experimental modules are labelled in-product.</p></div>
    </section>

    <section id="plans" className="fp-pro-plans"><div className="fp-pro-section-heading"><div><div className="fp-pro-kicker">PLANS & ACCESS</div><h2>Choose how your team runs FaultPlane.</h2><p>Core remains the open-source foundation. Pro is a separate private, paid product.</p></div><div className="fp-cycle-toggle"><button className={cycle==="monthly"?"active":""} onClick={()=>setCycle("monthly")}>Monthly</button><button className={cycle==="annual"?"active":""} onClick={()=>setCycle("annual")}>Annual</button></div></div>
      <div className="fp-plan-grid">
        <article className="fp-plan-card"><div className="fp-plan-label">COMMUNITY</div><h3>Core</h3><div className="fp-plan-price">Free</div><p className="fp-plan-desc">Open-source runtime foundation for self-hosted evaluation and development.</p><ul><li><Check size={14}/> Public Core repository</li><li><Check size={14}/> Basic runtime and telemetry</li><li><Check size={14}/> Community-led evaluation</li><li><Check size={14}/> Self-hosted deployment</li></ul><a className="fp-plan-outline" href="https://github.com/devloperdevesh/FaultPlane" target="_blank" rel="noreferrer">View Core on GitHub <ArrowRight size={14}/></a></article>
        <article className="fp-plan-card featured"><div className="fp-plan-ribbon">PRIVATE · PAID</div><div className="fp-plan-label">ADVANCED RESILIENCE</div><h3>FaultPlane Pro</h3><div className="fp-plan-price">{cycle==="monthly"?"Custom":"Custom"}<small> / {cycle}</small></div><p className="fp-plan-desc">Proprietary capabilities for teams that need deeper recovery workflows and operational controls.</p><ul>{proFeatures.slice(0,6).map(([f,href])=><li key={f}><Check size={14}/><Link href={href} className="fp-pro-feature-link">{f}</Link></li>)}</ul><Link className="fp-plan-cta" href={`/checkout?plan=${cycle}`}>Continue to payment options <ArrowRight size={15}/></Link><div className="fp-plan-foot"><LockKeyhole size={12}/> Pro source is private and not included in Core.</div></article>
        <article className="fp-plan-card"><div className="fp-plan-label">ORGANIZATIONS</div><h3>Enterprise</h3><div className="fp-plan-price">Talk to us</div><p className="fp-plan-desc">For custom deployment, security, support and integration requirements.</p><ul><li><Check size={14}/> Deployment scoping</li><li><Check size={14}/> Architecture review</li><li><Check size={14}/> Security requirements discussion</li><li><Check size={14}/> Custom commercial terms</li></ul><a className="fp-plan-outline" href="https://github.com/devloperdevesh/FaultPlane" target="_blank" rel="noreferrer">Discuss enterprise <ArrowRight size={14}/></a></article>
      </div>
      <p className="fp-pro-disclaimer">Commercial note: Pro pricing and payment processing must be configured before accepting live payments. No payment is taken from this preview until a valid provider checkout URL is configured.</p>
    </section>
  </main>;
}
