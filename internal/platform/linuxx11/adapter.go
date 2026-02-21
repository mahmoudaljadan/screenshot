package linuxx11

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type Adapter struct{}

func NewAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) SessionType() string { return "x11" }
func (a *Adapter) WaylandBeta() bool   { return false }

func (a *Adapter) Capture(ctx context.Context, mode string, outputPath string) (core.CaptureResult, error) {
	if err := runCaptureTool(ctx, mode, outputPath); err != nil {
		return core.CaptureResult{}, err
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
		SessionID: "capture-x11",
	}, nil
}

func runCaptureTool(ctx context.Context, mode, outputPath string) error {
	if _, err := exec.LookPath("maim"); err == nil {
		args := []string{}
		if mode == "region" {
			args = append(args, "-s")
		}
		args = append(args, outputPath)
		if out, runErr := exec.CommandContext(ctx, "maim", args...).CombinedOutput(); runErr != nil {
			return &core.AppError{Code: core.ErrCaptureUnavailable, Message: "maim failed: " + string(out)}
		}
		return nil
	}
	if _, err := exec.LookPath("scrot"); err == nil {
		args := []string{}
		if mode == "region" {
			args = append(args, "-s")
		}
		args = append(args, outputPath)
		if out, runErr := exec.CommandContext(ctx, "scrot", args...).CombinedOutput(); runErr != nil {
			return &core.AppError{Code: core.ErrCaptureUnavailable, Message: "scrot failed: " + string(out)}
		}
		return nil
	}
	if _, err := exec.LookPath("import"); err == nil {
		if out, runErr := exec.CommandContext(ctx, "import", "-window", "root", outputPath).CombinedOutput(); runErr != nil {
			return &core.AppError{Code: core.ErrCaptureUnavailable, Message: "import failed: " + string(out)}
		}
		return nil
	}
	return &core.AppError{
		Code: core.ErrCaptureUnavailable,
		Message: fmt.Sprintf("no X11 screenshot tool found. Install one of: %s",
			"maim, scrot, or imagemagick(import)"),
	}
}
