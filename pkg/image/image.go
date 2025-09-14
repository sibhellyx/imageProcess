package image

import (
	"github.com/sibhellyx/imageProccesor/internal/models"
)

func ProccesImage(id int, task *models.ImageTask) error {
	for _, action := range task.Actions {
		switch action.Type {
		// case models.ActionTypeDownload:
		// 	output, err := actions.DownloadImageWithResty(task.DownloadPath, task.Name)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	fmt.Println(output)
		// 	task.Path = output

		case models.ActionTypeResize:
			// params := action.GetResizeParams()
			// err := actions.Resize(task.Path, params)
			// if err != nil {
			// 	return err
			// }

		}
	}

	return nil
	// return "new" + path
}
