package image

import (
	"fmt"

	"github.com/sibhellyx/imageProccesor/internal/models"
	"github.com/sibhellyx/imageProccesor/pkg/actions"
)

func ProccesImage(id int, task *models.ImageTask) error {
	// fmt.Println("Worker ", id, "procces image ", path)
	// time.Sleep(time.Duration(rand.IntN(10)) * time.Second)
	if task.Actions[0].Type == models.ActionTypeDownload {
		output, err := actions.DownloadImageWithResty(task.DownloadPath, task.Name)
		if err != nil {
			return err
		}
		fmt.Println(output)
		task.Path = output
	}
	return nil
	// return "new" + path
}
