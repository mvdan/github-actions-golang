package version

import (
	"fmt"
	"runtime"
)

func print() {
	fmt.Printf("Go version: %s\n", runtime.Version())
}
