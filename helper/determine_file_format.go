package helper

import (
	"strings"

	"github.com/funmi4194/ecommerce/primer"
)

func DetermineFileFormat(name string) string {
	extension := ExtractExtension(name)
	if format, ok := primer.ImageFormat[extension]; ok {
		return strings.ToUpper(format.String())
	}
	return ""
}
