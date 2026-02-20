# go-wails-shot

Flameshot-like screenshot tool scaffold built with Go + Wails.

## Current status
- Wails desktop entrypoint implemented and buildable
- Backend API implemented (`StartCapture`, `SaveAnnotated`, `GetAppState`, `SetPreference`)
- Platform capture adapters scaffolded (`darwin`, `linux/x11`, `linux/wayland beta`)
- Deterministic export pipeline from operation log implemented in Go
- Frontend canvas editor scaffolded for annotation and operation logging
- AI-oriented docs + ADR + task templates included

## Commands
```bash
# backend smoke check
go run ./cmd/app

# unit tests
GOCACHE=/tmp/go-build GOMODCACHE=/tmp/gomodcache go test ./...

# sync static frontend files (src -> dist)
scripts/sync_frontend_dist.sh

# run desktop app in dev mode
wails dev

# build desktop app
wails build -debug
```

## Notes
- Frontend is currently plain static assets (`frontend/src`) copied to `frontend/dist`.
- Wails bindings are generated into `frontend/wailsjs` during `wails dev/build`.
- Wayland is beta by design for v1 scope.

See `docs/` for architecture, ADRs, and task contracts.
