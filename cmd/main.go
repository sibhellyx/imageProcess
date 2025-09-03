package main

import (
	"fmt"

	"github.com/sibhellyx/imageProccesor/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
}
