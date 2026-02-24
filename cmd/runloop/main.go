package main

import (
	"github.com/alecthomas/kong"
)

type cli struct {
	Version VersionCmd `cmd:"" default:"1"`
}

type VersionCmd struct{}

func (c *VersionCmd) Run() error {
	return nil
}

func main() {
	var cli cli
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
