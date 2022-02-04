package stack

type Stack struct {
	Items []int
}

func (s *Stack) Pop() int {
	if !s.IsEmpty() {
		f := s.Items[len(s.Items)-1]
		arr := s.Items[:len(s.Items)-1]
		s.Items = arr
		return f
	} else {
		return 0
	}
}

func (s *Stack) Push(n int) {
	arr := append(s.Items, n)
	s.Items = arr
}

func (s *Stack) Peek() int {
	if !s.IsEmpty() {
		return s.Items[len(s.Items)-1]
	} else {
		return 0
	}
}

func (s *Stack) IsEmpty() bool {
	return len(s.Items) == 0
}
