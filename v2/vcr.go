package traindown

type node struct {
	r    rune
	next *node
}

type vcr struct {
	start *node
}

func newVCR() vcr {
	return vcr{}
}

func (s *vcr) push(r rune) {
	n := &node{r: r}

	if s.start == nil {
		s.start = n
	} else {
		n.next = s.start
		s.start = n
	}
}

func (s *vcr) pop() rune {
	if s.start == nil {
		return EOFRune
	}

	n := s.start
	s.start = n.next

	return n.r
}

func (s *vcr) clear() {
	s.start = nil
}
