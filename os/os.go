package os

import (
	"runtime"

	"github.com/ekrem95/github-actions-golang/log"
)

func print() {
	log.Infof("GOOS: %s\n", runtime.GOOS)
}
