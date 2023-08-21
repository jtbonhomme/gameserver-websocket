package utils

type Stack[T any] []T

// IsEmpty: check if stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

// Len returns stack length
func (s *Stack[T]) Len() int {
	return len(*s)
}

// Push a new value onto the stack
func (s *Stack[T]) Push(t T) {
	*s = append(*s, t) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
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

// Return top element of stack. Return false if stack is empty.
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
