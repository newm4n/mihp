package helper

import "testing"

func TestRandomName(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(RandomName())
	}
}
