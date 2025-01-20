package core

type ChessGame struct {
	board      Board
	HumanColor Color
}

func NewGame(humanColor Color) ChessGame {
	// board, err := BoardFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	board, err := BoardFromFen("rn1qkbn1/ppp1pppp/2b5/3p1r1B/1KP2P2/7Q/PBPPP1PP/RN4NR w KQkq - 0 1")
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

func (g *ChessGame) MakeMove(m Move) bool {
	return true
}

func (g *ChessGame) CalculateLegalMoves() []Move {
	return []Move{}
}

func (g *ChessGame) GetPieceMoves(sq Square) (BitBoard, bool) {
	return g.board.getPieceMoves(sq)
}
