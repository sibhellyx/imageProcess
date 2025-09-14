package models

import (
	"fmt"
	"strings"
)

type ConvertParams struct {
	Format string
}

func (a Action) ValidateConvertParams() error {
	format, ok := a.Params["format"]
	if !ok {
		return fmt.Errorf("format is required")
	}
	err := validateFormat(format.(string))
	if err != nil {
		return err
	}
	return nil
}

func (a Action) GetConvertParams() ConvertParams {
	return ConvertParams{
		Format: strings.ToLower(a.Params["format"].(string)),
	}
}

func validateFormat(format string) error {
	switch strings.ToLower(format) {
	case "jpeg", "jpg", "png", "webp", "gif", "bmp", "tiff", "tif":
		return nil
	default:
		return fmt.Errorf("unknown format")
	}
}
