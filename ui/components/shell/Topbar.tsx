"use client";

import { useEffect, useMemo, useState } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Bell, ChevronDown, Command, CreditCard, HelpCircle, KeyRound, LogOut, Menu, Moon, Search, Settings, ShieldCheck, Sun, UserRound, Users, X } from "lucide-react";

const destinations = [
  ["Dashboard", "/dashboard"], ["Clusters", "/clusters"], ["Workers", "/workers"], ["Workloads", "/workflows"],
  ["Recovery", "/recovery"], ["Recovery rules", "/recovery-rules"], ["Checkpoints", "/checkpoints"],
  ["Incidents", "/incidents"], ["Failure Lab", "/failure-lab"], ["Metrics", "/metrics"], ["Telemetry", "/telemetry"],
  ["Team members", "/team"], ["API keys", "/api-keys"], ["Audit log", "/audit"], ["Billing", "/billing"], ["FaultPlane Pro", "/pro"], ["Settings", "/settings"],
] as const;

export default function Topbar() {
  const pathname = usePathname();
  const router = useRouter();
  const [query, setQuery] = useState("");
  const [searchOpen, setSearchOpen] = useState(false);
  const [profileOpen, setProfileOpen] = useState(false);
  const [noticeOpen, setNoticeOpen] = useState(false);
  const [dark, setDark] = useState(false);
  const [mobileNav, setMobileNav] = useState(false);
  const [range, setRange] = useState("Last 24 hours");
  const [notifications, setNotifications] = useState<string[]>(["Demo workspace connected", "Recovery policy updated", "Weekly runtime report is ready"]);
  const results = useMemo(() => destinations.filter(([label]) => label.toLowerCase().includes(query.toLowerCase())).slice(0, 7), [query]);

  useEffect(() => {
    const saved = localStorage.getItem("fp-theme");
    const isDark = saved === "dark";
    const themeTimer = window.setTimeout(() => setDark(isDark), 0);
    document.documentElement.dataset.theme = isDark ? "dark" : "light";
    const onKey = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") { event.preventDefault(); setSearchOpen(true); }
      if (event.key === "Escape") { setSearchOpen(false); setProfileOpen(false); setNoticeOpen(false); }
    };
    window.addEventListener("keydown", onKey);
    return () => { window.clearTimeout(themeTimer); window.removeEventListener("keydown", onKey); };
  }, []);

  const toggleTheme = () => {
    const next = !dark; setDark(next); document.documentElement.dataset.theme = next ? "dark" : "light"; localStorage.setItem("fp-theme", next ? "dark" : "light");
  };
  const go = (href: string) => { router.push(href); setSearchOpen(false); setProfileOpen(false); setMobileNav(false); setQuery(""); };

  return <header className="fp-topbar">
    <button className="fp-icon-button fp-mobile-menu" aria-label="Open navigation" onClick={() => setMobileNav(v => !v)}><Menu size={18}/></button>
    <div className="fp-search-wrap">
      <button className="fp-search-trigger" onClick={() => setSearchOpen(true)}><Search size={17}/><span>Search workloads, clusters, events, rules...</span><kbd><Command size={11}/> K</kbd></button>
    </div>
    <div className="fp-topbar-actions">
      <button className="fp-select-chip" onClick={() => go("/clusters")}><span className="fp-dot"/>All Clusters <ChevronDown size={14}/></button>
      <select className="fp-range-select" aria-label="Time range" value={range} onChange={e => setRange(e.target.value)}><option>Last 1 hour</option><option>Last 6 hours</option><option>Last 24 hours</option><option>Last 7 days</option><option>Last 30 days</option></select>
      <button className="fp-icon-button" aria-label="Toggle theme" onClick={toggleTheme}>{dark ? <Sun size={18}/> : <Moon size={18}/>}</button>
      <div className="fp-popover-anchor"><button className="fp-icon-button fp-bell" aria-label="Notifications" onClick={() => setNoticeOpen(v => !v)}><Bell size={18}/><i/></button>{noticeOpen && <div className="fp-popover fp-notice"><strong>Notifications</strong><div className="fp-notification-list">{notifications.length ? notifications.map((n,i)=><p className="fp-notification-item" key={`${n}-${i}`}>{n}<small>Just now · Demo event</small></p>) : <p className="fp-muted">You’re all caught up.</p>}</div><button onClick={() => {setNotifications([]);}}>Mark all as read</button><button onClick={() => go("/incidents")}>Open incident center →</button></div>}</div>
      <div className="fp-popover-anchor"><button className="fp-avatar" aria-label="Open profile menu" onClick={() => setProfileOpen(v => !v)}>D</button>{profileOpen && <div className="fp-popover fp-profile"><div className="fp-profile-head"><span className="fp-avatar large">D</span><div><strong>Devesh Chauhan</strong><small>Workspace owner</small></div></div><button onClick={() => go("/account")}><UserRound size={15}/> My profile</button><button onClick={() => go("/team")}><Users size={15}/> Team & members</button><button onClick={() => go("/api-keys")}><KeyRound size={15}/> API keys</button><button onClick={() => go("/billing")}><CreditCard size={15}/> Billing & plan</button><button onClick={() => go("/settings")}><Settings size={15}/> Settings</button><div className="fp-popover-divider"/><button onClick={() => { setProfileOpen(false); localStorage.removeItem("fp-demo-session"); go("/login"); }}><LogOut size={15}/> Sign out</button></div>}</div>
    </div>
    {mobileNav && <div className="fp-mobile-nav"><button className="fp-mobile-close" onClick={() => setMobileNav(false)}><X size={18}/></button>{destinations.map(([label, href]) => <button key={href} className={pathname === href ? "active" : ""} onClick={() => go(href)}>{label}</button>)}</div>}
    {searchOpen && <div className="fp-modal-backdrop" onMouseDown={e => { if (e.target === e.currentTarget) setSearchOpen(false); }}><div className="fp-command-modal"><div className="fp-command-input"><Search size={18}/><input autoFocus placeholder="Search pages and actions..." value={query} onChange={e => setQuery(e.target.value)}/><button onClick={() => setSearchOpen(false)}>ESC</button></div><div className="fp-command-results">{results.map(([label, href]) => <button key={href} onClick={() => go(href)}><span>{label}</span><span>Open <kbd>↵</kbd></span></button>)}{results.length === 0 && <p>No matching destinations.</p>}</div><footer><span><HelpCircle size={13}/> Navigate anywhere in FaultPlane</span><span>↑↓ select · ↵ open · esc close</span></footer></div></div>}
  </header>;
}
