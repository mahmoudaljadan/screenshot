package app

import (
	"context"
	"runtime"

	capture "github.com/mohamoundaljadan/screenshot/internal/capture"
	exporter "github.com/mohamoundaljadan/screenshot/internal/export"
)

type Service struct {
	capture *capture.Manager
	export  *exporter.Service
	prefs   *PreferenceStore
	version string
}

func NewService(captureManager *capture.Manager, exportService *exporter.Service, prefs *PreferenceStore, version string) *Service {
	return &Service{capture: captureManager, export: exportService, prefs: prefs, version: version}
}

func (s *Service) StartCapture(mode string) (CaptureResult, error) {
	return s.capture.Capture(context.Background(), mode)
}

func (s *Service) SaveAnnotated(req ExportRequest) (ExportResult, error) {
	return s.export.Export(context.Background(), req)
}

func (s *Service) GetAppState() (AppState, error) {
	sessionType, waylandBeta := s.capture.SessionInfo()
	return AppState{
		Platform:    runtime.GOOS,
		SessionType: sessionType,
		WaylandBeta: waylandBeta,
		Version:     s.version,
	}, nil
}

func (s *Service) SetPreference(key string, value string) error {
	return s.prefs.Set(key, value)
}
