package main

import "fmt"

// main parses options and runs the command determined by those options
func main() {
	opts := &opts{}
	if cmd, err := opts.parse(); err != nil {
		fmt.Println("ERROR:", err)
	} else {
		cmd.run(opts)
	}
}
