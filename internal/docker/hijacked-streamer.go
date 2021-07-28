package docker

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/ioutils"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/term"

	"wpld/internal/docker/streams"
)

// see: https://github.com/docker/cli/blob/master/cli/command/container/hijack.go

type HijackedStreamer struct {
	in  *streams.InStream
	out *streams.OutStream
	err io.Writer

	stdin  io.ReadCloser
	stdout io.Writer
	stderr io.Writer

	resp types.HijackedResponse
	tty  bool
}

func NewHijackedStreamer(resp types.HijackedResponse, tty bool) *HijackedStreamer {
	stdin, stdout, stderr := term.StdStreams()

	return &HijackedStreamer{
		in:  streams.NewInStream(stdin),
		out: streams.NewOutStream(stdout),
		err: stderr,

		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,

		resp: resp,
		tty:  tty,
	}
}

func (hs HijackedStreamer) Stream(ctx context.Context) error {
	restoreInput, err := hs.SetupInput()
	if err != nil {
		return err
	} else {
		defer restoreInput()
	}

	inputDone, inputDetached := hs.BeginInputStream(restoreInput)
	outputDone := hs.BeginOutputStream(restoreInput)

	select {
	case err := <-outputDone:
		return err
	case <-inputDone:
		if hs.stdout != nil || hs.stderr != nil {
			select {
			case err := <-outputDone:
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		return nil
	case err := <-inputDetached:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (hs *HijackedStreamer) SetupInput() (func(), error) {
	if hs.in == nil || !hs.tty {
		return func() {}, nil
	}

	if err := hs.SetRawTerminal(); err != nil {
		return nil, fmt.Errorf("unable to set IO streams as raw terminal: %s", err)
	}

	var restoreOnce sync.Once
	restore := func() {
		restoreOnce.Do(func() {
			_ = hs.in.RestoreTerminal()
			_ = hs.out.RestoreTerminal()

			if hs.stdin != nil && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
				_ = hs.stdin.Close()
			}
		})
	}

	// The default escape key sequence: ctrl-p, ctrl-q
	proxyReader := term.NewEscapeProxy(hs.stdin, []byte{16, 17})
	hs.stdin = ioutils.NewReadCloserWrapper(proxyReader, hs.stdin.Close)

	return restore, nil
}

func (hs HijackedStreamer) SetRawTerminal() error {
	if err := hs.in.SetRawTerminal(); err != nil {
		return err
	}

	if err := hs.out.SetRawTerminal(); err != nil {
		return err
	}

	return nil
}

func (hs HijackedStreamer) BeginInputStream(restoreInput func()) (<-chan struct{}, <-chan error) {
	inputDone := make(chan struct{})
	inputDetached := make(chan error)

	go func() {
		if hs.stdin != nil {
			_, err := io.Copy(hs.resp.Conn, hs.stdin)

			// restore input as soon as possible
			restoreInput()

			if err != nil {
				if _, ok := err.(term.EscapeError); ok {
					inputDetached <- err
					return
				} else {
					// TODO: log error using logrus.debug
					// logrus.Debugf("Error sendStdin: %s", err)
				}
			}
		}

		if err := hs.resp.CloseWrite(); err != nil {
			// TODO: log error using logrus.debug
			// logrus.Debugf("Couldn't send EOF: %s", err)
		}

		close(inputDone)
	}()

	return inputDone, inputDetached
}

func (hs HijackedStreamer) BeginOutputStream(restoreInput func()) <-chan error {
	if hs.stdout == nil && hs.stderr == nil {
		return nil
	}

	outputDone := make(chan error)

	go func() {
		var err error

		if hs.stdout != nil && hs.tty {
			_, err = io.Copy(hs.stdout, hs.resp.Reader)
			// restore input as soon as possible
			restoreInput()
		} else {
			_, err = stdcopy.StdCopy(hs.stdout, hs.stderr, hs.resp.Reader)
		}

		if err != nil {
			// TODO: log error using logrus.debug
			// logrus.Debugf("Error receiveStdout: %s", err)
		}

		outputDone <- err
	}()

	return outputDone
}
