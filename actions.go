// Copyright (c) 2019, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package actions

import (
	"fmt"
	"runtime"
)

func Demo() {
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("GOOS: %s\n", runtime.GOOS)
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
}
