# Project Brief

## Goal
Build a cross-platform screenshot and annotation desktop app with Go backend and Wails frontend.

## MVP
- Capture screen image
- Annotate on canvas (rect, line, arrow, text, blur, pixelate)
- Save deterministic output (Go-side render)

## Constraints
- Linux X11 first-class support
- Wayland beta support with explicit preflight diagnostics
- No global hotkeys in v1
- Dev builds first; packaging/signing postponed

## Defaults
- Default output format: PNG
- Telemetry: disabled by default
