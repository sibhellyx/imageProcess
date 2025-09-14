package actions

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/sibhellyx/imageProccesor/internal/models"
)

func Convert(img image.Image, params models.ConvertParams) image.Image {
	switch params.Format {
	case "jpeg", "jpg":
		if hasTransparency(img) {
			img = imaging.Overlay(
				imaging.New(img.Bounds().Dx(), img.Bounds().Dy(), color.White),
				img,
				image.Point{},
				1.0,
			)
		}
	}
	return img

}

func hasTransparency(img image.Image) bool {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 0xFFFF {
				return true
			}
		}
	}
	return false
}
