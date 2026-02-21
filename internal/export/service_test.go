package exporter

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

func TestExportDeterministic(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "base.png")
	writeBaseImage(t, base)

	req := core.ExportRequest{
		BaseImagePath: base,
		Format:        "png",
		OutputPath:    filepath.Join(tmp, "annotated.png"),
		Ops: []core.AnnotationOp{
			{ID: "b", Kind: "line", Z: 2, Payload: json.RawMessage(`{"x1":10,"y1":10,"x2":60,"y2":60,"color":"#00ff00","strokeWidth":3}`)},
			{ID: "a", Kind: "rect", Z: 1, Payload: json.RawMessage(`{"x":20,"y":15,"w":40,"h":30,"color":"#ff0000","strokeWidth":2}`)},
		},
	}

	svc := NewService()
	result, err := svc.Export(context.Background(), req)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if result.OutputPath == "" {
		t.Fatal("expected output path")
	}

	h1 := hashFile(t, result.OutputPath)
	result2, err := svc.Export(context.Background(), req)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	h2 := hashFile(t, result2.OutputPath)
	if h1 != h2 {
		t.Fatalf("expected deterministic output hash, got %s vs %s", h1, h2)
	}
}

func writeBaseImage(t *testing.T, p string) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 80, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 80; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x * 3), G: uint8(y * 3), B: 128, A: 255})
		}
	}
	f, err := os.Create(p)
	if err != nil {
		t.Fatalf("create base: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("encode base: %v", err)
	}
}

func hashFile(t *testing.T, p string) string {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read hash file: %v", err)
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
