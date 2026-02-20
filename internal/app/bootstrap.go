package app

import (
	capture "github.com/mohamoundaljadan/screenshot/internal/capture"
	exporter "github.com/mohamoundaljadan/screenshot/internal/export"
)

func NewDefaultService(version string) (*Service, error) {
	captureManager := capture.NewManager("")
	exportService := exporter.NewService()
	prefs, err := NewPreferenceStore("")
	if err != nil {
		return nil, err
	}
	return NewService(captureManager, exportService, prefs, version), nil
}
