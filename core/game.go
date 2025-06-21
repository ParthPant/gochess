package core

import (
	"fmt"
	"log/slog"

	"github.com/ParthPant/gochess/util"
)

type ChessGame struct {
	HumanColor Color
	board      Board
	history    util.Stack[Board]
}

func NewGame(humanColor Color) ChessGame {
	board, err := BoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	if err != nil {
		panic("Error while constructing default fen board.")
	}
	if board.hash != board.calculateHash() {
		panic("Error: Zobrist has not set while construction.")
	}
	return ChessGame{
		humanColor,
		board,
		util.NewStack[Board](),
	}
}

func (g *ChessGame) Board() *Board {
	return &g.board
}

// MakeMove will make a new move in the Game.
// promPiece will be used only if the move is a promotion.
func (g *ChessGame) MakeMove(from Square, to Square, promPiece promotedPiece) (Move, bool) {
	legal_moves := g.GetLegalPieceMovesBB(from)
	if !legal_moves.IsSet(to) {
		slog.Info("Move is Illegal", "from", from.ToStr(), "to", to.ToStr())
		return Move{}, false
	}
	move, ok := g.board.inferMove(from, to)
	if !ok {
		slog.Info("Unable to infer move.", "from", from.ToStr(), "to", from.ToStr())
		return Move{}, false
	}
	if move.IsPromotion() {
		move.SetPromPiece(promPiece)
	}
	return move, g.makeMoveImpl(move)
}

func (g *ChessGame) UndoPreviousMove() {
	prev_state, ok := g.history.Pop()
	if !ok {
		slog.Error("No more history to undo.")
		return
	}
	slog.Info("Restored game to the previous state.")
	g.board = prev_state
}

// Implementation of MakeMove. It will also make a new entry in the history
func (g *ChessGame) makeMoveImpl(m Move) bool {
	slog.Debug("Making Move:", "move", m.ToStr())
	// Move is first made on a copy of the board.
	board_copy, valid := g.board.makeMove(m)
	if valid {
		slog.Info("Valid Move.", "move", m.ToStr())
		// if the move is valid the board is replaced with it's copy
		g.history.Push(g.board)
		g.board = board_copy

		calculated_hash := g.board.calculateHash()
		slog.Debug("Board Hash",
			slog.String("hash", fmt.Sprintf("0x%x", g.board.hash)),
			slog.String("calculated", fmt.Sprintf("0x%x", calculated_hash)))
		if g.board.hash != calculated_hash {
			panic("Error: Hashes do not match")
		}
	} else {
		slog.Info("Invalid Move.", "move", m.ToStr())
	}
	return valid
}

func (g *ChessGame) GetLegalPieceMoves(sq Square) MoveList {
	move_list, _ := g.board.getLegalMoves(sq)
	return move_list
}

func (g *ChessGame) GetLegalPieceMovesBB(sq Square) BitBoard {
	move_list, _ := g.board.getLegalMoves(sq)
	return move_list.ToBB()
}

func (g *ChessGame) GetAllLegalMoves(side Color) MoveList {
	return g.board.getAllLegalMoves(side)
}
