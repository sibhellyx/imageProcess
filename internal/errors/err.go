package errors

import "errors"

var (
	ErrQueueTasksFull = errors.New("task queue for proccesing full")
)
