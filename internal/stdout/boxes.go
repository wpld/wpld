package stdout

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func Box(lines []string, formatted []string) {
	longest := 0
	for _, line := range lines {
		length := len(line)
		if length > longest {
			longest = length
		}
	}

	if longest == 0 {
		return
	}

	padding := 5
	fullLen := longest + padding*2

	pipeLine := color.MagentaString("│")
	emptyLine := pipeLine + strings.Repeat(" ", fullLen) + pipeLine

	fmt.Println(color.MagentaString("╭%s╮", strings.Repeat("─", fullLen)))
	fmt.Println(emptyLine)

	for i, line := range formatted {
		fmt.Print(
			color.MagentaString("│"),
			strings.Repeat(" ", padding),
			line,
			strings.Repeat(" ", padding+longest-len(lines[i])),
			color.MagentaString("│"),
			"\n",
		)
	}

	fmt.Println(emptyLine)
	fmt.Println(color.MagentaString("╰%s╯", strings.Repeat("─", fullLen)))
}
