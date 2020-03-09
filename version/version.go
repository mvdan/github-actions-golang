package version

import (
	"fmt"
	"runtime"
)

func Print() {
	fmt.Printf("Go version: %s\n", runtime.Version())
}
