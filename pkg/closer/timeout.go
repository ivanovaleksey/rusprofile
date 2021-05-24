package closer

import (
	"errors"
	"io"
	"time"
)

var ErrTimeout = errors.New("timeout")

type TimeoutCloser struct {
	timeout time.Duration
	closer  io.Closer
}

func NewTimeoutCloser(closer io.Closer, to time.Duration) *TimeoutCloser {
	return &TimeoutCloser{
		timeout: to,
		closer:  closer,
	}
}

func (t *TimeoutCloser) Close() error {
	done := make(chan error, 1)

	go func() {
		done <- t.closer.Close()
	}()

	select {
	case err := <-done:
		if err != nil {
			return err
		}
		return nil
	case <-time.After(t.timeout):
		return ErrTimeout
	}
}
