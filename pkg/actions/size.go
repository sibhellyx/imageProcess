package actions

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
	"github.com/sibhellyx/imageProccesor/internal/models"
)

func Resize(path string, params models.ResizeParams) (image.Image, error) {
	img, err := openImage(path)
	if err != nil {
		return nil, fmt.Errorf("error of open image: %w", err)
	}

	filter := getInterpolation(params.Interpolation)

	dstImg := imaging.Resize(img, params.Width, params.Height, filter)

	return dstImg, nil
}
