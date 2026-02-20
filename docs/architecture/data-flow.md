# Data Flow

1. `StartCapture(mode)` called from frontend.
2. Capture manager selects adapter by OS/session.
3. Adapter captures base image to temp path and returns `CaptureResult`.
4. Frontend loads base image and builds operation log from user edits.
5. `SaveAnnotated(req)` called with base image path + ops.
6. Export service validates ops, sorts deterministically, applies ops in Go renderer.
7. Export service writes PNG/JPEG and returns output metadata.
