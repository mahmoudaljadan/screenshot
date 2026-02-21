package permissions

import (
	"context"
	"runtime"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type Checker interface {
	CheckCapture(ctx context.Context) (core.CapturePermissionStatus, error)
	OpenCaptureSettings(ctx context.Context) error
}

func NewChecker() Checker {
	switch runtime.GOOS {
	case "darwin":
		return newDarwinChecker()
	default:
		return newDefaultChecker()
	}
}
