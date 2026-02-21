package permissions

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type darwinChecker struct{}

func newDarwinChecker() Checker { return &darwinChecker{} }

func (d *darwinChecker) CheckCapture(ctx context.Context) (core.CapturePermissionStatus, error) {
	tmp := filepath.Join(os.TempDir(), "permission-check-"+time.Now().Format("20060102T150405.000")+".png")
	defer os.Remove(tmp)

	cmd := exec.CommandContext(ctx, "screencapture", "-x", tmp)
	out, err := cmd.CombinedOutput()
	if err == nil {
		log.Printf("[perm][darwin] screencapture permission check: granted")
		return core.CapturePermissionStatus{
			Granted:      true,
			CanPrompt:    true,
			Message:      "Screen Recording permission granted",
			SettingsHint: "System Settings > Privacy & Security > Screen Recording",
		}, nil
	}

	msg := strings.TrimSpace(string(out))
	lower := strings.ToLower(msg)
	if strings.Contains(lower, "screen recording") ||
		strings.Contains(lower, "permission") ||
		strings.Contains(lower, "not authorized") ||
		strings.Contains(lower, "could not create image from display") ||
		strings.Contains(lower, "could not create image from rect") {
		log.Printf("[perm][darwin] screencapture permission check: denied out=%q", msg)
		return core.CapturePermissionStatus{
			Granted:      false,
			CanPrompt:    true,
			Message:      "Screen Recording permission is required",
			SettingsHint: "System Settings > Privacy & Security > Screen Recording",
		}, nil
	}

	// Non-permission failure: do not hard-block capture; return unknown-as-granted with message.
	log.Printf("[perm][darwin] screencapture permission check unknown failure err=%v out=%q", err, msg)
	return core.CapturePermissionStatus{
		Granted:      true,
		CanPrompt:    true,
		Message:      "Permission check could not be confirmed; capture will be attempted",
		SettingsHint: "System Settings > Privacy & Security > Screen Recording",
	}, nil
}

func (d *darwinChecker) OpenCaptureSettings(_ context.Context) error {
	log.Printf("[perm][darwin] opening Screen Recording settings")
	cmd := exec.Command("open", "x-apple.systempreferences:com.apple.preference.security?Privacy_ScreenCapture")
	return cmd.Run()
}
