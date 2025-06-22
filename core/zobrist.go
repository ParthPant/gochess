package core

import (
	"log/slog"

	"github.com/ParthPant/gochess/util"
)

var ZobPieceKeys [12][64]uint64
var ZobEpKeys [64]uint64
var ZobCastleKeys [16]uint64
var ZobBlackToMoveKey uint64

func init() {
	slog.Info("Generating Zobrist keys.")
	var prng util.PRNG
	prng.Seed(2342342)
	for piece := range 12 {
		for sq := range 64 {
			ZobPieceKeys[piece][sq] = prng.SparseRand64()
		}
	}
	for sq := range 64 {
		ZobEpKeys[sq] = prng.SparseRand64()
	}
	for i := range 16 {
		ZobCastleKeys[i] = prng.SparseRand64()
	}
	ZobBlackToMoveKey = prng.SparseRand64()
}
