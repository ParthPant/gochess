package core

type AI interface {
	GetBestMove(b *Board) (Move, bool)
}

type NegaMaxAI struct {
	evalMethod func(b *Board) int32
	depth      uint8
}

func NewNegaMaxAI() NegaMaxAI {
	return NegaMaxAI{
		evaluateBoard,
		5,
	}
}

func (nmax *NegaMaxAI) GetBestMove(b *Board) (Move, bool) {
	value := MinScore
	var bestMove Move
	found := false
	for _, move := range b.getAllLegalMoves(b.activeColor) {
		if board_copy, ok := b.makeMove(move); ok {
			move_score := -nmax.negamax(board_copy, nmax.depth-1, MinScore, MaxScore)
			if move_score > value {
				value = move_score
				bestMove = move
				found = true
			}
		}
	}
	return bestMove, found
}

func (nmax *NegaMaxAI) negamax(b Board, depth uint8, alpha int32, beta int32) int32 {
	if depth == 0 {
		return nmax.evalMethod(&b)
	}
	if b.isActiveSideInCheck() {
		return MatingScore
	}
	value := MinScore
	for _, move := range b.getAllLegalMoves(b.activeColor) {
		if board_copy, ok := b.makeMove(move); ok {
			value = max(value, -nmax.negamax(board_copy, depth-1, -beta, -alpha))
			alpha = max(alpha, value)
			if alpha >= beta {
				break
			}
		}
	}
	return value
}
