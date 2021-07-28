package streams

import (
	"os"

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

func (cs CommonStream) SetRawTerminal() error {
	if os.Getenv("NORAW") != "" || !cs.isTerm {
		return nil
	}

	if state, err := term.MakeRaw(cs.fd); err != nil {
		return err
	} else {
		cs.state = state
	}

	return nil
}

func (cs CommonStream) RestoreTerminal() error {
	if cs.state != nil {
		return term.RestoreTerminal(cs.fd, cs.state)
	} else {
		return nil
	}
}
