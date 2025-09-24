package ops

import (
	"context"
	"unsafe"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
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

		chunk := exercises[start:end]
		if opErr = validateExercises(chunk); opErr != nil {
			return
		}

		var numRows int64
		_ = dal.BulkCreateExercisesParams(types.Exercise{})
		numRows, opErr = dal.Query1x2(
			dal.Q.BulkCreateExercises, queries, ctxt,
			*(*[]dal.BulkCreateExercisesParams)(unsafe.Pointer(&chunk)),
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotAddExercisesErr, dal.FormatErr(opErr),
			)
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

func EnsureExercisesExist(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	exercises ...types.Exercise,
) (opErr error) {
	names := make([]string, min(len(exercises), int(state.Global.BatchSize)))
	// TODO - would be nice to have the types be the enum types but sqlc
	// continues to be a gigantic pain in my ass
	kindIDs := make([]int32, min(len(exercises), int(state.Global.BatchSize)))
	focusIDs := make([]int32, min(len(exercises), int(state.Global.BatchSize)))

	for start, end := range batchIndexes(exercises, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		chunk := exercises[start:end]
		if opErr = validateExercises(chunk); opErr != nil {
			return
		}

		for i, e := range chunk {
			names[i] = e.Name
			kindIDs[i] = int32(e.KindID)
			focusIDs[i] = int32(e.FocusID)
		}

		opErr = dal.Query1x1(
			dal.Q.EnsureExercisesExist, queries, ctxt,
			dal.EnsureExercisesExistParams{
				Names:   names[:len(chunk)],
				Kinds:   kindIDs[:len(chunk)],
				Focuses: focusIDs[:len(chunk)],
			},
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotAddExercisesErr, dal.FormatErr(opErr),
			)
			return
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Ensured exercises exist",
			"NumExercises", len(chunk),
		)
	}

	return
}

func validateExercises(exercises []types.Exercise) (opErr error) {
	for _, iterEd := range exercises {
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
	return
}

func CreateExercisesFromCSV(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	exercises := []types.Exercise{}
	opts.ReuseRecord = true

	for _, file := range files {
		if opErr = sbcsv.LoadCSVFile(file, &sbcsv.LoadOpts{
			Opts:          opts,
			RequestedCols: sbcsv.ReqColsForStruct[types.Exercise](),
			Op:            sbcsv.RowToStructOp(&exercises),
		}); opErr != nil {
			return opErr
		}
	}

	return CreateExercises(ctxt, state, queries, exercises...)
}

func ReadNumExercises(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
) (res int64, opErr error) {
	res, opErr = dal.Query0x2(dal.Q.GetNumExercises, queries, ctxt)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotGetNumExercisesErr, dal.FormatErr(opErr),
		)
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
		rawData, opErr = dal.Query1x2(
			dal.Q.GetExercisesByName, queries, ctxt, names[start:end],
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotFindRequestedExerciseErr, dal.FormatErr(opErr),
			)
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
			"Read exercises from name",
			"Num", len(rawData),
		)
	}

	return
}

func FindExercisesByName(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	names ...string,
) (res []types.Found[types.Exercise], opErr error) {
	res = make([]types.Found[types.Exercise], len(names))

	for start, end := range batchIndexes(names, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return
		default:
		}

		var rawData []dal.FindExercisesByNameRow
		rawData, opErr = dal.Query1x2(
			dal.Q.FindExercisesByName, queries, ctxt, names[start:end],
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotFindRequestedExerciseErr, dal.FormatErr(opErr),
			)
			return
		}

		rawDataIdx := 0
		for i := 0; i < end-start; i++ {
			res[i+start].Found = (rawDataIdx < len(rawData) && rawData[rawDataIdx].Ord-1 == int64(i))
			if res[i+start].Found {
				res[i+start].Value = types.Exercise{
					Name:    rawData[rawDataIdx].Name,
					KindID:  rawData[rawDataIdx].KindID,
					FocusID: rawData[rawDataIdx].FocusID,
				}
				rawDataIdx++
			}
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Found exercises from name",
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
				types.CouldNotUpdateRequestedExerciseErr, dal.FormatErr(opErr),
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
		opErr = sberr.AppendError(
			types.CouldNotDeleteRequestedExerciseErr, dal.FormatErr(opErr),
		)
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
