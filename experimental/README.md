# FaultPlane Experimental Features

This directory contains optional Phase 7 capabilities.

Experimental features are isolated from the production recovery core.

Current feature areas:

- finops
- speculative
- time_travel
- wasm
- gitops

Rules:

1. Experimental code must not silently alter core recovery behaviour.
2. Experimental APIs must have explicit boundaries.
3. Runtime values must come from real system state or explicit external integrations.
4. No fabricated production metrics.
5. Experimental functionality must be independently testable.
6. Features remain opt-in until promoted to the core product.
