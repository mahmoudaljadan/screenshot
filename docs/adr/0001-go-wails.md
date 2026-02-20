# ADR 0001: Go + Wails stack

## Status
Accepted

## Context
Need a desktop screenshot app with strong Go backend ergonomics and maintainable cross-platform UI.

## Decision
Use Go for backend domain/platform logic and Wails for desktop shell + webview frontend.

## Consequences
- Faster backend iteration for Go-centric development.
- Platform-specific capture code still required.
- Webview canvas is suitable for annotation UX but final export stays Go-side for deterministic output.
