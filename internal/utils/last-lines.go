package utils

import (
	"io"
	"os"
)

// GetCursorForLastLines returns cursor for last lines position
func GetCursorForLastLines(filepath string, n uint64) int64 {
	fileHandle, err := os.Open(filepath)
	if err != nil {
		return 0
	}
	defer func(fileHandle *os.File) {
		_ = fileHandle.Close()
	}(fileHandle)

	var lines = 0
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {

		if filesize == 0 {
			break
		}

		cursor--
		_, _ = fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		_, _ = fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // line found
			lines++

			if lines == int(n) {
				cursor++
				break
			}
		}

		if cursor == -filesize {
			if lines == 0 {
				cursor = 0
			}
			break
		}
	}

	return cursor
}
