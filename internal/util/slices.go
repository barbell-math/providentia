package util

import (
	"errors"
	"iter"
	"os"
	"path"
	"time"

	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"golang.org/x/exp/constraints"
)

type (
	FilesWithExtInDirOpts struct {
		DirsAllowed           bool
		OtherFileTypesAllowed bool
	}
)

var (
	FilesWithExtInDirErr = errors.New("Could not get files in dir")
)

func SliceClamp[S []T, T any, U constraints.Integer](s S, _len U) S {
	if len(s) < int(_len) {
		s = make([]T, _len)
	}
	return s[:_len]
}

func SliceSeq2Err[S []T, T any](s S) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for _, v := range s {
			if !yield(v, nil) {
				return
			}
		}
	}
}

func FilterSeq2Err[T any](
	s iter.Seq2[T, error],
	filter func(v *T, e *error) bool,
) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for v, e := range s {
			if (filter(&v, &e) || e != nil) && !yield(v, e) {
				return
			}
		}
	}
}

func DateEqual(date1 time.Time, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func FilesWithExtInDir(
	dir string,
	ext string,
	opts FilesWithExtInDirOpts,
) iter.Seq2[string, error] {
	var err error
	var dirEntries []os.DirEntry
	if dir != "" {
		dirEntries, err = os.ReadDir(dir)
	}

	return func(yield func(string, error) bool) {
		if err != nil {
			yield("", sberr.AppendError(FilesWithExtInDirErr, err))
			return
		}
		for _, entry := range dirEntries {
			if entry.IsDir() {
				if !opts.DirsAllowed {
					yield("", sberr.Wrap(
						FilesWithExtInDirErr,
						"Supplied dir (%s) contained a directory",
						dir,
					))
					return
				}
				continue
			}

			name := entry.Name()
			if path.Ext(name) != ext {
				if !opts.OtherFileTypesAllowed {
					yield("", sberr.Wrap(
						FilesWithExtInDirErr,
						"Supplied dir (%s) contained non-csv files",
						dir,
					))
					return
				}
				continue
			}

			if !yield(path.Join(dir, entry.Name()), nil) {
				return
			}
		}
	}
}
