package migrations

import (
	"context"
	"embed"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbsqlm "code.barbellmath.net/barbell-math/smoothbrain-sqlmigrate"
	"github.com/jackc/pgx/v5"
)

//go:embed *.sql
var SqlMigrations embed.FS
var PostOps = map[sbsqlm.Migration]sbsqlm.PostMigrationOp[*types.State]{
	0: func(ctxt context.Context, tx pgx.Tx, state *types.State) error {
		if err := dal.CreateExerciseFocusWithID(
			ctxt, state, tx, ExerciseFocusSetupData,
		); err != nil {
			return err
		}
		if err := dal.CreateExerciseKindWithID(
			ctxt, state, tx, ExerciseKindSetupData,
		); err != nil {
			return err
		}
		if err := dal.CreateExercisesWithID(
			ctxt, state, tx, ExerciseSetupData,
		); err != nil {
			return err
		}

		// TODO - update to be like the rest!
		// cnt, err = q.BulkCreateModelsWithID(ctxt, ModelSetupData)
		// if err != nil {
		// 	return err
		// } else if cnt != int64(len(ModelSetupData)) {
		// 	return sberr.Wrap(
		// 		SetupMissingDataErr,
		// 		"Setting up the providentia.model table: Expected %d rows to be created, got %d",
		// 		len(ModelSetupData), cnt,
		// 	)
		// }
		// if err := q.UpdateModelSerialCount(ctxt); err != nil {
		// 	return sberr.Wrap(
		// 		SerialUpdateCmdErr,
		// 		"Setting up the providentia.model table",
		// 	)
		// }

		// cnt, err = q.BulkCreateHyperparams(ctxt, HyperparamsSetupData)
		// if err != nil {
		// 	return err
		// } else if cnt != int64(len(HyperparamsSetupData)) {
		// 	return sberr.Wrap(
		// 		SetupMissingDataErr,
		// 		"Setting up the providentia.hyperparams table: Expected %d rows to be created, got %d",
		// 		len(HyperparamsSetupData), cnt,
		// 	)
		// }
		// if err := q.UpdateHyperparamsSerialCount(ctxt); err != nil {
		// 	return sberr.Wrap(
		// 		SerialUpdateCmdErr,
		// 		"Setting up the providentia.hyperparams table",
		// 	)
		// }

		return nil
	},
}

func RunMigrations(ctxt context.Context, state *types.State) (opErr error) {
	m := sbsqlm.Migrations[*types.State]{}
	if opErr = m.Load(SqlMigrations, ".", PostOps); opErr != nil {
		return
	}
	if opErr = m.Run(ctxt, state.DB, state); opErr != nil {
		return
	}

	return
}
