package vo

type SetVo[K comparable] struct {
	values map[K]struct{}
}

func NewSetVo[K comparable](elements ...K) *SetVo[K] {
	set := &SetVo[K]{values: make(map[K]struct{})}
	for _, e := range elements {
		set.values[e] = struct{}{}
	}
	return set
}

func NewEmptySetVo[K comparable]() *SetVo[K] {
	return &SetVo[K]{values: make(map[K]struct{})}
}

func (s *SetVo[K]) Add(k K) {
	s.values[k] = struct{}{}
}

func (s *SetVo[K]) Remove(k K) {
	delete(s.values, k)
}

func (s *SetVo[K]) Contains(k K) bool {
	_, ok := s.values[k]
	return ok
}

func (s *SetVo[K]) IsEmpty() bool {
	return len(s.values) == 0
}

func (s *SetVo[K]) Values() []K {
	keys := make([]K, 0, len(s.values))
	for k := range s.values {
		keys = append(keys, k)
	}
	return keys
}

func (s *SetVo[K]) Size() int {
	return len(s.values)
}
