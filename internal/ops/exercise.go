package ops

import (
	"context"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
)

func CreateExercises(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	exercises ...types.Exercise,
) (opErr error) {
	for start, end := range batchIndexes(exercises, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		for i := start; i < end; i++ {
			iterEd := exercises[i]
			if iterEd.Name == "" {
				opErr = sberr.AppendError(
					types.InvalidExerciseErr, types.MissingExerciseNameErr,
				)
				return
			}
			if !types.ExerciseFocus(iterEd.FocusID).IsValid() {
				opErr = sberr.AppendError(
					types.InvalidExerciseErr,
					types.ErrInvalidExerciseFocus,
				)
				return
			}
			if !types.ExerciseKind(iterEd.KindID).IsValid() {
				opErr = sberr.AppendError(
					types.InvalidExerciseErr,
					types.ErrInvalidExerciseKind,
				)
				return
			}
		}

		var numRows int64
		chunk := exercises[start:end]
		_ = dal.BulkCreateExercisesParams(types.Exercise{})
		numRows, opErr = dal.Query1x2(
			dal.Q.BulkCreateExercises, queries, ctxt,
			*(*[]dal.BulkCreateExercisesParams)(unsafe.Pointer(&chunk)),
		)
		if opErr != nil {
			opErr = sberr.AppendError(types.CouldNotAddExercisesErr, opErr)
			return
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Added new exercises",
			"NumRows", numRows,
		)
	}

	return
}

func ReadNumExercises(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
) (res int64, opErr error) {
	res, opErr = dal.Query0x2(dal.Q.GetNumExercises, queries, ctxt)
	if opErr != nil {
		opErr = sberr.AppendError(types.CouldNotGetNumExercisesErr, opErr)
		return
	}
	state.Log.Log(ctxt, sblog.VLevel(3), "Read num exercises")
	return
}

func ReadExercisesByName(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	names ...string,
) (res []types.Exercise, opErr error) {
	res = make([]types.Exercise, len(names))

	for start, end := range batchIndexes(names, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		var rawData []dal.GetExercisesByNameRow
		rawData, opErr = dal.Query1x2(dal.Q.GetExercisesByName, queries, ctxt, names[start:end])
		if opErr != nil {
			opErr = sberr.AppendError(types.CouldNotFindRequestedExerciseErr, opErr)
			return
		}
		if len(rawData) != end-start {
			opErr = types.CouldNotFindRequestedExerciseErr
			return
		}

		_ = dal.BulkCreateExercisesParams(types.Exercise{})
		copy(
			res[start:end],
			*(*[]types.Exercise)(unsafe.Pointer(&rawData)),
		)

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Read clients from email",
			"Num", len(rawData),
		)
	}

	return
}

func UpdateExercises(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	exercises ...types.Exercise,
) (opErr error) {
	cntr := 0
	for _, e := range exercises {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		_ = dal.UpdateExerciseByNameParams(types.Exercise{})
		opErr = dal.Query1x1(
			dal.Q.UpdateExerciseByName, queries, ctxt,
			*(*dal.UpdateExerciseByNameParams)(unsafe.Pointer(&e)),
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotUpdateRequestedExerciseErr, opErr,
			)
			return
		}
		cntr++
	}
	if cntr != len(exercises) {
		opErr = sberr.AppendError(
			types.CouldNotUpdateRequestedExerciseErr,
			types.CouldNotFindRequestedExerciseErr,
		)
		return
	}
	state.Log.Log(ctxt, sblog.VLevel(3), "Updated clients", "Num", cntr)
	return
}

func DeleteExercises(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	names ...string,
) (opErr error) {
	// Deleting all referenced/referencing data is handled by cascade rules

	var count int64
	count, opErr = dal.Query1x2(dal.Q.DeleteExercisesByName, queries, ctxt, names)
	if opErr != nil {
		opErr = sberr.AppendError(types.CouldNotDeleteRequestedExerciseErr, opErr)
		return
	}
	if count != int64(len(names)) {
		opErr = sberr.AppendError(
			types.CouldNotDeleteRequestedExerciseErr,
			types.CouldNotFindRequestedExerciseErr,
		)
	}

	state.Log.Log(ctxt, sblog.VLevel(3), "Deleted exercises", "Num", count)
	return
}
