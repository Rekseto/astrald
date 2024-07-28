package astral

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// ErrRejected - the query was rejected by the target
var ErrRejected = errors.New("query rejected")

// ErrAborted - query was aborted and routing did not finish
var ErrAborted = errors.New("query aborted")

// ErrTimeout - query timed out
var ErrTimeout = errors.New("query timeout")

// ErrZoneExcluded - operation requires zones excluded from the scope
var ErrZoneExcluded = errors.New("zone excluded")

// ErrRouteNotFound - failed to route the query to the destination
type ErrRouteNotFound struct {
	Router Router
	Fails  []error
}

func RouteNotFound(r Router, errors ...error) (io.WriteCloser, error) {
	return nil, &ErrRouteNotFound{
		Router: r,
		Fails:  errors,
	}
}

func (e *ErrRouteNotFound) Error() string {
	return "route not found"
}

func (e *ErrRouteNotFound) Trace() string {
	var buf = &bytes.Buffer{}

	var p func(indent int, e *ErrRouteNotFound)

	p = func(indent int, e *ErrRouteNotFound) {
		in := strings.Repeat("  ", indent)
		fmt.Fprintf(buf, "%s%s: %s (%d suberrors)\n", in, reflect.TypeOf(e.Router), e.Error(), len(e.Fails))
		for _, sub := range e.Fails {
			if rnf, ok := sub.(*ErrRouteNotFound); ok {
				p(indent+1, rnf)
			} else {
				fmt.Fprintf(buf, "%s  %s\n", in, sub.Error())
			}
		}
	}

	p(0, e)

	return string(buf.Bytes())
}

func (e *ErrRouteNotFound) Is(other error) bool {
	_, ok := other.(*ErrRouteNotFound)
	return ok
}
