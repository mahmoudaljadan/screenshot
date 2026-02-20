# TASK-0001 Foundation

## Goal
Set up project skeleton, backend API, and AI documentation baseline.

## Non-goals
- Signed builds
- Full Wails packaging pipeline

## Files to touch
- `go.mod`
- `cmd/app/main.go`
- `internal/*`
- `frontend/src/*`
- `docs/*`

## API/type changes
- Initial `Service` APIs and shared DTOs in `internal/app/types.go`

## Acceptance criteria
- `go test ./...` passes
- `go run ./cmd/app` prints app state
- Docs and ADR files exist

## Test checklist
- Unit tests for op validation
- Unit tests for exporter deterministic behavior

## Rollback notes
If this task regresses startup, rollback `cmd/app/main.go` and keep docs-only changes.
