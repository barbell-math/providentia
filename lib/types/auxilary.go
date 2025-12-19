package types

type (
	Found[T any] struct {
		Found bool
		Value T
	}
)
