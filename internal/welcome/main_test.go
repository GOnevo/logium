package welcome

import (
	"github.com/gonevo/logium/internal/logline"
	"sync"
	"testing"
)

func BenchmarkWelcome_Append(b *testing.B) {
	w := New(5)
	l := &logline.LogLine{Log: "log", Line: "line"}

	for i := 0; i < b.N; i++ {
		w.Append("key", l)
	}

	if len(w.Get()["key"]) > 5 {
		b.Fatal("invalid size")
	}
}

func TestWelcome_Append(t *testing.T) {
	t.Parallel()
	w := New(5)

	l := logline.LogLine{Log: "log", Line: "line"}

	var wg sync.WaitGroup

	wg.Add(3)

	pushLine(w, "path1", &l, &wg)
	pushLine(w, "path1", &l, &wg)
	pushLine(w, "path2", &l, &wg)

	wg.Wait()

	if len(w.lines) != 2 {
		t.Error("incorrect size of map")
	}

	if len(w.lines["path1"]) != 2 {
		t.Error("incorrect size of path1 lines")
	}

	if len(w.lines["path2"]) != 1 {
		t.Error("incorrect size of path2 lines")
	}

	wg.Add(4)

	pushLine(w, "path1", &l, &wg)
	pushLine(w, "path1", &l, &wg)
	pushLine(w, "path1", &l, &wg)
	pushLine(w, "path1", &l, &wg)

	wg.Wait()

	if len(w.lines["path1"]) != 5 {
		t.Error("incorrect size of path1 lines")
	}
}

func pushLine(w *Welcome, key string, l *logline.LogLine, wg *sync.WaitGroup) {
	go func() {
		w.Append(key, l)
		wg.Done()
	}()
}
