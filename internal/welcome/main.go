package welcome

import (
	"github.com/gonevo/logium/internal/logline"
	"sync"
)

type LogLineSlice []*logline.LogLine

// LogLines type
type LogLines map[string]LogLineSlice

// Welcome type
type Welcome struct {
	lines LogLines
	limit int
	sync.Mutex
}

// New returns new Welcome
func New(limit int) *Welcome {
	return &Welcome{limit: limit, lines: make(LogLines)}
}

// Get returns slice of LogLines
func (w *Welcome) Get() LogLines {
	return w.lines
}

// Last returns limit
func (w *Welcome) Last() int {
	return w.limit
}

// Append adds new LogLine to the map
func (w *Welcome) Append(path string, l *logline.LogLine) {
	w.Lock()
	defer w.Unlock()

	sl, ok := w.lines[path]
	if !ok {
		sl = make([]*logline.LogLine, 0, w.Last())
	}

	if len(sl)+1 > w.limit {
		sl = sl[1:]
	}
	w.lines[path] = append(sl, l)
}
