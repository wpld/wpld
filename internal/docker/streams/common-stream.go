package streams

import (
	"github.com/moby/term"
)

type CommonStream struct {
	fd     uintptr
	isTerm bool
	state  *term.State
}

func (cs CommonStream) FD() uintptr {
	return cs.fd
}

func (cs CommonStream) RestoreTerminal() error {
	if cs.state != nil {
		return term.RestoreTerminal(cs.fd, cs.state)
	} else {
		return nil
	}
}
