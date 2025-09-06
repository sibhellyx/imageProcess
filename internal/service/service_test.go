package service

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/sibhellyx/imageProccesor/internal/errors"
	"github.com/sibhellyx/imageProccesor/internal/repository"
)

type mockWorkerPool struct {
	mu             sync.Mutex
	created        bool
	handled        []string
	shutdownCalled bool
	handleDelay    time.Duration
}

func (m *mockWorkerPool) Create() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.created = true
}

func (m *mockWorkerPool) Handle(imagePath string) <-chan string {
	m.mu.Lock()
	m.handled = append(m.handled, imagePath)
	delay := m.handleDelay
	m.mu.Unlock()

	resultChan := make(chan string, 1)
	go func() {
		if delay > 0 {
			time.Sleep(delay)
		}
		resultChan <- "processed: " + imagePath
		close(resultChan)
	}()
	return resultChan
}

func (m *mockWorkerPool) Wait() {}

func (m *mockWorkerPool) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shutdownCalled = true
}

func (m *mockWorkerPool) Stats() {}

func (m *mockWorkerPool) getHandled() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string{}, m.handled...) // Возвращаем копию
}

func (m *mockWorkerPool) isCreated() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.created
}

func (m *mockWorkerPool) isShutdown() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.shutdownCalled
}

func (m *mockWorkerPool) setHandleDelay(delay time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handleDelay = delay
}

type mockRepository struct {
	mu    sync.Mutex
	paths []string
}

func (m *mockRepository) Save(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.paths = append(m.paths, path)
	return nil
}

func (m *mockRepository) GetPaths() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string{}, m.paths...)
}

func createMockRepo() *repository.Repository {
	// Возвращаем nil, так как в текущей реализации репозиторий не используется
	return nil
}

func TestNewService(t *testing.T) {
	mockPool := &mockWorkerPool{}
	mockRepo := createMockRepo()

	service := NewService(mockRepo, mockPool, 10)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		service.Shutdown(ctx)
	}()

	if service == nil {
		t.Fatal("Service should not be nil")
	}

	if !mockPool.isCreated() {
		t.Error("WorkerPool.Create() should be called")
	}

	// Проверяем, что сервис не в состоянии shutdown
	service.mu.Lock()
	isShutdown := service.shutdown
	service.mu.Unlock()

	if isShutdown {
		t.Error("Service should not be in shutdown state initially")
	}
}

func TestService_Create(t *testing.T) {
	t.Run("successful task creation", func(t *testing.T) {
		mockPool := &mockWorkerPool{}
		mockRepo := createMockRepo()

		service := NewService(mockRepo, mockPool, 10)
		defer shutdownService(t, service)

		err := service.Create(context.Background(), "image1.jpg")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// Даем время для обработки
		time.Sleep(100 * time.Millisecond)

		handled := mockPool.getHandled()
		if len(handled) != 1 {
			t.Errorf("Expected 1 handled task, got %d", len(handled))
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		mockPool := &mockWorkerPool{}
		mockRepo := createMockRepo()

		service := NewService(mockRepo, mockPool, 1) // Маленькая очередь
		defer shutdownService(t, service)

		// Заполняем очередь
		service.Create(context.Background(), "image1.jpg")

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := service.Create(ctx, "image2.jpg")
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled, got %v", err)
		}
	})

	t.Run("service shutdown", func(t *testing.T) {
		mockPool := &mockWorkerPool{}
		mockRepo := createMockRepo()

		service := NewService(mockRepo, mockPool, 10)

		// Вызываем shutdown в отдельной горутине
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			service.Shutdown(ctx)
		}()

		// Ждем немного чтобы shutdown успел начаться
		time.Sleep(50 * time.Millisecond)

		err := service.Create(context.Background(), "image1.jpg")
		if err != errors.ErrServerShuttingDown {
			t.Errorf("Expected ErrServerShuttingDown, got %v", err)
		}
	})

	t.Run("queue full", func(t *testing.T) {
		mockPool := &mockWorkerPool{}
		mockPool.setHandleDelay(200 * time.Millisecond) // Замедляем обработку
		mockRepo := createMockRepo()

		service := NewService(mockRepo, mockPool, 1) // Очень маленькая очередь
		defer shutdownService(t, service)

		// Первая задача должна пройти
		err := service.Create(context.Background(), "image1.jpg")
		if err != nil {
			t.Errorf("First task should succeed, got %v", err)
		}

		// Вторая задача должна получить ошибку переполнения
		err = service.Create(context.Background(), "image2.jpg")
		if err != errors.ErrQueueTasksFull {
			t.Errorf("Expected ErrQueueTasksFull, got %v", err)
		}
	})
}

