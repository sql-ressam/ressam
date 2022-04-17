package help

// Ref returns a reference to T in.
func Ref[T any](in T) *T {
	return &in
}
