package darwin

import (
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type Adapter struct{}

func NewAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) SessionType() string { return "darwin" }
func (a *Adapter) WaylandBeta() bool   { return false }

func (a *Adapter) Capture(ctx context.Context, _ string, outputPath string) (core.CaptureResult, error) {
	cmd := exec.CommandContext(ctx, "screencapture", "-x", outputPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return core.CaptureResult{}, &core.AppError{Code: core.ErrCaptureUnavailable, Message: "screencapture failed: " + string(out)}
	}
	f, err := os.Open(outputPath)
	if err != nil {
		return core.CaptureResult{}, &core.AppError{Code: core.ErrReadFailed, Message: err.Error()}
	}
	defer f.Close()
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		return core.CaptureResult{}, &core.AppError{Code: core.ErrDecodeFailed, Message: err.Error()}
	}
	return core.CaptureResult{
		ImagePath: outputPath,
		Width:     img.Width,
		Height:    img.Height,
		Displays: []core.DisplayInfo{{
			ID:     "primary",
			X:      0,
			Y:      0,
			Width:  img.Width,
			Height: img.Height,
			Scale:  1,
		}},
		SessionID: "capture-darwin",
	}, nil
}
