package streams

import (
	"errors"
	"io"
	"runtime"

	"github.com/moby/term"
)

type InStream struct {
	CommonStream
	stream io.ReadCloser
}

func NewInStream(stream io.ReadCloser) *InStream {
	fd, isTerm := term.GetFdInfo(stream)

	return &InStream{
		CommonStream: CommonStream{
			fd:     fd,
			isTerm: isTerm,
		},
		stream: stream,
	}
}

func (s InStream) Read(p []byte) (int, error) {
	return s.stream.Read(p)
}

func (s InStream) Close() error {
	return s.stream.Close()
}

func (s InStream) CheckTTY(attachStdin, ttyMode bool) error {
	if ttyMode && attachStdin && !s.CommonStream.isTerm {
		if runtime.GOOS == "windows" {
			return errors.New("the input device is not a TTY.  If you are using mintty, try prefixing the command with 'winpty'")
		} else {
			return errors.New("the input device is not a TTY")
		}
	}

	return nil
}
