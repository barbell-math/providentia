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

	Vec2[T ~float64] struct {
		X T
		Y T
	}
)

func (v *Vec2[T]) ScanPoint(newVal pgtype.Point) error {
	*v = Vec2[T]{X: T(newVal.P.X), Y: T(newVal.P.Y)}
	return nil
}

func (v Vec2[T]) PointValue() (pgtype.Point, error) {
	return pgtype.Point{
		P:     pgtype.Vec2{X: float64(v.X), Y: float64(v.Y)},
		Valid: true,
	}, nil
}
