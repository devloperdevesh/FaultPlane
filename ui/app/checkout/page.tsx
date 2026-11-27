"use client";
import Link from "next/link";
import { useEffect, useState } from "react";
import { ArrowLeft, CheckCircle2, CreditCard, LockKeyhole, ShieldCheck, Sparkles } from "lucide-react";

export default function CheckoutPage() {
  const [plan, setPlan] = useState<"monthly"|"annual">("monthly");
  useEffect(() => { const timer = window.setTimeout(() => setPlan(new URLSearchParams(window.location.search).get("plan") === "annual" ? "annual" : "monthly"), 0); return () => window.clearTimeout(timer); }, []);
  const [method, setMethod] = useState<"stripe"|"razorpay">("stripe");
  const [complete, setComplete] = useState(false);
  const [busy, setBusy] = useState(false);
  const activateDemo = () => {
    setBusy(true);
    window.setTimeout(() => {
      localStorage.setItem("fp-demo-pro-entitlement", JSON.stringify({ active: true, plan, activatedAt: new Date().toISOString(), mode: "demo" }));
      localStorage.setItem("fp-demo-notice", "FaultPlane Pro demo access activated locally");
      setBusy(false); setComplete(true);
    }, 650);
  };
  return <main className="fp-checkout-page"><div className="fp-checkout-card">
    <Link href="/pro" className="fp-checkout-back"><ArrowLeft size={14}/> Back to Pro plans</Link>
    <div className="fp-checkout-mark"><LockKeyhole size={20}/></div><div className="fp-pro-kicker">FAULTPLANE PRO · CHECKOUT PREVIEW</div>
    {complete ? <><div className="fp-checkout-success"><CheckCircle2 size={30}/><h1>Pro demo access activated</h1><p>This is a local preview entitlement only. No payment was collected and no production licence was issued.</p><Link href="/dashboard" className="fp-plan-cta">Explore dashboard <Sparkles size={15}/></Link></div></> : <>
      <h1>Choose your Pro access</h1><p className="fp-checkout-lead">Selected plan: <strong>{plan === "annual" ? "Annual" : "Monthly"} · FaultPlane Pro</strong></p>
      <div className="fp-checkout-summary"><div><span>Product</span><b>FaultPlane Pro</b></div><div><span>Billing cycle</span><b>{plan === "annual" ? "Annual" : "Monthly"}</b></div><div><span>Distribution</span><b>Private / Proprietary</b></div><div><span>Preview price</span><b>Demo · No charge</b></div></div>
      <h2>Payment provider preview</h2>
      <div className="fp-payment-methods"><button className={method==="stripe"?"selected":""} onClick={()=>setMethod("stripe")}><CreditCard size={17}/> Stripe</button><button className={method==="razorpay"?"selected":""} onClick={()=>setMethod("razorpay")}><CreditCard size={17}/> Razorpay</button></div>
      <button className="fp-pay-button stripe" disabled={busy} onClick={activateDemo}>{busy?"Preparing demo access…":`Preview ${method === "stripe" ? "Stripe" : "Razorpay"} checkout`}</button>
      <div className="fp-payment-pending"><ShieldCheck size={18}/><div><b>Demo payment mode</b><p>This button simulates the customer journey and activates a local demo entitlement. Connect a real provider and server-side webhook before accepting money.</p></div></div>
    </>}
    <div className="fp-checkout-security"><ShieldCheck size={15}/><span>Production Pro access must be granted only after a verified payment webhook and server-side licence check.</span></div>
  </div></main>;
}
