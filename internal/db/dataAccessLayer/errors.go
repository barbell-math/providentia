package dal

import (
	"errors"

	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5/pgconn"
)

func FormatErr(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Where != "" {
		err = sberr.Wrap(
			err,
			"Detail: %s\nWhere: %s",
			pgErr.Detail, pgErr.Where,
		)
	} else {
		err = sberr.Wrap(err, "Detail: %s", pgErr.Detail)
	}
	return err
}
