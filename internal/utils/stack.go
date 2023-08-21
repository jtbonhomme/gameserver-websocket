package utils

import (
	"fmt"
	"os"
)

type Stack[T any] []T

// IsEmpty checks if stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Len returns stack length
func (s *Stack[T]) Len() int {
	return len(*s)
}

// push pushes an array of new values onto the stack.
func (s *Stack[T]) push(t []T) {
	*s = append(*s, t...) // Simply append the new value to the end of the stack
}

// Push pushes a new value onto the stack.
// Example:
//
//	stack := utils.Stack[int]{}
//	stack.Push(1)
//	v := []int{2, 3, 4, 5, 6}
//	stack.Push(v...)
//
// --> &utils.Stack[int]{1, 2, 3, 4, 5, 6}
func (s *Stack[T]) Push(t ...T) {
	s.push(t)
}

// Pop removes and return top element of stack. Return false if stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	var element T
	if s.IsEmpty() {
		return element, false
	} else {
		index := len(*s) - 1  // Get the index of the top most element.
		element = (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]     // Remove it from the stack by slicing it off.
		return element, true
	}
}

// Top returns top element of stack. Return false if stack is empty.
func (s *Stack[T]) Top() (T, bool) {
	var element T
	if s.IsEmpty() {
		return element, false
	} else {
		index := len(*s) - 1  // Get the index of the top most element.
		element = (*s)[index] // Index into the slice and obtain the element.
		return element, true
	}
}

// Dump displays out the content of stack on stderr.
// Use it for debug only.
func (s *Stack[T]) Dump() {
	fmt.Fprintf(os.Stderr, "%#v\n", s)
}
