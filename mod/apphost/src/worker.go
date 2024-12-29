package apphost

import (
	"context"
	"errors"
	"io"
	"strings"
)

func (mod *Module) worker(ctx context.Context) error {
	for conn := range mod.conns {
		var done = make(chan struct{})

		go func() {
			select {
			case <-ctx.Done():
				conn.Close()
			case <-done:
			}
		}()

		session := NewSession(mod, conn, mod.log)
		err := session.Serve(ctx)

		switch {
		case err == nil:
		case errors.Is(err, io.EOF):
		case strings.Contains(err.Error(), "connection closed"):
		case strings.Contains(err.Error(), "use of closed network connection"):
		case strings.Contains(err.Error(), "read/write on closed pipe"):
		default:
			mod.log.Error("serve error: %v", err)
		}

		conn.Close()
		close(done)
	}
	return nil
}
