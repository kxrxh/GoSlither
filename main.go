package main

import (
	"os"

	"github.com/kxrxh/goslither/game"
	"github.com/kxrxh/goslither/terminal"
)

func main() {
	err := terminal.SetStdinNonBlocking()
	if err != nil {
		panic(err)
	}

	fd := int(os.Stdin.Fd())
	initState, err := terminal.SetRawMode(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, initState)
	
	game.StartNew()
}
