package utils

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack()

	stack.Push(1)
	stack.Push(2)

	if stack.Top() != 2 {
		t.Error("stack.Top() expected 2!")
	}

	if stack.Pop() != 2 {
		t.Error("stack.Pop() expected 2!")
	}
	if stack.Empty() {
		t.Error("stack.Empty() expected false!")
	}

	stack.Pop()
	if !stack.Empty() {
		t.Error("stack.Empty() expected true!")
	}
}
