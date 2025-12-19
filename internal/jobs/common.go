package jobs

import (
	"fmt"
	"sync/atomic"
)

var (
	UID_CNTR atomic.Uint64
)

func formatJobLogLine(name string, uid uint64, msg string) string {
	return fmt.Sprintf("JOB: %s (UID: %d): %s", name, uid, msg)
}
