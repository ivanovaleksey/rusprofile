package closer

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type CloseFunc func() error

type Closer struct {
	funcs []CloseFunc
	mu    sync.Mutex

	once   sync.Once
	closed chan struct{}
}

func New(signals ...os.Signal) *Closer {
	cl := &Closer{
		closed: make(chan struct{}),
	}
	if len(signals) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, signals...)

			sign := <-ch
			fmt.Printf("got signal: %s\n", sign)

			signal.Stop(ch)
			cl.CloseAll()
		}()
	}
	return cl
}

func (c *Closer) Add(closer CloseFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, closer)
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.closed)

		c.mu.Lock()
		defer c.mu.Unlock()

		errs := make(chan error, len(c.funcs))
		for _, closeFunc := range c.funcs {
			go func(fn CloseFunc) {
				errs <- fn()
			}(closeFunc)
		}
		for i := 0; i < cap(errs); i++ {
			err := <-errs
			if err != nil {
				fmt.Printf("closer error: %v\n", err)
			}
		}
	})
}

func (c *Closer) Wait() {
	<-c.closed
}
