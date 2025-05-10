package set

type Set[T comparable] map[T]struct{}

func New[T comparable](xs ...T) Set[T] {
	s := make(Set[T])
	for _, xs := range xs {
		s.Add(xs)
	}
	return s
}

func (s Set[T]) Add(x T) {
	s[x] = struct{}{}
}

func (s Set[T]) Includes(x T) bool {
	_, ok := s[x]
	return ok
}

func (s Set[T]) Remove(x T) {
	delete(s, x)
}

func (s Set[T]) ToSlice() []T {
	xs := make([]T, 0, len(s))
	for x := range s {
		xs = append(xs, x)
	}
	return xs
}

func (s Set[T]) Sub(other Set[T]) Set[T] {
	res := New[T]()
	for x := range s {
		if !other.Includes(x) {
			res.Add(x)
		}
	}
	return res
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	res := New[T]()
	for x := range s {
		if other.Includes(x) {
			res.Add(x)
		}
	}
	return res
}

func (s Set[T]) Equal(other Set[T]) bool {
	if len(s) != len(other) {
		return false
	}
	for x := range s {
		if !other.Includes(x) {
			return false
		}
	}
	return true
}

func Filter[T any](slice []T, condition func(T) bool) []T {
	res := make([]T, 0, len(slice))
	for _, t := range slice {
		if condition(t) {
			res = append(res, t)
		}
	}
	return res
}
