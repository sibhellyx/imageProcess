package models

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
