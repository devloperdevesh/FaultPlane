# FaultPlane Beta Program

## Overview

FaultPlane is currently seeking engineers interested in evaluating
network-failure recovery for long-running workloads.

The beta focuses on reproducible testing and technical feedback.

## Test Environment

Before testing, record:

- Operating system
- Linux kernel
- CPU
- Go version
- FaultPlane commit
- Deployment method

Do not share:

- API keys
- Passwords
- Access tokens
- Private keys
- Confidential production data

## Beta Test Flow

1. Install FaultPlane in a controlled environment.
2. Run the baseline workload.
3. Introduce a controlled connection failure.
4. Observe FaultPlane.
5. Verify fallback behavior.
6. Verify workload continuity.
7. Record measured results.
8. Report the result.

## Feedback

Please report:

- Installation experience
- Failure detection result
- Fallback behavior
- Workload continuity
- Observed latency
- Any data-loss symptoms
- Errors or unexpected behavior
- Improvements you would like to see

## Publishing Results

Only publish results that were actually measured.

Do not publish customer-sensitive information.

If a result is quoted publicly, obtain explicit permission from
the person who provided it.

## Beta Status

FaultPlane is under active development.

Beta feedback is used to improve reliability, documentation,
testing, and production readiness.
