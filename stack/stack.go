package stack

type Stack[T any] struct {
	data  []T
	elems int
}

// Add to the stack
func (s *Stack[T]) Push(v T) {
	s.data = append(s.data, v)
	s.elems += 1
}

// Remove and return the top of the stack
// Second return value is false if this is called on an empty stack
func (s *Stack[T]) Pop() (T, bool) {
	var v T
	l := len(s.data)
	if l == 0 {
		return v, false
	}

	// Get the last value
	v = s.data[l-1]
	// Remove the last value
	s.data = s.data[:l-1]

	s.elems -= 1

	return v, true
}

// Return the top of the stack without removing it
// Second value returns false if stack is empty
func (s *Stack[T]) Peek() (T, bool) {
	var v T
	l := len(s.data)

	if l == 0 {
		return v, false
	}

	return s.data[l-1], true
}

func (s *Stack[T]) Get(i int) (T, bool) {
	var v T
	l := len(s.data)
	if l == 0 {
		return v, false
	}

	return s.data[i], true
}

// Get the amount of elements in the stack
func (s *Stack[T]) Size() int {
	return s.elems
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}
