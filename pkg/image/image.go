package image

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func ProccesImage(id int, path string) string {
	fmt.Println("Worker ", id, "procces image ", path)
	time.Sleep(time.Duration(rand.IntN(10)) * time.Second)
	return "new" + path
}
