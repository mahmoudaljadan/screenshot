# System Map

## Components
- Frontend (`frontend/src`): canvas UX, operation log creation, undo/redo, save trigger
- App service (`internal/app`): Wails-bound API surface and preferences
- Capture manager (`internal/capture`): adapter selection and capture orchestration
- Platform adapters (`internal/platform/*`): platform-specific capture logic and preflight checks
- Annotation engine (`internal/annotate`): op validation/sorting and rendering
- Export service (`internal/export`): decode base image, replay ops, encode output

## Interface contracts
- Frontend sends `ExportRequest` containing `ops`
- Backend validates op kinds/payloads and fails with explicit error codes
- Backend writes final image and returns `ExportResult`
