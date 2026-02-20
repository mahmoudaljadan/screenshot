# Spec: MVP Capture + Annotate + Save

## Functional scope
- Capture: fullscreen and region mode request path (platform-dependent implementation)
- Tools: rectangle, line, arrow, text, blur, pixelate
- Editing: undo/redo
- Export: PNG/JPEG

## Non-goals
- Cloud upload/share
- Global hotkey
- Advanced blur quality tuning

## API
- `StartCapture(mode string) (CaptureResult, error)`
- `SaveAnnotated(req ExportRequest) (ExportResult, error)`
- `GetAppState() (AppState, error)`
- `SetPreference(key string, value string) error`

## Error model
Use `AppError` codes:
- `ERR_INVALID_OP_KIND`
- `ERR_INVALID_OP_PAYLOAD`
- `ERR_CAPTURE_UNAVAILABLE`
- `ERR_WAYLAND_PREREQUISITE`
- `ERR_ENCODE_FAILED`
- `ERR_DECODE_FAILED`
- `ERR_READ_FAILED`
- `ERR_WRITE_FAILED`
