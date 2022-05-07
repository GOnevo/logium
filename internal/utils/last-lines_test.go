package utils

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetLastLines(t *testing.T) {
	defer func() {
		_ = os.Remove("test.file")
	}()
	t.Run("3 lines", func(t *testing.T) {
		WriteTestContent("a\nb\nc\nd\ne\n")

		cursor := GetCursorForLastLines("test.file", 3)
		if lines := ReadFrom(cursor); !reflect.DeepEqual(lines, []string{"c", "d", "e", ""}) {
			t.Error(len(lines), lines)
		}
	})
	t.Run("empty last line and 2", func(t *testing.T) {
		WriteTestContent("1\n2\n\n")

		cursor := GetCursorForLastLines("test.file", 3)
		if lines := ReadFrom(cursor); !reflect.DeepEqual(lines, []string{"1", "2", "", ""}) {
			t.Error(len(lines), lines)
		}
	})
	t.Run("empty last line", func(t *testing.T) {
		WriteTestContent("\n")

		cursor := GetCursorForLastLines("test.file", 3)
		if lines := ReadFrom(cursor); !reflect.DeepEqual(lines, []string{""}) {
			t.Error(len(lines), lines)
		}
	})
}

func WriteTestContent(s string) {
	_ = ioutil.WriteFile("test.file", []byte(s), 0644)
}

func ReadFrom(cursor int64) []string {
	file, err := ioutil.ReadFile("test.file")
	if err != nil {
		return nil
	}

	content := string(file)

	return strings.Split(content[len(content)+int(cursor):], "\n")
}
