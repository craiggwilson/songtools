package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

var cli = flags.NewNamedParser("songtool", flags.Default)

func main() {

	if _, err := cli.Parse(); err != nil {
		os.Exit(1)
	}
}
