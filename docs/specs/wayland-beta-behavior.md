# Spec: Wayland Beta Behavior

## Positioning
Wayland support is beta in v1 and clearly labeled in state/UI.

## Preconditions
- `grim` available for screenshot capture
- `xdg-desktop-portal` available
- `slurp` required for region capture path

## Failure behavior
- Missing prerequisites must return `ERR_WAYLAND_PREREQUISITE` with actionable text.
- Capture execution failures return `ERR_CAPTURE_UNAVAILABLE` with command stderr.

## UX guidance
- Surface prerequisites in settings/help and capture errors.
- Keep X11 as first-class Linux path for initial stable releases.
