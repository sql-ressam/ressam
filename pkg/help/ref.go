package help

func Ref[T any](in T) *T {
	return &in
}
