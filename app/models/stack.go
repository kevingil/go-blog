package models

type stack []interface{}

func (s *stack) Push(v interface{}) {
	*s = append(*s, v)
}

func (s *stack) Pop() interface{} {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *stack) Len() int {
	return len(*s)
}

func (s *stack) IsEmpty() bool {
	return len(*s) == 0
}
