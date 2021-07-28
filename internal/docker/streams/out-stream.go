package streams

import (
	"io"

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
