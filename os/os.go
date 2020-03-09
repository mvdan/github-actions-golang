package os

import (
	"runtime"

	"dummy.module/actions/log"
)

func Print() {
	log.Infof("GOOS: %s\n", runtime.GOOS)
}
