package onexit

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dc0d/onexit/fnpq"
)

var (
	_pq       fnpq.PriorityQueue
	_pqmx     sync.Mutex
	_pqhandle sync.Once
	_pqdone   = make(chan struct{})
)

func init() {
	onSignal(final)
}

func final() {
	_pqhandle.Do(func() {
		defer close(_pqdone)
		_pqmx.Lock()
		defer _pqmx.Unlock()
		for _pq.Len() > 0 {
			item := fnpq.Pop(&_pq)
			item.Action()
		}
	})
}

// Done all registered deferred funcs are executed on signal.
// Call <-Done() at the end of main function (for example).
func Done() <-chan struct{} { return _pqdone }

// Register a function to be executed on app exit (by receiving an os signal),
// based on priority.
func Register(action func(), priority ...int) {
	_pqmx.Lock()
	defer _pqmx.Unlock()
	_priority := 0
	if len(priority) > 0 {
		_priority = priority[0]
	}
	t := fnpq.NewItem(action, _priority)
	fnpq.Push(&_pq, t)
}

func onSignal(f func(), sig ...os.Signal) {
	if f == nil {
		return
	}
	sigc := make(chan os.Signal, 1)
	if len(sig) > 0 {
		signal.Notify(sigc, sig...)
	} else {
		signal.Notify(sigc,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

// ForceExit ensures all registered deferred funcs are executed,
// and waits for completion of all of them, then calls os.Exit(code)
// with the provided code.
func ForceExit(code int) {
	forceExit(code, true)
}

func forceExit(code int, callExit bool) {
	final()
	<-Done()
	if callExit {
		os.Exit(code)
	}
}
