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

	Optional[T any] struct {
		Present bool
		Value   T
	}

	Vec2[T ~float64, U ~float64] struct {
		X T
		Y U
	}

	PointInTime[T ~float64, U ~float64] struct {
		Time  T
		Value U
	}

	Split struct {
		StartIdx int64
		EndIdx   int64
	}
)
