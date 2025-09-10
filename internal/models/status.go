package models

type StatusImage string

const (
	StatusDownloading StatusImage = "Downloading"
	StatusPending     StatusImage = "Pending"
	StatusProcessing  StatusImage = "Proccesing"
	StatusCompleted   StatusImage = "Completed"
	StatusFailed      StatusImage = "Failed"
	StatusUnknow      StatusImage = "Unknown"
)

func (s StatusImage) IsValid() bool {
	switch s {
	case StatusDownloading,
		StatusPending,
		StatusProcessing,
		StatusCompleted,
		StatusFailed:
		return true
	default:
		return false
	}
}

func ParseStatus(status string) StatusImage {
	switch status {
	case "Pending":
		return StatusPending
	case "Downloading":
		return StatusDownloading
	case "Processing":
		return StatusProcessing
	case "Completed":
		return StatusCompleted
	case "Failed":
		return StatusFailed
	default:
		return StatusUnknow
	}
}
