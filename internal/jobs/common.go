package jobs

import (
	"fmt"
	"sync/atomic"
)

var (
	UID_CNTR atomic.Int64
)

func formatJobLogLine(name string, uid int64, msg string) string {
	return fmt.Sprintf("JOB: %s (UID: %d): %s", name, uid, msg)
}
