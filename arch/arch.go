package arch

import (
	"fmt"
	"runtime"
)

func Print() {
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
}
