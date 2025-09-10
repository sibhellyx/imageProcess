package pool

// import (
// 	"fmt"
// 	"strconv"
// 	"sync"
// 	"testing"
// 	"time"
// )

// func mockHandler(workerId int, imagePath string) string {
// 	time.Sleep(10 * time.Millisecond)
// 	return "processed_" + imagePath + "_by_" + strconv.Itoa(workerId)

// }

// func TestNewPool(t *testing.T) {
// 	pool := NewPool(mockHandler, 5)

// 	if pool == nil {
// 		t.Errorf("Expected pool was created")
// 	}

// 	if len(pool.workers) != 5 {
// 		t.Errorf("Expected 5 workers, got %d", len(pool.workers))
// 	}

// 	if cap(pool.pool) != 5 {
// 		t.Errorf("Expected pool channel capacity 5, got %d", cap(pool.pool))
// 	}

// }

// func TestPool_Create(t *testing.T) {
// 	pool := NewPool(mockHandler, 2)
// 	pool.Create()

// 	if len(pool.pool) != 2 {
// 		t.Errorf("Expected 2 workers in pool, got %d", len(pool.pool))
// 	}
// }

// func TestPool_Handle(t *testing.T) {
// 	pool := NewPool(mockHandler, 2)
// 	pool.Create()

// 	resultChan := pool.Handle("test_image.jpg")

// 	select {
// 	case res := <-resultChan:
// 		expectedPrefix := "processed_test_image.jpg_by_"
// 		if res[:len(expectedPrefix)] != expectedPrefix {
// 			t.Errorf("Expected result to start with %s, got %s", expectedPrefix, res)
// 		}
// 	case <-time.After(100 * time.Millisecond):
// 		t.Error("Timeout waiting for result")
// 	}

// }

// func TestPool_Handle_Shutdown(t *testing.T) {
// 	pool := NewPool(mockHandler, 2)
// 	pool.Create()

// 	pool.Shutdown()
// 	resultChan := pool.Handle("test_image.jpg")

// 	select {
// 	case _, ok := <-resultChan:
// 		if ok {
// 			t.Error("Expected channel to be closed after shutdown")
// 		}
// 	case <-time.After(15 * time.Millisecond):
// 		t.Error("Timeout waiting for result")
// 	}
// }

// func TestPool_ConcurrentHandling(t *testing.T) {
// 	pool := NewPool(mockHandler, 2)
// 	pool.Create()

// 	var mu sync.Mutex
// 	var wg sync.WaitGroup

// 	results := make([]string, 0)

// 	// Запускаем несколько concurrent задач
// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func(index int) {
// 			defer wg.Done()

// 			resultChan := pool.Handle(fmt.Sprintf("image_%d.jpg", index))
// 			result := <-resultChan

// 			mu.Lock()
// 			results = append(results, result)
// 			mu.Unlock()
// 		}(i)
// 	}

// 	wg.Wait()

// 	if len(results) != 10 {
// 		t.Errorf("Expected 10 results, got %d", len(results))
// 	}

// }

// func TestPool_Shutdown(t *testing.T) {
// 	pool := NewPool(mockHandler, 2)
// 	pool.Create()

// 	for i := 0; i < 5; i++ {
// 		go func(index int) {
// 			pool.Handle(fmt.Sprintf("image_%d.jpg", index))
// 		}(i)
// 	}

// 	time.Sleep(10 * time.Millisecond)

// 	pool.Shutdown()

// 	if !pool.shutdown {
// 		t.Error("Expected pool to be marked as shutdown")
// 	}

// 	resultChan := pool.Handle("new_image.jpg")
// 	select {
// 	case _, ok := <-resultChan:
// 		if ok {
// 			t.Error("Expected closed channel for new tasks after shutdown")
// 		}
// 	case <-time.After(10 * time.Millisecond):
// 		t.Error("Expected immediate response")
// 	}
// }
