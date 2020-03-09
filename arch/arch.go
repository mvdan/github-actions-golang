package arch

import (
	"fmt"
	"runtime"
)

func print() {
	fmt.Printf("GOARCH: %s\n", runtime.GOARCH)
}
