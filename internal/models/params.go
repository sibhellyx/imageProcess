package models

const DefaultInterpolation = "Lanczos"

type ResizeParams struct {
	Width         int
	Height        int
	Interpolation string
}
