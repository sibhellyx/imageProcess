package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

func DownloadImageWithResty(url, filename string) (string, error) {
	client := resty.New()

	tempFile := filepath.Join(baseDir, filename+".tmp")

	resp, err := client.R().
		SetOutput(tempFile).
		SetHeader("User-Agent", "Mozilla/5.0 (compatible; Go-Image-Downloader/1.0)").
		Get(url)

	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("HTTP error: %s", resp.Status())
	}

	contentType := resp.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		os.Remove(tempFile)
		return "", fmt.Errorf("not an image: %s", contentType)
	}

	extension := getExtensionFromContentType(contentType)
	if extension == "" {
		os.Remove(tempFile)
		return "", fmt.Errorf("unsupported image type: %s", contentType)
	}

	finalFileName := filename + extension
	outputPath := filepath.Join(baseDir, finalFileName)

	err = os.Rename(tempFile, outputPath)
	if err != nil {
		os.Remove(tempFile)
		return "", fmt.Errorf("failed to rename file: %w", err)
	}

	return outputPath, nil
}

func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	case "image/bmp":
		return ".bmp"
	case "image/tiff":
		return ".tiff"
	default:
		// Пытаемся извлечь расширение из типа, если он нестандартный
		if strings.HasPrefix(contentType, "image/") {
			parts := strings.Split(contentType, "/")
			if len(parts) > 1 {
				return "." + parts[1]
			}
		}
		return ""
	}
}
