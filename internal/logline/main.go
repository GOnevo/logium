package logline

import (
	"fmt"
)

// Sep is separator symbol
var Sep = "|"

// LogLine type
type LogLine struct {
	Log  string
	Line string
}

// New returns LogLine instance
func New(name string, line string) *LogLine {
	return &LogLine{Log: name, Line: line}
}

// ToBytes converts structure to the byte slice
func (l *LogLine) ToBytes() []byte {
	return []byte(fmt.Sprintf("%s%s%s", l.Log, Sep, l.Line))
}
