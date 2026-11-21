# FaultPlane frontend update

## Included
- Refreshed responsive white/blue control-plane shell based on the supplied reference.
- Collapsible navigation with grouped workspace, resilience, observability, Pro and administration sections.
- Global command search (`Ctrl/Cmd + K`), theme toggle (persisted in local storage), profile menu, notification panel and responsive mobile navigation.
- Expanded dashboard with explicit Demo / Runtime modes. Demo figures are marked as illustrative; runtime mode reads the existing `/api/metrics` endpoint and does not silently substitute demo metrics when the API is unavailable.
- New frontend routes: `/clusters`, `/account`, `/team`, `/api-keys`, `/audit`, `/billing`, `/integrations`, `/incidents`, `/recovery-rules`, `/failure-lab`, `/login`, `/signup`, `/forgot-password`.
- Basic page-level filtering and local preview row creation for administrative screens.
- Reduced-motion-friendly, responsive CSS transitions and layout behavior.

## Important integration boundary
This ZIP is a frontend update. It does **not** create a secure identity provider, database, payment processor, entitlement server, API-key vault, audit persistence, team invitation delivery, Kubernetes agent, or recovery engine. Authentication forms deliberately explain that an identity provider is not configured. Sign-out is not represented as a successful operation. Administrative row creation is a local preview only and explicitly says it is not persisted.

Before production use, implement and test:
1. Server-side authentication, secure cookies, CSRF protections, password reset, MFA and session revocation.
2. Organization tenancy and server-enforced RBAC on every API and resource.
3. API key generation with one-time secret display, hashed storage, scopes, expiry, rotation and revocation.
4. Billing provider webhooks and server-side Core/Pro entitlement checks. Never rely on frontend route hiding for licensing.
5. Persistent audit events, invitation workflows, integrations and notification delivery.
6. Runtime endpoints for clusters, workloads, policy execution, checkpoints, logs, traces and recovery verification.
7. Safe failure-injection controls with explicit authorization, environment allowlists, blast-radius limits and kill switch.

## Core / Pro source boundary
The package metadata is marked `UNLICENSED` pending a deliberate source split. Do not publish this frontend repository as MIT while it contains proprietary Pro implementation. Keep open-source Core in its own repository/package with its own licence and publish only the Core source. Keep Pro source, build artifacts, licensing checks and server-side entitlement logic in private repositories/services. Review all dependencies and the Modern-Production-Research repository licence/ownership before commercial distribution.

## Run locally
```bash
npm ci
npm run dev
```

Set `NEXT_PUBLIC_API_URL` to the reachable FaultPlane API base URL for runtime mode. Without the backend, use Demo mode. Run `npm run lint`, `npm test` and `npm run build` after installing dependencies in a network-enabled environment.

## Pro storefront update
- Added `/pro` product and plans page with monthly/annual selection, Core vs proprietary Pro positioning, capability list and CTA.
- Added `/checkout?plan=monthly|annual` with Stripe and Razorpay hosted payment-link buttons.
- Payment buttons are intentionally shown only when matching `NEXT_PUBLIC_FAULTPLANE_*_URL` values are configured. This avoids fake payment success and prevents collecting card details in the frontend.
- Pro code remains marked `UNLICENSED`. Keep proprietary implementation and licensing enforcement in private repositories/services.
- Updated sidebar brand mark to use the supplied FaultPlane logo URL.
- Added refined card hover depth, entrance timing, status pulse and reduced-motion support.

### Payment integration still required
1. Create the real Pro products/prices and hosted payment links in Stripe and/or Razorpay.
2. Add the resulting links to `.env.local` using `.env.example` as the key list.
3. Implement server-side webhook signature verification, customer/subscription records, entitlement grants/revocations and customer portal handling before treating a payment as an activated Pro licence.
4. Keep private Pro source and release artifacts out of the public Core repository. A frontend-only lock is not a licensing boundary.

Stripe Checkout supports hosted one-time and subscription payment flows; provider-hosted checkout keeps card entry outside this frontend. See https://stripe.com/payments/checkout.
