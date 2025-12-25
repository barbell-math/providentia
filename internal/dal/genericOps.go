package dal

import (
	"context"
	"fmt"
	"strings"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	genericPoint struct {
		X float64
		Y float64
	}

	genericCreateOpts[T any] struct {
		Data        []T
		ValueGetter func(v *T, res *[]any) error
		TableName   string
		Columns     []string
		Err         error
	}

	genericCreateReturningIdVal[T any] struct {
		Val T
		Id  int64
	}
	genericCreateReturningIdOpts[T any] struct {
		Data                    []genericCreateReturningIdVal[T]
		ValueGetter             func(v *genericCreateReturningIdVal[T], res *[]any) error
		TableName               string
		Columns                 []string
		ModifyValuePlaceholders func(placeholders []string) []string
		Err                     error
	}

	genericReadTotalNumOpts struct {
		TableName string
		Res       *int64
	}

	genericReadByUniqueIdOpts[T any, U any] struct {
		Ids        []T
		Res        *[]U
		TableName  string
		Columns    []string
		UniqueCol  string
		IdsSqlType string
		Err        error
	}

	genericFindByUniqueIdOpts[T any, U types.Found[V], V any] struct {
		Ids           []T
		Res           *[]U
		TableName     string
		Columns       []string
		UniqueCol     string
		IdsSqlType    string
		SetScanValues func(v *V, res []any)
		Err           error
	}

	genericDeleteByUniqueIdOpts[T any] struct {
		Ids       []T
		TableName string
		UniqueCol string
		Err       error
	}
)

const (
	updateSerialIdSql = `
SELECT SETVAL(
	pg_get_serial_sequence('providentia.%s', 'id'),
	(SELECT MAX(id) FROM providentia.%s) + 1
);
`

	createReturningIdSql = `
INSERT INTO providentia.%s (%s) VALUES (%s) RETURNING id;
`

	ensureExistSql = `
INSERT INTO providentia.%s (%s) VALUES (%s) ON CONFLICT (%s) DO NOTHING;
`

	readTotalNumSql = `SELECT COUNT(*) FROM providentia.%s;`

	readByUniqueIdSql = `
SELECT %s FROM providentia.%s JOIN UNNEST($1::%s[])
WITH ORDINALITY t(%s, ord)
USING (%s) ORDER BY ord;
`

	findByUniqueIdSql = `
SELECT ord::INT8, %s FROM providentia.%s JOIN UNNEST($1::%s[])
WITH ORDINALITY t(%s, ord) USING (%s) ORDER BY ord;
`

	deleteByUniqueIdSql = `DELETE FROM providentia.%s WHERE %s = $1;`
)

func (v *genericPoint) ScanPoint(newVal pgtype.Point) error {
	*v = genericPoint{X: newVal.P.X, Y: newVal.P.Y}
	return nil
}

func (v genericPoint) PointValue() (pgtype.Point, error) {
	return pgtype.Point{P: pgtype.Vec2{X: v.X, Y: v.Y}, Valid: true}, nil
}

