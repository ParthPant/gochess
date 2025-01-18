package main

import (
	"fmt"
	"log/slog"

	"github.com/ParthPant/gochess/core"
)

func main() {
	_, err := core.BoardFromFen("4k2r/6r1/8/8/8/8/3R4/R3K3 b - a8 54 2")
	if err != nil {
		slog.Error(fmt.Sprintf("Error: %s", err))
	}
	// board.Print()
}
