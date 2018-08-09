package main

type Set struct {
	m map[string]bool
}

func NewSet() *Set {
	s := &Set{}
	s.m = make(map[string]bool)
	return s
}

func (s *Set) Add(value string) {
	// Never add something that has been removed
	
	_, e := s.m[value]
	if(!e) {
		s.m[value] = true
	}
}

func (s *Set) Remove(value string) {
	// Tombstone it forever
	
	s.m[value] = false
	//delete(s.m, value)
}

func (s *Set) Contains(value string) bool {
	x := s.m[value]
	return x
}

func (s *Set) CleanUp() {
	for k, v := range s.m {
		if(!v) {
			delete(s.m, k)
		}
	}
}

func (s *Set) SetToSlice() []string {
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}
