package actions

import (
	"image"

	"github.com/disintegration/imaging"
)

const baseDir = "./downloads"

// func for open img.
func openImage(path string) (image.Image, error) {
	img, err := imaging.Open(baseDir + "/" + path)
	if err != nil {
		return nil, err
	}
	return img, nil
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
