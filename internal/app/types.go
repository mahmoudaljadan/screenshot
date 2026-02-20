package app

import "github.com/mohamoundaljadan/screenshot/internal/core"

type DisplayInfo = core.DisplayInfo
type CaptureResult = core.CaptureResult
type AnnotationOp = core.AnnotationOp
type ExportRequest = core.ExportRequest
type ExportResult = core.ExportResult
type AppState = core.AppState
type AppError = core.AppError

const (
	ErrInvalidOpKind       = core.ErrInvalidOpKind
	ErrInvalidOpPayload    = core.ErrInvalidOpPayload
	ErrCaptureUnavailable  = core.ErrCaptureUnavailable
	ErrWaylandPrerequisite = core.ErrWaylandPrerequisite
	ErrEncodeFailed        = core.ErrEncodeFailed
	ErrDecodeFailed        = core.ErrDecodeFailed
	ErrWriteFailed         = core.ErrWriteFailed
	ErrReadFailed          = core.ErrReadFailed
)
