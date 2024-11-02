package primer

import "github.com/opensaucerer/imgconv"

// ImageFormat is a map of supported image formats.
var ImageFormat = map[string]imgconv.Format{
	"jpeg": imgconv.JPEG,
	"jpg":  imgconv.JPEG,
	"png":  imgconv.PNG,
}
