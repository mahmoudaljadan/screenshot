package core

import "encoding/json"

type DisplayInfo struct {
	ID     string `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Scale  int    `json:"scale"`
}

type CaptureResult struct {
	ImagePath string        `json:"imagePath"`
	Width     int           `json:"width"`
	Height    int           `json:"height"`
	Displays  []DisplayInfo `json:"displays"`
	SessionID string        `json:"sessionId"`
}

type AnnotationOp struct {
	ID      string          `json:"id"`
	Kind    string          `json:"kind"`
	Z       int             `json:"z"`
	Payload json.RawMessage `json:"payload"`
}

type ExportRequest struct {
	BaseImagePath string         `json:"baseImagePath"`
	Ops           []AnnotationOp `json:"ops"`
	Format        string         `json:"format"`
	Quality       int            `json:"quality"`
	OutputPath    string         `json:"outputPath"`
}

type ExportResult struct {
	OutputPath string `json:"outputPath"`
	Bytes      int64  `json:"bytes"`
	Format     string `json:"format"`
}

type AppState struct {
	Platform    string `json:"platform"`
	SessionType string `json:"sessionType"`
	WaylandBeta bool   `json:"waylandBeta"`
	Version     string `json:"version"`
}

type CapturePermissionStatus struct {
	Granted      bool   `json:"granted"`
	CanPrompt    bool   `json:"canPrompt"`
	Message      string `json:"message"`
	SettingsHint string `json:"settingsHint"`
}

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Code + ": " + e.Message
}

const (
	ErrInvalidOpKind       = "ERR_INVALID_OP_KIND"
	ErrInvalidOpPayload    = "ERR_INVALID_OP_PAYLOAD"
	ErrCaptureUnavailable  = "ERR_CAPTURE_UNAVAILABLE"
	ErrCapturePermission   = "ERR_CAPTURE_PERMISSION"
	ErrWaylandPrerequisite = "ERR_WAYLAND_PREREQUISITE"
	ErrEncodeFailed        = "ERR_ENCODE_FAILED"
	ErrDecodeFailed        = "ERR_DECODE_FAILED"
	ErrWriteFailed         = "ERR_WRITE_FAILED"
	ErrReadFailed          = "ERR_READ_FAILED"
)
