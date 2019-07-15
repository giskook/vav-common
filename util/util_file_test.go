package util

import (
	"testing"
)

func TestMkfifo(t *testing.T) {
	Mkfifo("./current/pipe_example")
}
