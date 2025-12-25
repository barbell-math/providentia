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

func NewTracelogWithAdapter(logger *slog.Logger, level int) *tracelog.TraceLog {
	return &tracelog.TraceLog{
		Logger:   &PgxLogAdapter{Logger: logger},
		LogLevel: translateLogLevel(level),
	}
}

func translateLogLevel(level int) tracelog.LogLevel {
	if level < 4 {
		// Always log warnings and errors, even if supplied level is 0
		return tracelog.LogLevelWarn
	}
	if level < 5 {
		return tracelog.LogLevelInfo
	}
	return tracelog.LogLevelDebug
}

func (p *PgxLogAdapter) Log(
	ctxt context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
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
