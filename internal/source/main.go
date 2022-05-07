package source

import (
	"github.com/gonevo/logium/internal/utils"
	"github.com/gonevo/logium/internal/welcome"
	"github.com/nxadm/tail"
	"io"
)

// Source type
type Source struct {
	Name string
	Tail *tail.Tail
}

// Sources type
type Sources map[string]*Source

// Collection type
type Collection struct {
	sources Sources
}

// New returns new Collection of sources
func New() *Collection {
	return &Collection{sources: make(Sources)}
}

// Init reads last lines and adds tail to the collection
func (c *Collection) Init(path string, name string, w *welcome.Welcome) {
	cursor := utils.GetCursorForLastLines(path, uint64(w.Last()))
	t, err := tail.TailFile(path, tail.Config{
		Follow:   true,
		ReOpen:   true,
		Location: &tail.SeekInfo{Offset: cursor, Whence: io.SeekEnd},
	})
	if err != nil {
		return
	}

	c.add(path, &Source{Name: name, Tail: t})
}

// Get returns slice of sources
func (c *Collection) Get() Sources {
	return c.sources
}

// Stop stops all tails
func (c *Collection) Stop() {
	for _, source := range c.sources {
		_ = source.Tail.Stop()
	}
}

func (c *Collection) add(path string, s *Source) {
	c.sources[path] = s
}
