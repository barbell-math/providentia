package dal

import (
	"time"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5/pgtype"
)

func TimeToPGDate(t time.Time) pgtype.Date {
	return pgtype.Date{
		Time:             t,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func Vec2ToPoint[T ~float64, U ~float64](v types.Vec2[T, U]) pgtype.Point {
	return pgtype.Point{
		P:     pgtype.Vec2{X: float64(v.X), Y: float64(v.Y)},
		Valid: true,
	}
}
