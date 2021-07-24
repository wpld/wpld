package stdout

import (
	"fmt"

	"github.com/fatih/color"
)

func Success(msg string) {
	if IsTerm() {
		fmt.Print(color.GreenString("✔"), " ", msg, "\n")
	} else {
		fmt.Print("√ ", msg, "\n")
	}
}

func Error(msg string) {
	if IsTerm() {
		fmt.Print(color.RedString("✖"), " ", msg, "\n")
	} else {
		fmt.Print("× ", msg, "\n")
	}
}
