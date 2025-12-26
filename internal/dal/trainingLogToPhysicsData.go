package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

type (
	trainingLogToPhysicsData struct {
		TrainingLogId *int64 `db:"training_log_id"`
		PhysicsId     *int64 `db:"physics_id"`
		SetNum        int32  `db:"set_num"`
	}
)

const (
	trainingLogToPhysicsDataTableName = "training_log_to_physics_data"
)

func createTrainingLogToPhysicsMappings(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []trainingLogToPhysicsData,
) error {
	return genericCreate(
		ctxt, state, tx, &genericCreateOpts[trainingLogToPhysicsData]{
			TableName: trainingLogToPhysicsDataTableName,
			Columns:   []string{"training_log_id", "physics_id", "set_num"},
			ValueGetter: func(v *trainingLogToPhysicsData, res *[]any) error {
				*res = util.SliceClamp(*res, 3)
				(*res)[0] = v.TrainingLogId
				(*res)[1] = v.PhysicsId
				(*res)[2] = v.SetNum
				return nil
			},
			Data: data,
			Err:  types.CouldNotCreateAllTrainingLogPhysicsDataMappingsErr,
		},
	)
}
