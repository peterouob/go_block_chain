package main

import (
	"block_chain/cli"
	"os"
)

func main() {
	defer os.Exit(1)
	cmd := cli.CommandLine{}
	cmd.Run()
}
