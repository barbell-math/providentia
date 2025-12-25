package dal

import (
	"context"
	"log/slog"

	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5/tracelog"
)

type (
	PgxLogAdapter struct {
		*slog.Logger
	}
)

func (p *PgxLogAdapter) Log(
	ctxt context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	// TODO - try to make smarter so that the list is only created if the log
	// level is high enough
	args := make([]any, len(data)*2)
	cntr := 0
	for k, v := range data {
		args[cntr] = k
		args[cntr+1] = v
		cntr += 2
	}

	switch level {
	case tracelog.LogLevelDebug:
		p.Logger.Log(ctxt, sblog.VLevel(5), "PGX: "+msg, args...)
	case tracelog.LogLevelInfo:
		p.Logger.Log(ctxt, sblog.VLevel(4), "PGX: "+msg, args...)
	case tracelog.LogLevelWarn:
		p.Logger.Warn("PGX: "+msg, args...)
	case tracelog.LogLevelError:
		p.Logger.Error("PGX: "+msg, args...)
	}
}
