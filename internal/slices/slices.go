package slices

func Map[S ~[]E1, E1, E2 any](s S, f func(v E1) E2) []E2 {
	var r []E2
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}
