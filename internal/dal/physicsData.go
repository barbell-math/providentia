package dal

import (
	"context"
	"fmt"
	"unsafe"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

const (
	physicsDataTableName = "physics_data"

	barPathCalcIdSelectSql = `(
	SELECT providentia.model.id FROM providentia.hyperparams
	JOIN providentia.model
		ON providentia.model.id = providentia.hyperparams.model_id
	WHERE providentia.model.name='%s'
		AND providentia.hyperparams.version=$2
)`

	barPathTrackerIdSelectSql = `(
	SELECT providentia.model.id FROM providentia.hyperparams
	JOIN providentia.model
		ON providentia.model.id = providentia.hyperparams.model_id
	WHERE providentia.model.name='%s'
		AND providentia.hyperparams.version=$3
)`
)

// TODO - make sure this work through workout tests in the tests dir
func createPhysicsDataReturningIds(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []genericCreateReturningIdVal[*types.PhysicsData],
) error {
	return genericCreateReturningId(
		ctxt, state, tx, &genericCreateReturningIdOpts[*types.PhysicsData]{
			TableName: physicsDataTableName,
			Columns: []string{
				"path", "bar_path_calc_id", "bar_path_track_id",
				"time",
				"position", "velocity", "acceleration", "jerk",
				"force", "impulse", "work", "power",
				"rep_splits",
				"min_vel", "max_vel",
				"min_acc", "max_acc",
				"min_force", "max_force",
				"min_impulse", "max_impulse",
				"avg_work", "min_work", "max_work",
				"avg_power", "min_power", "max_power",
			},
			ValueGetter: func(
				v *genericCreateReturningIdVal[*types.PhysicsData],
				res *[]any,
			) error {
				*res = make([]any, 27)
				(*res)[0] = v.Val.VideoPath
				(*res)[1] = v.Val.BarPathCalcVersion
				(*res)[2] = v.Val.BarPathTrackerVersion
				(*res)[3] = v.Val.Time
				(*res)[4] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.Position))
				(*res)[5] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.Velocity))
				(*res)[6] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.Acceleration))
				(*res)[7] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.Jerk))
				(*res)[8] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.Force))
				(*res)[9] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.Impulse))
				(*res)[10] = v.Val.Work
				(*res)[11] = v.Val.Power
				(*res)[12] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.RepSplits))
				(*res)[13] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MinVel))
				(*res)[14] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MaxVel))
				(*res)[15] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MinAcc))
				(*res)[16] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MaxAcc))
				(*res)[17] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MinForce))
				(*res)[18] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MaxForce))
				(*res)[19] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MinImpulse))
				(*res)[20] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MaxImpulse))
				(*res)[21] = v.Val.AvgWork
				(*res)[22] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MinWork))
				(*res)[23] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MaxWork))
				(*res)[24] = v.Val.AvgPower
				(*res)[25] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MinPower))
				(*res)[26] = *(*[]genericPoint)(unsafe.Pointer(&v.Val.MaxPower))
				return nil
			},
			ModifyValuePlaceholders: func(placeholders []string) []string {
				placeholders[0] = "$1::TEXT"
				placeholders[1] = fmt.Sprintf(
					barPathCalcIdSelectSql, types.BarPathCalc,
				)
				placeholders[2] = fmt.Sprintf(
					barPathTrackerIdSelectSql, types.BarPathTracker,
				)
				return placeholders
			},
			Data: data,
			Err:  types.CouldNotCreateAllPhysicsDataErr,
		},
	)
}
