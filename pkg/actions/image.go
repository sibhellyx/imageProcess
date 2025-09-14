package actions

import (
	"fmt"
	"image"

	"github.com/disintegration/imaging"
)

const baseDir = "./downloads"

// func for open img.
func OpenImage(path string) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// func for save img
func SaveImage(img image.Image, path string) error {
	err := imaging.Save(img, path)
	if err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}
	return nil
}

// get type interpolation.
func getInterpolation(typeInterpolation string) imaging.ResampleFilter {
	switch typeInterpolation {
	case "NearestNeighbor":
		return imaging.NearestNeighbor
	case "Liner":
		return imaging.Linear
	case "MitchelNetravali":
		return imaging.MitchellNetravali
	case "CatmullRom":
		return imaging.CatmullRom
	case "Lanczos":
		return imaging.Lanczos
	default:
		return imaging.Lanczos
	}
}
