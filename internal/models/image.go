package models

import (
	"fmt"
	"time"
)

type ImageRequest struct {
	Url  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`

	Path string `json:"path,omitempty"`

	Actions []Action `json:"actions"`
}

func (r *ImageRequest) Validate() (bool, error) {
	if r.Url == "" && r.Path == "" {
		return false, fmt.Errorf("either URL or file_path must be specified")
	}
	if r.Url != "" && r.Path != "" {
		return false, fmt.Errorf("cannot specify both URL and file_path")
	}
	for _, action := range r.Actions {
		valid := action.Type.IsValid()
		if !valid {
			return false, fmt.Errorf("action type incorrect")
		}
	}
	if r.Url != "" {
		if r.Name == "" {
			return false, fmt.Errorf("name is required when downloading from Url")
		}
		if len(r.Actions) != 1 && r.Actions[0].Type != ActionTypeDownload {
			return false, fmt.Errorf("action must be one and must be downloading when downloading from Url")
		}
		return true, nil
	}

	return false, nil
}

type ImageTask struct {
	Name         string
	DownloadPath string
	Path         string
	Status       StatusImage
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Actions      []Action
}

type ImageTaskAnswer struct {
	Output string
	Err    error
}
