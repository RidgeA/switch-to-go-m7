package crawler

import "sync"

type (
	set struct {
		sync.RWMutex
		s    map[string]struct{}
	}
)

func newSet() *set {
	return &set{
		s:    make(map[string]struct{}),
	}
}

func (s *set) Add(v string) {
	s.Lock()
	defer s.Unlock()
	s.s[v] = struct{}{}
}

func (s *set) Has(v string) bool {
	s.RLock()
	defer s.RUnlock()
	_, exists := s.s[v]
	return exists
}

func (s *set) Del(v string) {
	s.Lock()
	defer s.Unlock()
	delete(s.s, v)
}

func (s *set) Empty() bool {
	s.Lock()
	defer s.Unlock()
	return len(s.s) == 0
}
