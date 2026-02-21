package permissions

import (
	"context"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type defaultChecker struct{}

func newDefaultChecker() Checker { return &defaultChecker{} }

func (d *defaultChecker) CheckCapture(_ context.Context) (core.CapturePermissionStatus, error) {
	return core.CapturePermissionStatus{Granted: true, CanPrompt: false, Message: "capture permission check is not required on this platform"}, nil
}

func (d *defaultChecker) OpenCaptureSettings(_ context.Context) error {
	return nil
}
