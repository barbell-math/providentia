package types

import "github.com/jackc/pgx/v5/pgtype"

type (
	RPE float32

	Kilogram float64

	Second float64

	Meter        float64
	MeterPerSec  float64
	MeterPerSec2 float64
	MeterPerSec3 float64

	Newton    float64
	NewtonSec float64

	Joule float64
	Watt  float64

	// V2[T ~float64]               = Vec2[T, T]
	Vec2[T ~float64, U ~float64] struct {
		X T
		Y U
	}

	Split struct {
		StartIdx int64
		EndIdx   int64
	}

	PointInTime[T ~float64, U ~float64] struct {
		Time  T
		Value U
	}
)

func (v *Vec2[T, U]) ScanPoint(newVal pgtype.Point) error {
	*v = Vec2[T, U]{X: T(newVal.P.X), Y: U(newVal.P.Y)}
	return nil
}

func (v Vec2[T, U]) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: float64(v.X), Y: float64(v.Y)},
		Valid: true,
	}, nil
}

func (v *Split) ScanPoint(newVal pgtype.Point) error {
	*v = Split{StartIdx: int64(newVal.P.X), EndIdx: int64(newVal.P.Y)}
	return nil
}

func (v Split) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: float64(v.StartIdx), Y: float64(v.EndIdx)},
		Valid: true,
	}, nil
}

func (v *PointInTime[T, U]) ScanPoint(newVal pgtype.Point) error {
	*v = PointInTime[T, U]{Time: T(newVal.P.X), Value: U(newVal.P.Y)}
	return nil
}

func (v PointInTime[T, U]) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: float64(v.Time), Y: float64(v.Value)},
		Valid: true,
	}, nil
}
