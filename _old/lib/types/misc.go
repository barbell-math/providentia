package types

import (
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	Found[T any] struct {
		Found bool
		Value T
	}

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

	CSVDataDirOptions struct {
		sbcsv.Opts
		*BarPathCalcHyperparams
		*BarPathTrackerHyperparams
		ClientCreateType      CreateFuncType
		ClientDir             string
		ExerciseCreateType    CreateFuncType
		ExerciseDir           string
		HyperparamsCreateType CreateFuncType
		HyperparamsDir        string
		WorkoutCreateType     CreateFuncType
		WorkoutDir            string
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
