package main

import (
	"testing"
)

func TestParall_one(t *testing.T) {
	t.Parallel()
	for i := 0; i < 20; i++ {
		t.Logf("one %d\n", i)
	}
}

func TestParall_two(t *testing.T) {
	t.Parallel()
	for i := 0; i < 20; i++ {
		t.Logf("two %d\n", i)
	}
}
