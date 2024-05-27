package main

import (
	"fmt"
	"os"

	"github.com/blugnu/test-report/internal"
)

var osExit = os.Exit

// main parses options and runs the command determined by those options
func main() {
	opts := &internal.Options{}
	if cmd, err := opts.Parse(); err != nil {
		fmt.Println("ERROR:", err)
	} else {
		osExit(cmd.Run(opts))
	}
}
