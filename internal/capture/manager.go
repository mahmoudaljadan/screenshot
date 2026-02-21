package capture

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mohamoundaljadan/screenshot/internal/core"
	"github.com/mohamoundaljadan/screenshot/internal/platform/darwin"
	"github.com/mohamoundaljadan/screenshot/internal/platform/linuxwayland"
	"github.com/mohamoundaljadan/screenshot/internal/platform/linuxx11"
	"github.com/mohamoundaljadan/screenshot/internal/platform/permissions"
)

type Adapter interface {
	SessionType() string
	WaylandBeta() bool
	Capture(ctx context.Context, mode, outputPath string) (core.CaptureResult, error)
}

type Manager struct {
	adapter Adapter
	perm    permissions.Checker
	tmpDir  string
}

func NewManager(tmpDir string) *Manager {
	if tmpDir == "" {
		tmpDir = os.TempDir()
	}
	return &Manager{
		adapter: chooseAdapter(),
		perm:    permissions.NewChecker(),
		tmpDir:  tmpDir,
	}
}

func (m *Manager) SessionInfo() (string, bool) {
	return m.adapter.SessionType(), m.adapter.WaylandBeta()
}

func (m *Manager) Capture(ctx context.Context, mode string) (core.CaptureResult, error) {
	timestamp := time.Now().UTC().Format("20060102T150405.000Z")
	outputPath := filepath.Join(m.tmpDir, fmt.Sprintf("capture-%s.png", timestamp))
	return m.adapter.Capture(ctx, mode, outputPath)
}

func (m *Manager) CheckCapturePermission(ctx context.Context) (core.CapturePermissionStatus, error) {
	return m.perm.CheckCapture(ctx)
}

func (m *Manager) OpenCapturePermissionSettings(ctx context.Context) error {
	return m.perm.OpenCaptureSettings(ctx)
}

func chooseAdapter() Adapter {
	switch runtime.GOOS {
	case "darwin":
		return darwin.NewAdapter()
	case "linux":
		if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
			return linuxwayland.NewAdapter()
		}
		return linuxx11.NewAdapter()
	default:
		return unsupportedAdapter{sessionType: runtime.GOOS}
	}
}

type unsupportedAdapter struct {
	sessionType string
}

func (u unsupportedAdapter) SessionType() string { return u.sessionType }
func (u unsupportedAdapter) WaylandBeta() bool   { return false }

func (u unsupportedAdapter) Capture(_ context.Context, _, _ string) (core.CaptureResult, error) {
	return core.CaptureResult{}, &core.AppError{Code: core.ErrCaptureUnavailable, Message: "screen capture is not supported on this platform yet"}
}
