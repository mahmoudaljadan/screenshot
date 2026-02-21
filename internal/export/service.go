package exporter

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohamoundaljadan/screenshot/internal/annotate"
	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Export(_ context.Context, req core.ExportRequest) (core.ExportResult, error) {
	if err := annotate.ValidateOps(req.Ops); err != nil {
		return core.ExportResult{}, err
	}
	f, err := os.Open(req.BaseImagePath)
	if err != nil {
		return core.ExportResult{}, &core.AppError{Code: core.ErrReadFailed, Message: err.Error()}
	}
	img, _, err := image.Decode(f)
	_ = f.Close()
	if err != nil {
		return core.ExportResult{}, &core.AppError{Code: core.ErrDecodeFailed, Message: err.Error()}
	}

	rgba := image.NewRGBA(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	annotate.SortOps(req.Ops)
	if err := annotate.ApplyOps(rgba, req.Ops); err != nil {
		return core.ExportResult{}, err
	}

	format := strings.ToLower(req.Format)
	if format == "" {
		format = "png"
	}
	if req.OutputPath == "" {
		req.OutputPath = defaultOutputPath(format)
	}
	if err := os.MkdirAll(filepath.Dir(req.OutputPath), 0o755); err != nil {
		return core.ExportResult{}, &core.AppError{Code: core.ErrWriteFailed, Message: err.Error()}
	}

	out, err := os.Create(req.OutputPath)
	if err != nil {
		return core.ExportResult{}, &core.AppError{Code: core.ErrWriteFailed, Message: err.Error()}
	}
	defer out.Close()

	switch format {
	case "png":
		err = png.Encode(out, rgba)
	case "jpg", "jpeg":
		quality := req.Quality
		if quality <= 0 || quality > 100 {
			quality = 90
		}
		err = jpeg.Encode(out, rgba, &jpeg.Options{Quality: quality})
	default:
		return core.ExportResult{}, &core.AppError{Code: core.ErrEncodeFailed, Message: "unsupported format: " + format}
	}
	if err != nil {
		return core.ExportResult{}, &core.AppError{Code: core.ErrEncodeFailed, Message: err.Error()}
	}

	stat, err := out.Stat()
	if err != nil {
		return core.ExportResult{}, &core.AppError{Code: core.ErrReadFailed, Message: err.Error()}
	}
	log.Printf("[export] saved format=%s bytes=%d output=%s", format, stat.Size(), req.OutputPath)
	return core.ExportResult{OutputPath: req.OutputPath, Bytes: stat.Size(), Format: format}, nil
}

func defaultOutputPath(format string) string {
	ext := ".png"
	if format == "jpg" || format == "jpeg" {
		ext = ".jpg"
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		ts := time.Now().Format("20060102-150405")
		return filepath.Join(os.TempDir(), fmt.Sprintf("go-wails-shot-%s%s", ts, ext))
	}
	dir := filepath.Join(home, "Pictures", "go-wails-shot")
	ts := time.Now().Format("20060102-150405")
	return filepath.Join(dir, fmt.Sprintf("capture-%s%s", ts, ext))
}
