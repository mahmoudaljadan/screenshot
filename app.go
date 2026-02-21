package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	stdruntime "runtime"
	"time"

	appsvc "github.com/mohamoundaljadan/screenshot/internal/app"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx context.Context
	svc *appsvc.Service
}

func NewApp(version string) (*App, error) {
	svc, err := appsvc.NewDefaultService(version)
	if err != nil {
		return nil, err
	}
	return &App{svc: svc}, nil
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	runtime.WindowSetSize(ctx, 560, 140)
	runtime.WindowCenter(ctx)
}

func (a *App) StartCapture(mode string) (appsvc.CaptureResult, error) {
	log.Printf("[app] StartCapture requested mode=%s", mode)
	if a.ctx != nil {
		log.Printf("[app] hiding window before capture")
		runtime.Hide(a.ctx)
		// Allow the OS compositor to remove our app window before capture.
		time.Sleep(420 * time.Millisecond)
	}
	result, err := a.svc.StartCapture(mode)
	if a.ctx != nil {
		log.Printf("[app] showing window after capture")
		runtime.Show(a.ctx)
		runtime.WindowUnminimise(a.ctx)
	}
	if err != nil {
		log.Printf("[app] StartCapture failed mode=%s err=%v", mode, err)
	} else {
		log.Printf("[app] StartCapture success mode=%s image=%s", mode, result.ImagePath)
	}
	return result, err
}

func (a *App) SaveAnnotated(req appsvc.ExportRequest) (appsvc.ExportResult, error) {
	return a.svc.SaveAnnotated(req)
}

func (a *App) GetAppState() (appsvc.AppState, error) {
	return a.svc.GetAppState()
}

func (a *App) SetPreference(key string, value string) error {
	return a.svc.SetPreference(key, value)
}

func (a *App) EnterEditorMode() {
	if a.ctx == nil {
		return
	}
	log.Printf("[app] EnterEditorMode fullscreen")
	runtime.WindowFullscreen(a.ctx)
}

func (a *App) EnterScreenEditorMode() {
	a.EnterEditorMode()
}

func (a *App) EnterRegionEditorMode(imageWidth int, imageHeight int) {
	if a.ctx == nil {
		return
	}
	log.Printf("[app] EnterRegionEditorMode image=%dx%d", imageWidth, imageHeight)
	runtime.WindowUnfullscreen(a.ctx)

	// Keep most of the desktop visible: show a compact floating editor window.
	const maxW = 1280
	const maxH = 900
	const minW = 520
	const minH = 360
	w := imageWidth + 32
	h := imageHeight + 88
	if w < minW {
		w = minW
	}
	if h < minH {
		h = minH
	}
	if w > maxW {
		w = maxW
	}
	if h > maxH {
		h = maxH
	}
	runtime.WindowSetSize(a.ctx, w, h)
	runtime.WindowCenter(a.ctx)
	runtime.Show(a.ctx)
	runtime.WindowUnminimise(a.ctx)
}

func (a *App) ExitEditorMode() {
	if a.ctx == nil {
		return
	}
	log.Printf("[app] ExitEditorMode launcher-size")
	runtime.WindowUnfullscreen(a.ctx)
	runtime.WindowSetSize(a.ctx, 560, 140)
	runtime.WindowCenter(a.ctx)
}

func (a *App) CheckCapturePermission() (appsvc.CapturePermissionStatus, error) {
	status, err := a.svc.CheckCapturePermission()
	if err != nil {
		log.Printf("[app] CheckCapturePermission failed err=%v", err)
	} else {
		log.Printf("[app] CheckCapturePermission granted=%v message=%q", status.Granted, status.Message)
	}
	return status, err
}

func (a *App) OpenCapturePermissionSettings() error {
	log.Printf("[app] OpenCapturePermissionSettings requested")
	return a.svc.OpenCapturePermissionSettings()
}

func (a *App) LoadCaptureImage(path string) (string, error) {
	log.Printf("[app] LoadCaptureImage path=%s", path)
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "image/png"
	}
	encoded := base64.StdEncoding.EncodeToString(b)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
}

func (a *App) PromptSavePath(format string) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("app context is not ready")
	}

	ext := ".png"
	display := "PNG Image"
	pattern := "*.png"
	if format == "jpg" || format == "jpeg" {
		ext = ".jpg"
		display = "JPEG Image"
		pattern = "*.jpg;*.jpeg"
	}

	now := time.Now().Format("20060102-150405")
	defaultFilename := "capture-" + now + ext

	defaultDir := ""
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		switch stdruntime.GOOS {
		case "darwin", "linux", "windows":
			defaultDir = filepath.Join(home, "Pictures", "go-wails-shot")
		default:
			defaultDir = home
		}
	}

	log.Printf("[app] PromptSavePath format=%s defaultDir=%s", format, defaultDir)
	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Save Screenshot",
		DefaultDirectory: defaultDir,
		DefaultFilename:  defaultFilename,
		Filters: []runtime.FileFilter{
			{DisplayName: display, Pattern: pattern},
		},
	})
}
