package models

import (
	"fmt"

	"github.com/sibhellyx/imageProccesor/pkg/utils"
)

type ActionType string

const (
	ActionTypeResize    ActionType = "resize"
	ActionTypeConvert   ActionType = "convert"
	ActionTypeWatermark ActionType = "watermark"
	ActionTypeCrop      ActionType = "crop"
	ActionTypeRotate    ActionType = "rotate"
	ActionTypeCompress  ActionType = "compress"
)

type Action struct {
	Type   ActionType             `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}

func (a ActionType) IsValid() bool {
	switch a {
	case ActionTypeResize, ActionTypeConvert, ActionTypeWatermark,
		ActionTypeCrop, ActionTypeRotate, ActionTypeCompress:
		return true
	default:
		return false
	}
}

func ParseActionType(s string) (ActionType, bool) {
	action := ActionType(s)
	return action, action.IsValid()
}

func (a ActionType) String() string {
	return string(a)
}

func (a Action) ValidateParams() error {
	switch a.Type {
	case ActionTypeResize:
		return a.ValiateResizeParams()
	default:
		return nil
	}
}

// ------------------------//
// Params for resize
// width
// height
// interpolationType:
//   - "NearestNeighbor"
//   - "Liner"
//   - "MitchelNetravali"
//   - "CatmullRom"
//   - "Lanczos"
//
// ------------------------//
func (a Action) ValiateResizeParams() error {
	width, ok := a.Params["width"]
	if !ok {
		return fmt.Errorf("width is required")
	}

	height, ok := a.Params["height"]
	if !ok {
		return fmt.Errorf("height is required")
	}

	widthInt, err := utils.ConverToInt(width)
	if err != nil {
		return fmt.Errorf("width incorrect, err: %w", err)
	}

	heightInt, err := utils.ConverToInt(height)
	if err != nil {
		return fmt.Errorf("width in correct, err: %w", err)
	}

	if widthInt <= 0 {
		return fmt.Errorf("width must be positive")
	}
	if heightInt <= 0 {
		return fmt.Errorf("height must be positive")
	}

	return nil
}

func (a Action) GetResizeParams() ResizeParams {
	widthInt, _ := utils.ConverToInt(a.Params["width"])
	heightInt, _ := utils.ConverToInt(a.Params["height"])

	interpolationType := DefaultInterpolation
	interpolation, ok := a.Params["interpolationType"].(string)
	if ok {
		interpolationType = interpolation
	}

	return ResizeParams{
		Width:         widthInt,
		Height:        heightInt,
		Interpolation: interpolationType,
	}
}
