package docker

import (
	"wpld/internal/misc"
)

func GetBasicLabels() map[string]string {
	return map[string]string{
		"io.wpld.version": misc.VERSION,
	}
}
