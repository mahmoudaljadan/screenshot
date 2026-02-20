# ADR 0002: Deterministic operation-log rendering

## Status
Accepted

## Context
Frontend annotation may differ by browser/webview behavior. Export output must remain deterministic and testable.

## Decision
Frontend emits an operation log only. Backend replays operations in Go on export.

## Consequences
- Deterministic output across machines.
- Strong unit/integration testability.
- Requires keeping op schema and renderer synchronized.
