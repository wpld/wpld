package streams

import (
	"io"
	"os"

	"github.com/moby/term"
)

type OutStream struct {
	CommonStream
	stream io.Writer
}

func NewOutStream(stream io.Writer) *OutStream {
	fd, isTerm := term.GetFdInfo(stream)

	return &OutStream{
		CommonStream: CommonStream{
			fd:     fd,
			isTerm: isTerm,
		},
		stream: stream,
	}
}

func (s OutStream) Write(p []byte) (int, error) {
	return s.stream.Write(p)
}

func (s *OutStream) SetRawTerminal() error {
	if os.Getenv("NORAW") != "" || !s.CommonStream.isTerm {
		return nil
	}

	if state, err := term.SetRawTerminalOutput(s.CommonStream.fd); err != nil {
		return err
	} else {
		s.CommonStream.state = state
	}

	return nil
}
