package annotate

import (
	"encoding/json"
	"sort"

	"github.com/mohamoundaljadan/screenshot/internal/core"
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type RectPayload struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	W           int    `json:"w"`
	H           int    `json:"h"`
	Color       string `json:"color"`
	StrokeWidth int    `json:"strokeWidth"`
	Fill        bool   `json:"fill"`
}

type LinePayload struct {
	X1          int    `json:"x1"`
	Y1          int    `json:"y1"`
	X2          int    `json:"x2"`
	Y2          int    `json:"y2"`
	Color       string `json:"color"`
	StrokeWidth int    `json:"strokeWidth"`
}

type ArrowPayload struct {
	LinePayload
	HeadSize int `json:"headSize"`
}

type TextPayload struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Text  string `json:"text"`
	Color string `json:"color"`
	Size  int    `json:"size"`
}

type BlurPayload struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	W      int `json:"w"`
	H      int `json:"h"`
	Radius int `json:"radius"`
}

type PixelatePayload struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	W    int `json:"w"`
	H    int `json:"h"`
	Size int `json:"size"`
}

var knownKinds = map[string]struct{}{
	"rect":     {},
	"line":     {},
	"arrow":    {},
	"text":     {},
	"blur":     {},
	"pixelate": {},
}

func ValidateOps(ops []core.AnnotationOp) error {
	for _, op := range ops {
		if _, ok := knownKinds[op.Kind]; !ok {
			return &core.AppError{Code: core.ErrInvalidOpKind, Message: "unsupported op kind: " + op.Kind}
		}
		if err := validatePayload(op); err != nil {
			return err
		}
	}
	return nil
}

func SortOps(ops []core.AnnotationOp) {
	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Z == ops[j].Z {
			return ops[i].ID < ops[j].ID
		}
		return ops[i].Z < ops[j].Z
	})
}

func validatePayload(op core.AnnotationOp) error {
	if len(op.Payload) == 0 {
		return &core.AppError{Code: core.ErrInvalidOpPayload, Message: "empty payload for op: " + op.ID}
	}
	switch op.Kind {
	case "rect":
		var p RectPayload
		if err := json.Unmarshal(op.Payload, &p); err != nil {
			return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
		}
	case "line":
		var p LinePayload
		if err := json.Unmarshal(op.Payload, &p); err != nil {
			return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
		}
	case "arrow":
		var p ArrowPayload
		if err := json.Unmarshal(op.Payload, &p); err != nil {
			return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
		}
	case "text":
		var p TextPayload
		if err := json.Unmarshal(op.Payload, &p); err != nil {
			return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
		}
	case "blur":
		var p BlurPayload
		if err := json.Unmarshal(op.Payload, &p); err != nil {
			return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
		}
	case "pixelate":
		var p PixelatePayload
		if err := json.Unmarshal(op.Payload, &p); err != nil {
			return &core.AppError{Code: core.ErrInvalidOpPayload, Message: err.Error()}
		}
	}
	return nil
}
