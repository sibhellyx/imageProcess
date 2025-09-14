package actions

import (
	"image"

	"github.com/disintegration/imaging"
	"github.com/sibhellyx/imageProccesor/internal/models"
)

func Resize(img image.Image, params models.ResizeParams) (image.Image, error) {
	filter := getInterpolation(params.Interpolation)

	dstImg := imaging.Resize(img, params.Width, params.Height, filter)

	return dstImg, nil
}
