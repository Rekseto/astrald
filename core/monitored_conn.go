package core

import (
	"github.com/cryptopunkscc/astrald/astral"
	"sync"
	"sync/atomic"
	"time"
)

const (
	StateRouting       = "routing"
	StateOpen          = "open"
	StateClosingTarget = "closing-t"
	StateClosingCaller = "closing-c"
	StateClosed        = "closed"
)

var nextConnID atomic.Int64

type MonitoredConn struct {
	id            int64
	target        *MonitoredWriter
	caller        *MonitoredWriter
	query         *astral.Query
	establishedAt time.Time

	closeMu      sync.Mutex
	targetClosed bool
	callerClosed bool
	done         chan struct{}
}

func NewMonitoredConn(caller *MonitoredWriter, target *MonitoredWriter, query *astral.Query) *MonitoredConn {
	conn := &MonitoredConn{
		id:            nextConnID.Add(1),
		query:         query,
		done:          make(chan struct{}),
		establishedAt: time.Now(),
	}

	conn.SetCaller(caller)
	conn.SetTarget(target)

	return conn
}

func (conn *MonitoredConn) SetTarget(target *MonitoredWriter) {
	conn.target = target
	if target != nil {
		target.AfterClose = func(err error) {
			conn.onTargetClosed()
		}
	}
}

func (conn *MonitoredConn) SetCaller(caller *MonitoredWriter) {
	conn.caller = caller
	if caller != nil {
		caller.AfterClose = func(err error) {
			conn.onCallerClosed()
		}
	}
}

func (conn *MonitoredConn) ID() int {
	return int(conn.id)
}

func (conn *MonitoredConn) Target() *MonitoredWriter {
	return conn.target
}

func (conn *MonitoredConn) Caller() *MonitoredWriter {
	return conn.caller
}

func (conn *MonitoredConn) Query() *astral.Query {
	return conn.query
}

func (conn *MonitoredConn) SetQuery(query *astral.Query) {
	conn.query = query
}

func (conn *MonitoredConn) BytesOut() int {
	if conn.target == nil {
		return 0
	}
	return conn.target.Bytes()
}

func (conn *MonitoredConn) BytesIn() int {
	if conn.caller == nil {
		return 0
	}
	return conn.caller.Bytes()
}

func (conn *MonitoredConn) Done() <-chan struct{} {
	return conn.done
}

func (conn *MonitoredConn) State() string {
	conn.closeMu.Lock()
	defer conn.closeMu.Unlock()

	if conn.target == nil {
		return StateRouting
	}
	var c int
	if conn.callerClosed {
		c++
	}
	if conn.targetClosed {
		c++
	}
	switch c {
	case 0:
		return StateOpen
	case 1:
		if conn.targetClosed {
			return StateClosingTarget
		}
		return StateClosingCaller
	case 2:
		return StateClosed
	}
	panic("?")
}

func (conn *MonitoredConn) onTargetClosed() {
	conn.closeMu.Lock()
	defer conn.closeMu.Unlock()

	if conn.targetClosed {
		return
	}

	conn.targetClosed = true
	if conn.callerClosed {
		close(conn.done)
	}
}

func (conn *MonitoredConn) onCallerClosed() {
	conn.closeMu.Lock()
	defer conn.closeMu.Unlock()

	if conn.callerClosed {
		return
	}

	conn.callerClosed = true
	if conn.targetClosed {
		close(conn.done)
	}
}
