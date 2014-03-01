package main

import (
	"fmt"
	"github.com/laher/someutils"
	"os"
)

func main() {
	err := someutils.Call("mv", os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Returned error: %v\n", err)
		os.Exit(1)
	}

}
