package stdout

import (
	"os"

	"github.com/mattn/go-isatty"
)

var isterm bool

func init() {
	isterm = false

	fd := os.Stdout.Fd()
	if isatty.IsTerminal(fd) {
		isterm = true
	} else if isatty.IsCygwinTerminal(fd) {
		isterm = true
	}
}

func IsTerm() bool {
	return isterm
}
