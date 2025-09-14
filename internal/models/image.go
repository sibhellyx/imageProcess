package models

import (
	"fmt"
	"time"
)

type ImageRequestDownload struct {
	Url  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

func (r *ImageRequestDownload) Validate() error {
	if r.Name == "" || r.Url == "" {
		return fmt.Errorf("url and name are required")
	}
	return nil
}

type ImageRequestAction struct {
	Path    string   `json:"path"`
	Actions []Action `json:"actions"`
}

func (r *ImageRequestAction) Validate() error {
	if r.Path == "" {
		return fmt.Errorf("path is required for proccesing")
	}
	// check valid params for action
	for i, action := range r.Actions {
		valid := action.Type.IsValid()
		if !valid {
			return fmt.Errorf("invalid action type '%s' at position %d", action.Type, i+1)
		}
		// проверка валидности параметров
		err := action.ValidateParams()
		if err != nil {
			return fmt.Errorf("error validating params for action '%s' at position %d: %w",
				action.Type, i+1, err)

		}
	}

	return nil
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
