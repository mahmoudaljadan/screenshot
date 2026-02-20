package main

import (
	"context"

	appsvc "github.com/mohamoundaljadan/screenshot/internal/app"
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
}

func (a *App) StartCapture(mode string) (appsvc.CaptureResult, error) {
	return a.svc.StartCapture(mode)
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
