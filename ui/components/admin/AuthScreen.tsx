"use client";

import { useState, type FormEvent } from "react";
import Image from "next/image";
import Link from "next/link";
import { ArrowRight, LockKeyhole, Mail } from "lucide-react";
import { FaGithub } from "react-icons/fa";
import { FcGoogle } from "react-icons/fc";

const copy = {
  login: {
    title: "Welcome back",
    description: "Sign in to your FaultPlane workspace.",
    action: "Sign in",
  },
  signup: {
    title: "Create your workspace",
    description: "Set up an account to manage your AI workloads.",
    action: "Create account",
  },
  reset: {
    title: "Reset your password",
    description: "Enter the work email associated with your account.",
    action: "Send reset link",
  },
};

export default function AuthScreen({ mode }: { mode: keyof typeof copy }) {
  const [passwordVisible, setPasswordVisible] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const isLogin = mode === "login";
  const { title, description, action } = copy[mode];

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitted(true);
  }

  return (
    <main className="fp-auth-page">
      <div className="fp-auth-layout">
        <header className="fp-auth-brand">
          <Image src="/faultplane-logo.png" alt="" width={40} height={40} loading="eager" unoptimized />
          <div>
            <strong>FaultPlane</strong>
            <span>AI Runtime Control Plane</span>
          </div>
        </header>

        <section className="fp-auth-card" aria-labelledby="auth-title">
          <h1 id="auth-title">{title}</h1>
          <p className="fp-auth-description">{description}</p>

          {isLogin && (
            <>
              <div className="fp-auth-providers">
                <button type="button" aria-describedby="auth-provider-note" disabled>
                  <FaGithub aria-hidden="true" /> Continue with GitHub
                </button>
                <button type="button" aria-describedby="auth-provider-note" disabled>
                  <FcGoogle aria-hidden="true" /> Continue with Google
                </button>
              </div>
              <p id="auth-provider-note" className="fp-auth-note">
                GitHub and Google sign-in are unavailable until OAuth is configured.
              </p>
              <div className="fp-auth-divider">or continue with email</div>
            </>
          )}

          <form className="fp-auth-form" onSubmit={handleSubmit}>
            <div className="fp-auth-field">
              <label htmlFor="auth-email">Work email</label>
              <div className="fp-auth-input">
                <Mail size={17} aria-hidden="true" />
                <input
                  id="auth-email"
                  type="email"
                  autoComplete="email"
                  placeholder="you@company.com"
                  required
                />
              </div>
            </div>

            {mode !== "reset" && (
              <div className="fp-auth-field">
                <div className="fp-auth-field-heading">
                  <label htmlFor="auth-password">Password</label>
                  {isLogin && <Link href="/forgot-password">Forgot password?</Link>}
                </div>
                <div className="fp-auth-input">
                  <LockKeyhole size={17} aria-hidden="true" />
                  <input
                    id="auth-password"
                    type={passwordVisible ? "text" : "password"}
                    autoComplete={isLogin ? "current-password" : "new-password"}
                    placeholder="At least 8 characters"
                    minLength={8}
                    required
                  />
                  <button
                    type="button"
                    className="fp-auth-password-toggle"
                    aria-controls="auth-password"
                    onClick={() => setPasswordVisible((visible) => !visible)}
                  >
                    {passwordVisible ? "Hide password" : "Show password"}
                  </button>
                </div>
              </div>
            )}

            <p id="auth-preview-note" className="fp-auth-note">
              This preview requires an identity provider to authenticate. It does not send your credentials.
            </p>
            <button type="submit" className="fp-auth-submit" aria-describedby="auth-preview-note">
              {action} <ArrowRight size={17} aria-hidden="true" />
            </button>
            <div role="status" className={submitted ? "fp-auth-message" : undefined}>
              {submitted && "Authentication is unavailable. No credentials were sent."}
            </div>
          </form>

          <p className="fp-auth-bottom">
            {isLogin ? (
              <>New to FaultPlane? <Link href="/signup">Create an account</Link></>
            ) : (
              <>Back to <Link href="/login">sign in</Link></>
            )}
          </p>
        </section>
      </div>
    </main>
  );
}
