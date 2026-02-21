package darwin

import (
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type Adapter struct{}

func NewAdapter() *Adapter { return &Adapter{} }

func (a *Adapter) SessionType() string { return "darwin" }
func (a *Adapter) WaylandBeta() bool   { return false }

func (a *Adapter) Capture(ctx context.Context, mode string, outputPath string) (core.CaptureResult, error) {
	log.Printf("[capture][darwin] Start mode=%s output=%s", mode, outputPath)
	args := []string{"-x"}
	if mode == "region" {
		args = append(args, "-i")
	}
	args = append(args, outputPath)
	cmd := exec.CommandContext(ctx, "screencapture", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(out))
		log.Printf("[capture][darwin] screencapture failed err=%v out=%q", err, msg)
		lower := strings.ToLower(msg)
		if strings.Contains(lower, "screen recording") ||
			strings.Contains(lower, "permission") ||
			strings.Contains(lower, "not authorized") ||
			strings.Contains(lower, "could not create image from display") ||
			strings.Contains(lower, "could not create image from rect") {
			return core.CaptureResult{}, &core.AppError{
				Code:    core.ErrCapturePermission,
				Message: "Screen Recording permission is missing or blocked. Allow it in System Settings > Privacy & Security > Screen Recording, then restart the app.",
			}
		}
		return core.CaptureResult{}, &core.AppError{Code: core.ErrCaptureUnavailable, Message: "screencapture failed: " + msg}
	}
	f, err := os.Open(outputPath)
	if err != nil {
		log.Printf("[capture][darwin] failed to open output=%s err=%v", outputPath, err)
		return core.CaptureResult{}, &core.AppError{Code: core.ErrReadFailed, Message: err.Error()}
	}
	defer f.Close()
	img, _, err := image.DecodeConfig(f)
	if err != nil {
		log.Printf("[capture][darwin] failed to decode output=%s err=%v", outputPath, err)
		return core.CaptureResult{}, &core.AppError{Code: core.ErrDecodeFailed, Message: err.Error()}
	}
	log.Printf("[capture][darwin] success mode=%s size=%dx%d", mode, img.Width, img.Height)
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
