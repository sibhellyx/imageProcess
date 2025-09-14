package utils

import (
	"fmt"
	"path/filepath"
	"strconv"
)

func ConverToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// ChangeExtension изменяет расширение файла в пути
func ChangeExtension(path, newExt string) string {
	oldExt := filepath.Ext(path)
	return path[:len(path)-len(oldExt)] + "." + newExt
}
