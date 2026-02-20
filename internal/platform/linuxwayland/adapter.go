package linuxwayland

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

func (a *Adapter) SessionType() string { return "wayland" }
func (a *Adapter) WaylandBeta() bool   { return true }

func (a *Adapter) Capture(ctx context.Context, mode, outputPath string) (core.CaptureResult, error) {
	if err := preflight(); err != nil {
		return core.CaptureResult{}, err
	}

	args := []string{outputPath}
	if mode == "region" {
		if _, err := exec.LookPath("slurp"); err != nil {
			return core.CaptureResult{}, &core.AppError{Code: core.ErrWaylandPrerequisite, Message: "region capture requires slurp on Wayland"}
		}
		cmd := exec.CommandContext(ctx, "sh", "-c", "grim -g \"$(slurp)\" "+outputPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return core.CaptureResult{}, &core.AppError{Code: core.ErrCaptureUnavailable, Message: "grim/slurp failed: " + string(out)}
		}
	} else {
		cmd := exec.CommandContext(ctx, "grim", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			return core.CaptureResult{}, &core.AppError{Code: core.ErrCaptureUnavailable, Message: "grim failed: " + string(out)}
		}
	}

	f, err := os.Open(outputPath)
	if err != nil {
		return core.CaptureResult{}, &core.AppError{Code: core.ErrReadFailed, Message: err.Error()}
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return core.CaptureResult{}, &core.AppError{Code: core.ErrDecodeFailed, Message: err.Error()}
	}
	return core.CaptureResult{
		ImagePath: outputPath,
		Width:     cfg.Width,
		Height:    cfg.Height,
		Displays: []core.DisplayInfo{{
			ID:     "primary",
			X:      0,
			Y:      0,
			Width:  cfg.Width,
			Height: cfg.Height,
			Scale:  1,
		}},
		SessionID: "capture-wayland-beta",
	}, nil
}

func preflight() error {
	if _, err := exec.LookPath("grim"); err != nil {
		return &core.AppError{Code: core.ErrWaylandPrerequisite, Message: "Wayland beta needs grim installed"}
	}
	if _, err := exec.LookPath("xdg-desktop-portal"); err != nil {
		return &core.AppError{Code: core.ErrWaylandPrerequisite, Message: fmt.Sprintf("Wayland beta requires xdg-desktop-portal. install it for your distro and retry")}
	}
	return nil
}
