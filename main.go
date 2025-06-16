package main

import (
	"log/slog"
	"os"

	"github.com/ParthPant/gochess/core"
	"github.com/ParthPant/gochess/ui"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	g := ui.CreateGui(core.NewGame(core.White), 800)
	g.GameLoop()
}
