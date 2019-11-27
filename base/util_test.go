package base

import (
	"testing"
)

func TestMkfifo(t *testing.T) {
	t.Log(Mkfiof("./a"))
}
