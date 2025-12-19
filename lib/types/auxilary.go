package types

type (
	IdWrapper[T any] struct {
		Id  int64 `db:"id"`
		Val T
	}

	Found[T any] struct {
		Found bool
		Value T
	}
)
