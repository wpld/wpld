package utils_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wpld/utils"
)

func TestParseBytes(t *testing.T) {
	each := map[string]int64 {
		"1KiB": 1 << 10,
		"1mib": 1 << 20,
		"1GiB": 1 << 30,
		"128MiB": 128 * 1 << 20,
		"1kb": 1000,
		"1mb": 1000 * 1000,
		"1gb": 1000 * 1000 * 1000,
		"256Kb": 256 * 1000,
	}

	for from, to := range each {
		name := from
		expected := to
		t.Run(name, func(tst *testing.T) {
			actual, err := utils.ParseBytes(name)
			if assert.NoError(tst, err) {
				assert.Equal(tst, expected, actual, name)
			}
		})
	}
}
