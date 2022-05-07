package logline

import "testing"

func TestLogLine_ToBytes(t *testing.T) {
	testCase(t, "log", "name", "log|name")
	testCase(t, "log", "", "log|")
	testCase(t, "", "name", "|name")
	testCase(t, "", "", "|")
}

func testCase(t *testing.T, log string, line string, expected string) {
	l := *New(log, line)
	if s := string(l.ToBytes()); s != expected {
		t.Error(expected, s)
	}
}
