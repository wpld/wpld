package docker

import (
	"os"
	"path/filepath"
	"strings"
)

func NormalizeContainerBinds(binds []string) []string {
	normalized := make([]string, len(binds))

	for i, bind := range binds {
		parts := strings.SplitN(bind, ":", 2)
		if !filepath.IsAbs(parts[0]) {
			if abs, err := filepath.Abs(parts[0]); err == nil {
				if _, statErr := os.Stat(abs); !os.IsNotExist(statErr) {
					parts[0] = abs
				}
			}
		}

		normalized[i] = strings.Join(parts, ":")
	}

	return normalized
}
