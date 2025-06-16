package core

import "log/slog"

type ChessGame struct {
	board      Board
	HumanColor Color
}

func NewGame(humanColor Color) ChessGame {
	board, err := BoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	if err != nil {
		panic("Error whilce constructing default fen board.")
	}
	return ChessGame{
		board,
		humanColor,
	}
}

func (g *ChessGame) Board() *Board {
	return &g.board
}

// MakeMove will make a new move in the Game.
// promPiece will be used only if the move is a promotion.
func (g *ChessGame) MakeMove(from Square, to Square, promPiece promotedPiece) (Move, bool) {
	legal_moves, ok := g.GetLegalPieceMovesBB(from)
	if !ok {
		return Move{}, false
	}
	if !legal_moves.IsSet(to) {
		slog.Info("Move is Illegal")
		return Move{}, false
	}
	move, ok := g.board.inferMove(from, to)
	if !ok {
		slog.Debug("Invalid Move", "from", from.ToStr(), "to", from.ToStr())
		return Move{}, false
	}
	if move.IsPromotion() {
		move.SetPromPiece(promPiece)
	}
	return move, g.makeMove(move)
}

// TODO: This should be a method of Board.
func (g *ChessGame) makeMove(m Move) bool {
	slog.Debug("Making Move:", "move", m.ToStr())
	// Move is first made on a copy of the board.
	board_copy := g.board.MakeCopy()
	valid := board_copy.makeMove(m)
	if valid {
		slog.Debug("Valid Move.")
		// if the move is valid the board is replaced with it's copy
		g.board = board_copy
	} else {
		slog.Debug("Invalid Move.")
	}
	return valid
}

func (g *ChessGame) GetLegalMoves() []Move {
	return []Move{}
}

func (g *ChessGame) GetLegalPieceMovesBB(sq Square) (BitBoard, bool) {
	moves_bb, ok := g.board.getLegalMoves(sq)
	return moves_bb, ok
}
