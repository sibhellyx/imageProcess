package image

import (
	"fmt"

	"github.com/sibhellyx/imageProccesor/internal/models"
	"github.com/sibhellyx/imageProccesor/pkg/actions"
)

func ProccesImage(id int, task *models.ImageTask) error {
	img, err := actions.OpenImage(task.Path)
	if err != nil {
		return fmt.Errorf("error of open image: %w", err)
	}
	for _, action := range task.Actions {
		switch action.Type {
		case models.ActionTypeResize:
			params := action.GetResizeParams()
			img, err = actions.Resize(img, params)
			if err != nil {
				return fmt.Errorf("error of resize image: %w", err)
			}
		}
	}
	err = actions.SaveImage(img, task.Path)
	if err != nil {
		return fmt.Errorf("error of save image: %w", err)
	}
	return nil
}
