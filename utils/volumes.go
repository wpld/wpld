package utils

import (
	"path/filepath"
	"strings"
)

func NormalizeVolumeBind(volume string) string {
	parts := strings.SplitN(volume, ":", 2)
	if !filepath.IsAbs(parts[0]) {
		if abs, err := filepath.Abs(parts[0]); err == nil {
			parts[0] = abs
		}
	}

	return strings.Join(parts, ":")
}
