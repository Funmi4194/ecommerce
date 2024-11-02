package helper

import (
	"path/filepath"
	"strings"
)

func ExtractExtension(path string) string {
	return strings.ToLower(filepath.Ext(path)[1:])
}