func TestService_Shutdown(t *testing.T) {
	t.Run("graceful shutdown", func(t *testing.T) {
		mockPool := &mockWorkerPool{}
		mockRepo := createMockRepo()
		service := NewService(mockRepo, mockPool, 10)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := service.Shutdown(ctx)
		if err != nil {
			t.Errorf("Expected no error during shutdown, got %v", err)
		}

		if !mockPool.isShutdown() {
			t.Error("WorkerPool.Shutdown() should be called")
		}

		// Проверяем, что сервис помечен как shutdown
		service.mu.Lock()
		isShutdown := service.shutdown
		service.mu.Unlock()

		if !isShutdown {
			t.Error("Service should be marked as shutdown")
		}
	})

	t.Run("shutdown timeout", func(t *testing.T) {
		mockPool := &mockWorkerPool{}
		mockPool.setHandleDelay(2 * time.Second) // Долгая обработка
		mockRepo := createMockRepo()
		service := NewService(mockRepo, mockPool, 10)

		// Добавляем задачу
		service.Create(context.Background(), "slow-image.jpg")

		// Даем время задаче начать обрабатываться
		time.Sleep(100 * time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		err := service.Shutdown(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("Expected DeadlineExceeded, got %v", err)
		}
	})
}

func TestService_Integration(t *testing.T) {
	mockPool := &mockWorkerPool{}
	mockRepo := createMockRepo()
	service := NewService(mockRepo, mockPool, 10)
	defer shutdownService(t, service)

	paths := []string{"image1.jpg", "image2.jpg", "image3.jpg"}

	for _, path := range paths {
		err := service.Create(context.Background(), path)
		if err != nil {
			t.Errorf("Failed to create task for %s: %v", path, err)
		}
	}

	// Ждем завершения обработки
	time.Sleep(500 * time.Millisecond)

	handled := mockPool.getHandled()
	if len(handled) != len(paths) {
		t.Errorf("Expected %d handled tasks, got %d: %v", len(paths), len(handled), handled)
	}

	// Проверяем, что все задачи были обработаны
	for _, path := range paths {
		found := false
		for _, handledPath := range handled {
			if handledPath == path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Path %s was not handled", path)
		}
	}
}

func TestService_Concurrent(t *testing.T) {
	mockPool := &mockWorkerPool{}
	mockRepo := createMockRepo()
	service := NewService(mockRepo, mockPool, 1)
	defer shutdownService(t, service)

	var wg sync.WaitGroup
	const numTasks = 50

	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			err := service.Create(context.Background(), fmt.Sprintf("image_%d.jpg", index))
			if err != nil && err != errors.ErrQueueTasksFull {
				t.Errorf("Unexpected error for task %d: %v", index, err)
			}
		}(i)
	}

	wg.Wait()

	// Ждем завершения обработки
	time.Sleep(1 * time.Second)

	handled := mockPool.getHandled()
	t.Logf("Handled %d tasks out of %d", len(handled), numTasks)

	if len(handled) == 0 {
		t.Error("Expected at least some tasks to be handled")
	}
}

func TestService_ImmediateQueueFull(t *testing.T) {
	mockPool := &mockWorkerPool{}
	mockPool.setHandleDelay(100 * time.Millisecond)
	mockRepo := createMockRepo()

	// Сервис с очень маленькой очередью и быстрой обработкой
	service := NewService(mockRepo, mockPool, 2)
	defer shutdownService(t, service)

	// Быстро отправляем много задач
	var queueFullCount int
	for i := 0; i < 10; i++ {
		err := service.Create(context.Background(), fmt.Sprintf("image_%d.jpg", i))
		if err == errors.ErrQueueTasksFull {
			queueFullCount++
		} else if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Небольшая задержка между отправками
		time.Sleep(10 * time.Millisecond)
	}

	t.Logf("Queue full errors: %d", queueFullCount)

	if queueFullCount == 0 {
		t.Error("Expected at least some queue full errors")
	}

	// Даем время на обработку
	time.Sleep(500 * time.Millisecond)

	handled := mockPool.getHandled()
	t.Logf("Handled tasks: %v", handled)
}

// Вспомогательная функция для graceful shutdown
func shutdownService(t *testing.T, service *Service) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := service.Shutdown(ctx)
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Unexpected error during shutdown: %v", err)
	}
}
