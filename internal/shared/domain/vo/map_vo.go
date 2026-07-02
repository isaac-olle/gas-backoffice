package vo

type MapVo[K comparable, V any] struct {
	values map[K]V
}

func NewMapVo[K comparable, V any](initial map[K]V) *MapVo[K, V] {
	if initial == nil {
		initial = make(map[K]V)
	}
	return &MapVo[K, V]{values: initial}
}

func NewEmptyMapVo[K comparable, A any]() *MapVo[K, A] {
	return &MapVo[K, A]{values: make(map[K]A)}
}

func (rv *MapVo[K, V]) Set(key K, value V) {
	rv.values[key] = value
}

func (rv *MapVo[K, V]) Remove(key K) {
	delete(rv.values, key)
}

func (rv *MapVo[K, V]) Contains(key K) bool {
	_, ok := rv.values[key]
	return ok
}

func (rv *MapVo[K, V]) Get(key K) V {
	return rv.values[key]
}

func (rv *MapVo[K, V]) Values() map[K]V {
	return rv.values
}

func (rv *MapVo[K, V]) Size() int {
	return len(rv.values)
}
