package errors

import "errors"

var (
	ErrQueueTasksFull     = errors.New("task queue for proccesing full")
	ErrServerShuttingDown = errors.New("server shuttingdown")
	ErrShutdownTimeout    = errors.New("service shutdown timeout was limited")
)
