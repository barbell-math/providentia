package state

import (
	"context"
	"flag"
	"fmt"
	"log/slog"

	sbargp "github.com/barbell-math/smoothbrain-argparse"
	sblog "github.com/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	Conf struct {
		Logging sbargp.LoggingConf
		DB      sbargp.DBConf
	}

	State struct {
		Conf Conf
		DB   *pgxpool.Pool
		Log  *slog.Logger
	}

	ctxtKey struct{}
)

var (
	stateCtxtKey ctxtKey
)

func FromContext(ctxt context.Context) (*State, bool) {
	s, ok := ctxt.Value(stateCtxtKey).(*State)
	return s, ok
}

func Parse(ctxt context.Context, args []string) (context.Context, error) {
	var err error
	var poolConf *pgxpool.Config
	state := State{}

	if err = sbargp.Parse(&state.Conf, args, sbargp.ParserOpts[Conf]{
		ProgName: "providentia",
		RequiredArgs: []string{
			"DB.User", "DB.PswdEnvVar", "DB.Name",
		},
		ArgDefsSetter: func(conf *Conf, fs *flag.FlagSet) error {
			sbargp.Logging(fs, &conf.Logging, "Logging", sbargp.LoggingConf{
				Verbosity:       0,
				SaveTo:          "",
				Name:            "providentia",
				MaxNumLogs:      1,
				MaxLogSizeBytes: 1e6, // 1 MB
			})
			sbargp.DB(fs, &conf.DB, "DB", sbargp.DBConf{
				Host: "localhost",
				Port: 5432,
			})
			return nil
		},
	}); err != nil {
		goto done
	}

	if state.Log, err = sblog.New(sblog.Opts{
		CurVerbosityLevel: uint(state.Conf.Logging.Verbosity),
		RotateWriterOpts: sblog.RotateWriterOpts{
			LogDir:          string(state.Conf.Logging.SaveTo),
			LogName:         "providentia",
			MaxNumLogs:      uint(state.Conf.Logging.MaxNumLogs),
			MaxLogSizeBytes: uint64(state.Conf.Logging.MaxLogSizeBytes),
		},
	}); err != nil {
		return nil, err
	}

	if poolConf, err = pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		state.Conf.DB.Host,
		state.Conf.DB.Port,
		state.Conf.DB.User,
		state.Conf.DB.PswdEnvVar,
		state.Conf.DB.Name,
	)); err != nil {
		goto done
	}
	if state.DB, err = pgxpool.NewWithConfig(ctxt, poolConf); err != nil {
		goto done
	}
	if err = state.DB.Ping(ctxt); err != nil {
		goto done
	}

done:
	return context.WithValue(ctxt, stateCtxtKey, &state), err
}
