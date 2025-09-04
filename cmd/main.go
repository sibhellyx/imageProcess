package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/sibhellyx/imageProccesor/config"
	"github.com/sibhellyx/imageProccesor/internal/server"
	"github.com/sibhellyx/imageProccesor/internal/workerpool/pool"
)

func main() {
	// WorkerPoolTest()

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	server := server.NewServer(ctx, cfg)
	go func() {
		<-sigChan
		fmt.Println("\nПолучен сигнал завершения, начинаю graceful shutdown...")
		server.Shutdown()
		cancel()
	}()

	server.Serve()
}

func proccesImage(workerId int, imagePath string) {
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	fmt.Println("Worker ", workerId, "procces image ", imagePath)
}

func getImages() []string {
	iCount := rand.Intn(100)

	images := make([]string, 0, iCount)
	for i := 0; i < iCount; i++ {
		images = append(images, strconv.Itoa(i)+"ImagePath")
	}

	return images
}

func WorkerPoolTest() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	pool := pool.NewPool(proccesImage, 5)

l:
	for {
		select {
		case <-ctx.Done():
			break l
		default:
		}

		images := getImages()

		pool.Create()

		for _, image := range images {
			pool.Handle(image)
		}

		pool.Wait()
	}
	pool.Stats()
}
