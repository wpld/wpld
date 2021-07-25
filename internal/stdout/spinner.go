package stdout

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var s *spinner.Spinner

func init() {
	s = spinner.New(
		spinner.CharSets[14],
		100*time.Millisecond,
		spinner.WithWriter(os.Stderr),
		spinner.WithHiddenCursor(true),
	)
}

func StartSpinner(suffixAndPrefix ...string) {
	if IsTerm() {
		paramsCount := len(suffixAndPrefix)

		if paramsCount > 0 {
			s.Suffix = " " + suffixAndPrefix[0]
		} else {
			s.Suffix = ""
		}

		if paramsCount > 1 {
			s.Prefix = suffixAndPrefix[1] + " "
		} else {
			s.Prefix = ""
		}

		s.Start()
	}
}

func StopSpinner() {
	if IsTerm() && s.Active() {
		s.Prefix = ""
		s.Suffix = ""
		s.Stop()
	}
}