func genericCreateWithId[T any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericCreateOpts[T],
) error {
	if len(opts.Data) == 0 {
		return nil
	}

	genericCreate(ctxt, state, tx, opts)

	if _, err := tx.Exec(ctxt, fmt.Sprintf(
		updateSerialIdSql, opts.TableName, opts.TableName,
	)); err != nil {
		return sberr.AppendError(
			opts.Err, sberr.Wrap(err, "Failed to update serial index"),
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Updated serial count",
		"Table", opts.TableName,
	)
	return nil
}

func genericCreate[T any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericCreateOpts[T],
) error {
	if len(opts.Data) == 0 {
		return nil
	}

	cpy := CpyFromSlice[T]{Data: opts.Data, ValueGetter: opts.ValueGetter}
	if n, err := tx.CopyFrom(
		ctxt, pgx.Identifier{"providentia", opts.TableName},
		opts.Columns,
		&cpy,
	); err != nil {
		return sberr.AppendError(opts.Err, err)
	} else if n != int64(len(opts.Data)) {
		return sberr.Wrap(
			opts.Err,
			"Expected to create %d entires but only created %d, rolling back",
			len(opts.Data), n,
		)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		fmt.Sprintf("DAL: Created new %s entries", opts.TableName),
		"NumRows", len(opts.Data),
	)
	return nil
}

func genericCreateReturningId[T any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericCreateReturningIdOpts[T],
) error {
	if len(opts.Data) == 0 {
		return nil
	}

	commaSepCols := strings.Join(opts.Columns, ", ")
	sql := fmt.Sprintf(
		createReturningIdSql,
		opts.TableName, commaSepCols,
		strings.Join(opts.ModifyValuePlaceholders(
			defaultValuePlaceholders(len(opts.Columns)),
		), ", "),
	)
	cpy := CpyFromSlice[genericCreateReturningIdVal[T]]{
		Data: opts.Data, ValueGetter: opts.ValueGetter,
	}
	for start, end := range batchIndexes(opts.Data, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b := pgx.Batch{}
		for i := start; i < end; i++ {
			vals, err := cpy.Values()
			if err != nil {
				return sberr.AppendError(opts.Err, err)
			}
			b.Queue(sql, vals...)
			cpy.Next()
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if err := results.QueryRow().Scan(&opts.Data[start+i].Id); err != nil {
				results.Close()
				return sberr.AppendError(opts.Err, err)
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Created new %s entries returning id", opts.TableName),
			"NumRows", end-start,
		)
	}
	return nil
}

func genericEnsureExists[T any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericCreateOpts[T],
) error {
	if len(opts.Data) == 0 {
		return nil
	}

	commaSepCols := strings.Join(opts.Columns, ", ")
	sql := fmt.Sprintf(
		ensureExistSql,
		opts.TableName, commaSepCols,
		defaultValuePlaceholdersJoined(len(opts.Columns)), commaSepCols,
	)
	cpy := CpyFromSlice[T]{Data: opts.Data, ValueGetter: opts.ValueGetter}
	for start, end := range batchIndexes(opts.Data, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		created := int64(0)
		b := pgx.Batch{}
		for i := start; i < end; i++ {
			vals, err := cpy.Values()
			if err != nil {
				return sberr.AppendError(opts.Err, err)
			}
			b.Queue(sql, vals...)
			cpy.Next()
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return sberr.AppendError(opts.Err, err)
			} else {
				created += cmdTag.RowsAffected()
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Ensured %ss exist", opts.TableName),
			"NumAffectedRows/NumRows", fmt.Sprintf("%d/%d", created, end-start),
		)
	}
	return nil
}

func genericReadTotalNum(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericReadTotalNumOpts,
) error {
	row := tx.QueryRow(ctxt, fmt.Sprintf(readTotalNumSql, opts.TableName))
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		fmt.Sprintf("DAL: Read total num %ss", opts.TableName),
	)
	return row.Scan(opts.Res)
}

func genericReadByUniqueId[T any, U any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericReadByUniqueIdOpts[T, U],
) error {
	*opts.Res = util.SliceClamp(*opts.Res, len(opts.Ids))
	commaSepCols := strings.Join(opts.Columns, ", ")
	sql := fmt.Sprintf(
		readByUniqueIdSql, commaSepCols,
		opts.TableName, opts.IdsSqlType, opts.UniqueCol, opts.UniqueCol,
	)

	for start, end := range batchIndexes(opts.Ids, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		rows, err := tx.Query(ctxt, sql, opts.Ids[start:end])
		if err != nil {
			return sberr.AppendError(opts.Err, err)
		}

		cntr := start
		for rows.Next() {
			(*opts.Res)[cntr], err = pgx.RowToStructByName[U](rows)
			if err != nil {
				rows.Close()
				return sberr.AppendError(opts.Err, err)
			}
			cntr++
		}
		rows.Close()

		if cntr != end {
			return sberr.Wrap(
				opts.Err,
				"Only read %d entries out of batch of %d requests",
				cntr-start, end-start,
			)
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Read %ss by %s", opts.TableName, opts.UniqueCol),
			"NumRows", end-start,
		)
	}
	return nil
}

func genericFindByUniqueId[T any, U types.Found[V], V any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericFindByUniqueIdOpts[T, U, V],
) error {
	*opts.Res = util.SliceClamp(*opts.Res, len(opts.Ids))
	commaSepCols := strings.Join(opts.Columns, ", ")
	sql := fmt.Sprintf(
		findByUniqueIdSql, commaSepCols,
		opts.TableName, opts.IdsSqlType, opts.UniqueCol, opts.UniqueCol,
	)

	for start, end := range batchIndexes(opts.Ids, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		rows, err := tx.Query(ctxt, sql, opts.Ids[start:end])
		if err != nil {
			return sberr.AppendError(opts.Err, err)
		}

		var iterVal V
		ord := int64(0)
		found := int64(0)
		scanValues := make([]any, len(opts.Columns)+1)
		scanValues[0] = &ord
		opts.SetScanValues(&iterVal, scanValues[1:])

		for rows.Next() {
			if err := rows.Scan(scanValues...); err != nil {
				rows.Close()
				return sberr.AppendError(opts.Err, err)
			}
			(*opts.Res)[int64(start)+ord-1] = U{
				Value: iterVal,
				Found: true,
			}
			found++
		}
		rows.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Found %ss by %s", opts.TableName, opts.UniqueCol),
			"NumFound/NumRows", fmt.Sprintf("%d/%d", found, end-start),
		)
	}
	return nil
}

func genericDeleteByUniqueId[T any](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *genericDeleteByUniqueIdOpts[T],
) error {
	if len(opts.Ids) == 0 {
		return nil
	}

	// Deleting all referenced/referencing data is handled by cascade rules

	sql := fmt.Sprintf(deleteByUniqueIdSql, opts.TableName, opts.UniqueCol)

	for start, end := range batchIndexes(opts.Ids, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b := pgx.Batch{}
		for i := start; i < end; i++ {
			b.Queue(sql, opts.Ids[i])
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return sberr.AppendError(opts.Err, err)
			} else if cmdTag.RowsAffected() == 0 {
				results.Close()
				return sberr.Wrap(
					opts.Err,
					"Could not delete entry with id '%v' (Does id exist?)",
					opts.Ids[i],
				)
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Deleted %ss", opts.TableName),
			"NumRows", end-start,
		)
	}
	return nil
}
