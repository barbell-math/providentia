package migrations

import (
	"context"
	"embed"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbsql"
	"github.com/jackc/pgx/v5"
)

//go:embed *.sql
var SqlMigrations embed.FS
var PostOps = map[sbsql.Migration]sbsql.PostMigrationOp[*types.State]{
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
		if err := dal.CreateModelsWithID(
			ctxt, state, tx, ModelSetupData,
		); err != nil {
			return err
		}
		if err := dal.CreateHyperparams(
			ctxt, state, tx, BarPathTrackerHyperparamsSetupData,
		); err != nil {
			return err
		}
		if err := dal.CreateHyperparams(
			ctxt, state, tx, BarPathCalcHyperparamsSetupData,
		); err != nil {
			return err
		}

		return nil
	},
}

func RunMigrations(ctxt context.Context, state *types.State) (opErr error) {
	m := sbsql.Migrations[*types.State]{}
	if opErr = m.Load(SqlMigrations, ".", PostOps); opErr != nil {
		return
	}
	if opErr = m.Run(ctxt, state.DB, state); opErr != nil {
		return
	}

	return
}
