package annotate

import (
	"encoding/json"
	"testing"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

func TestValidateOpsRejectsUnknownKind(t *testing.T) {
	op := core.AnnotationOp{ID: "1", Kind: "circle", Payload: json.RawMessage(`{"x":1}`)}
	if err := ValidateOps([]core.AnnotationOp{op}); err == nil {
		t.Fatal("expected error for unknown op kind")
	}
}

func TestValidateOpsAcceptsKnownKinds(t *testing.T) {
	ops := []core.AnnotationOp{
		{ID: "1", Kind: "rect", Payload: json.RawMessage(`{"x":1,"y":2,"w":3,"h":4,"color":"#ff0000","strokeWidth":2}`)},
		{ID: "2", Kind: "line", Payload: json.RawMessage(`{"x1":1,"y1":2,"x2":3,"y2":4,"color":"#ff0000","strokeWidth":2}`)},
		{ID: "3", Kind: "arrow", Payload: json.RawMessage(`{"x1":1,"y1":2,"x2":3,"y2":4,"color":"#ff0000","strokeWidth":2,"headSize":12}`)},
		{ID: "4", Kind: "text", Payload: json.RawMessage(`{"x":1,"y":2,"text":"abc","color":"#ff0000","size":12}`)},
		{ID: "5", Kind: "blur", Payload: json.RawMessage(`{"x":1,"y":2,"w":3,"h":4,"radius":2}`)},
		{ID: "6", Kind: "pixelate", Payload: json.RawMessage(`{"x":1,"y":2,"w":3,"h":4,"size":8}`)},
	}
	if err := ValidateOps(ops); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
