package utils_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"wpld/utils"
)

func TestSlugify(t *testing.T) {
	each := map[string]string {
		"1KiB": "1kib",
		"!!!MY SITE!!!": "my-site",
		"My site!!": "my-site",
		"Test233 message": "test233-message",
	}

	for from, to := range each {
		name := from
		expected := to
		t.Run(name, func(tst *testing.T) {
			actual := utils.Slugify(name)
			assert.Equal(tst, expected, actual, name)
		})
	}
}
