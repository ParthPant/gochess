package main

import (
	"github.com/ParthPant/gochess/core"
	"github.com/ParthPant/gochess/ui"
)

func main() {
	g := ui.CreateGui(core.NewGame(core.White), 800)
	g.GameLoop()
}
