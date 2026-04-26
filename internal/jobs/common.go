package jobs

import (
	"fmt"
	"sync/atomic"
)

var (
	UID_CNTR atomic.Int64
)

func formatSchedulerLogLine(name string, uid int64, msg string) string {
	return fmt.Sprintf("SCHED: %s (UID: %d): %s", name, uid, msg)
}

func formatJobLogLine(name string, uid int64, msg string) string {
	return fmt.Sprintf("JOB: %s (UID: %d): %s", name, uid, msg)
}
