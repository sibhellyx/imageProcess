package image

import (
	"fmt"
	"os"

	"github.com/sibhellyx/imageProccesor/internal/models"
	"github.com/sibhellyx/imageProccesor/pkg/actions"
	"github.com/sibhellyx/imageProccesor/pkg/utils"
)

func ProccesImage(id int, task *models.ImageTask) error {
	img, err := actions.OpenImage(task.Path)
	if err != nil {
		return fmt.Errorf("error of open image: %w", err)
	}
	var targetFormat string
	for _, action := range task.Actions {
		switch action.Type {
		case models.ActionTypeResize:
			params := action.GetResizeParams()
			img, err = actions.Resize(img, params)
			if err != nil {
				return fmt.Errorf("error of resize image: %w", err)
			}
		case models.ActionTypeConvert:
			params := action.GetConvertParams()
			img = actions.Convert(img, params)
			targetFormat = params.Format
		}
	}
	if targetFormat != "" {
		os.Remove(task.Path)
		task.Path = utils.ChangeExtension(task.Path, targetFormat)
	}
	err = actions.SaveImage(img, task.Path)
	if err != nil {
		return fmt.Errorf("error of save image: %w", err)
	}
	return nil
}
