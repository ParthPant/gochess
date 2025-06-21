package core

import (
	"bytes"
	"fmt"
	"log/slog"
)

type epTarget struct {
	exists bool
	sq     Square
}

type Board struct {
	bitBoards     [12]BitBoard
	halfMoveClock uint
	fullMoveClock uint
	activeColor   Color
	epTarget      epTarget
	castlingFlags uint8
	hash          uint64
}

func (ept *epTarget) set(sq Square) {
	ept.exists = true
	ept.sq = sq
}

func (ept *epTarget) clear() {
	ept.exists = false
}

func (ept *epTarget) get() (Square, bool) {
	if ept.exists {
		return ept.sq, true
	}
	return 0, false
}

func (b *Board) Print() {
	var buf bytes.Buffer

	for y := 7; y >= 0; y-- {
		buf.WriteString(fmt.Sprintf("%d", y+1))
		for x := 0; x <= 7; x++ {
			sq := Square(y*8 + x)
			piece, ok := b.GetAtSq(sq)
			if ok {
				buf.WriteString(fmt.Sprintf("  %c", piece.UtfRune()))
			} else {
				buf.WriteString("  •")
			}
		}
		buf.WriteRune('\n')
	}
	buf.WriteString(fmt.Sprintf("   A  B  C  D  E  F  G  H"))
	ept, eptExists := b.epTarget.get()
	epTarget := "N/A"
	if eptExists {
		epTarget = ept.ToStr()
	}
	buf.WriteString(fmt.Sprintf("\nActive Color: %s", b.activeColor.ToStr()))
	buf.WriteString(fmt.Sprintf("\nEn-Passant Target: %s", epTarget))
	buf.WriteString(fmt.Sprintf("\nHalf Move: %d\tFull Move: %d", b.halfMoveClock, b.fullMoveClock))
	buf.WriteString(fmt.Sprintf("\nCastling Flags: %04b", b.castlingFlags))
	fmt.Println(buf.String())
}

func (b *Board) GetAtSq(sq Square) (Piece, bool) {
	for _, piece := range BoardPieces {
		if b.bitBoards[piece].IsSet(sq) {
			return piece, true
		}
	}
	return 0, false

}

// inferMove Creates a Move object with all the relevant flags
// using only the source and destination square
// IMPORTANT: For promotions, the promoted piece will always be set to Knight by default.
// Use SetPromPiece() to set the correct promoted piece.
// TODO: inferMove should do all sanity checks on a move
func (b *Board) inferMove(from Square, to Square) (Move, bool) {
	piece, occupied := b.GetAtSq(from)
	if !occupied {
		return Move{}, false
	}

	flags := uint8(0)
	_, occupied = b.GetAtSq(to)
	// set capture flags
	if occupied {
		flags |= captureMask
	}
	if piece == Pw || piece == Pb {
		ept, eptExists := b.epTarget.get()
		if eptExists && to == ept {
			flags |= epCapture
		}
	}

	// set double pawn push flags
	if piece == Pw {
		if (1<<from)&SecondRank > 0 && (1<<to)&FourthRank > 0 {
			flags |= doublePawnPush
		}
	} else if piece == Pb {
		if (1<<from)&SeventhRank > 0 && (1<<to)&FifthRank > 0 {
			flags |= doublePawnPush
		}
	}

	// set castle flags
	if piece == Kw {
		if from == E1 && to == G1 {
			flags |= kingCastle
		}
		if from == E1 && to == C1 {
			flags |= queenCastle
		}
	} else if piece == Kb {
		if from == E8 && to == G8 {
			flags |= kingCastle
		}
		if from == E8 && to == C8 {
			flags |= queenCastle
		}
	}

	// set promotion mask
	if piece == Pw {
		if (1<<to)&EighthRank > 0 {
			flags |= promotionMask
		}
	} else if piece == Pb {
		if (1<<to)&FirstRank > 0 {
			flags |= promotionMask
		}
	}

	return Move{
		flags,
		from,
		to,
	}, true
}

// makeMove can potentially mutate Board to an invalid state.
// State is guaranteed to be valid if the return value is True
// TODO: makeMove should not do any sanity checks on a move. Sanity checks should be done by inferMove
func (b Board) makeMove(m Move) (Board, bool) {
	// get the moving_piece which the player has moved
	moving_piece, occupied := b.GetAtSq(m.from)
	if !occupied {
		slog.Error("No piece to be moved.")
		return b, false
	}

	if moving_piece.GetColor() != b.activeColor {
		slog.Error("Piece is not of the active color", "color", b.activeColor.ToStr())
		return b, false
	}
	slog.Debug("Moving Piece", "piece", string(moving_piece.Char()))

	// Handle captured pieces first
	if m.IsCapture() && !m.IsEp() {
		captured_piece, occupied := b.GetAtSq(m.to)
		slog.Debug("Handling capture piece", "piece", string(captured_piece.Char()), "color", moving_piece.GetColor().ToStr())
		if occupied && captured_piece.GetColor() == b.activeColor {
			slog.Error("Captured piece is not of opponent color.")
			return b, false
		} else if !occupied {
			slog.Error("Captured position is not occupied.")
			return b, false
		}
		b.bitBoards[captured_piece] = b.bitBoards[captured_piece].UnSet(m.to)
		b.hash ^= ZobPieceKeys[captured_piece][m.to]
	} else if m.IsEp() {
		captured_square := Square(m.to - 8)
		if moving_piece == Pb {
			captured_square = Square(m.to + 8)
		}
		captured_piece, occupied := b.GetAtSq(captured_square)
		slog.Debug("Handling capture piece", "piece", string(captured_piece.Char()), "color", moving_piece.GetColor().ToStr())
		if !occupied {
			slog.Error("Invalid state: ep move has no piece to be captured.")
			return b, false
		}
		b.bitBoards[captured_piece] = b.bitBoards[captured_piece].UnSet(captured_square)
		b.hash ^= ZobPieceKeys[captured_piece][captured_square]
	}

	// Now handle the moving piece
	if m.IsQuiet() || m.IsDoublePawnPush() || (m.IsCapture() && !m.IsPromotion()) {
		b.bitBoards[moving_piece] = b.bitBoards[moving_piece].UnSet(m.from)
		b.hash ^= ZobPieceKeys[moving_piece][m.from]
		b.bitBoards[moving_piece] = b.bitBoards[moving_piece].Set(m.to)
		b.hash ^= ZobPieceKeys[moving_piece][m.to]
	} else if m.IsPromotion() {
		b.bitBoards[moving_piece] = b.bitBoards[moving_piece].UnSet(m.from)
		b.hash ^= ZobPieceKeys[moving_piece][m.from]
		promoted_piece := m.GetPromPiece().WithColor(b.activeColor)
		b.bitBoards[promoted_piece] = b.bitBoards[promoted_piece].Set(m.to)
		b.hash ^= ZobPieceKeys[promoted_piece][m.to]
	} else if m.IsKingCastle() {
		if b.activeColor == White {
			if !b.CanWhiteOO() {
				slog.Error("White can not castle kingside.")
				return b, false
			}
			b.bitBoards[Kw] = b.bitBoards[Kw].UnSet(E1)
			b.hash ^= ZobPieceKeys[Kw][E1]
			b.bitBoards[Kw] = b.bitBoards[Kw].Set(G1)
			b.hash ^= ZobPieceKeys[Kw][G1]
			b.bitBoards[Rw] = b.bitBoards[Rw].UnSet(H1)
			b.hash ^= ZobPieceKeys[Rw][H1]
			b.bitBoards[Rw] = b.bitBoards[Rw].Set(F1)
			b.hash ^= ZobPieceKeys[Rw][F1]
		} else {
			if !b.CanBlackOO() {
				slog.Error("Black can not castle kingside.")
				return b, false
			}
			b.bitBoards[Kb] = b.bitBoards[Kb].UnSet(E8)
			b.hash ^= ZobPieceKeys[Kb][E8]
			b.bitBoards[Kb] = b.bitBoards[Kb].Set(G8)
			b.hash ^= ZobPieceKeys[Kb][G8]
			b.bitBoards[Rb] = b.bitBoards[Rb].UnSet(H8)
			b.hash ^= ZobPieceKeys[Rb][H8]
			b.bitBoards[Rb] = b.bitBoards[Rb].Set(F8)
			b.hash ^= ZobPieceKeys[Rb][F8]
		}
	} else if m.IsQueenCastle() {
		if b.activeColor == White {
			if !b.CanWhiteOOO() {
				slog.Error("White cannot castle queenside")
				return b, false
			}
			b.bitBoards[Kw] = b.bitBoards[Kw].UnSet(E1)
			b.hash ^= ZobPieceKeys[Kw][E1]
			b.bitBoards[Kw] = b.bitBoards[Kw].Set(C1)
			b.hash ^= ZobPieceKeys[Kw][C1]
			b.bitBoards[Rw] = b.bitBoards[Rw].UnSet(A1)
			b.hash ^= ZobPieceKeys[Rw][A1]
			b.bitBoards[Rw] = b.bitBoards[Rw].Set(D1)
			b.hash ^= ZobPieceKeys[Rw][D1]
		} else {
			if !b.CanBlackOOO() {
				slog.Error("Black cannot castle queenside")
				return b, false
			}
			b.bitBoards[Kb] = b.bitBoards[Kb].UnSet(E8)
			b.hash ^= ZobPieceKeys[Kb][E8]
			b.bitBoards[Kb] = b.bitBoards[Kb].Set(C8)
			b.hash ^= ZobPieceKeys[Kb][C8]
			b.bitBoards[Rb] = b.bitBoards[Rb].UnSet(A8)
			b.hash ^= ZobPieceKeys[Rb][A8]
			b.bitBoards[Rb] = b.bitBoards[Rb].Set(D8)
			b.hash ^= ZobPieceKeys[Rb][D8]
		}
	}

	// update if castling is no longer possible
	b.hash ^= ZobCastleKeys[b.castlingFlags]
	if moving_piece == Kw {
		b.unsetWhiteOO()
		b.unsetWhiteOOO()
	} else if moving_piece == Kb {
		b.unsetBlackOO()
		b.unsetBlackOOO()
	} else if moving_piece == Rw {
		if m.from == A1 {
			b.unsetWhiteOOO()
		} else if m.from == H1 {
			b.unsetWhiteOO()
		}
	} else if moving_piece == Rb {
		if m.from == A8 {
			b.unsetBlackOOO()
		} else if m.from == H8 {
			b.unsetBlackOO()
		}
	}
	b.hash ^= ZobCastleKeys[b.castlingFlags]

	// EP updates
	if t, exists := b.epTarget.get(); exists {
		b.hash ^= ZobEpKeys[t]
	}
	if m.IsDoublePawnPush() {
		var t Square
		if b.activeColor == White {
			t = Square(m.to - 8)
			b.epTarget.set(t)
		} else {
			t = Square(m.to + 8)
			b.epTarget.set(t)
		}
	} else {
		b.epTarget.clear()
	}
	if t, exists := b.epTarget.get(); exists {
		b.hash ^= ZobEpKeys[t]
	}

	// increment move clocks
	b.halfMoveClock += 1
	if b.activeColor == Black {
		b.fullMoveClock += 1
	}

	// toggle active color
	b.activeColor = 1 ^ b.activeColor
	b.hash ^= ZobBlackToMoveKey

	return b, true
}

func (b *Board) isMoveLegal(m Move) bool {
	piece, occupied := b.GetAtSq(m.from)
	if !occupied {
		slog.Debug("Moving piece is not occupied.")
		return false
	}
	moving_color := piece.GetColor()

	board_copy, valid := b.makeMove(m)
	if !valid {
		slog.Debug("The move is invalid.")
		return false
	}

	// get the position of the king
	var king_sq Square
	var king_occupied bool
	if moving_color == White {
		king_sq, king_occupied = board_copy.bitBoards[Kw].Peek()
	} else {
		king_sq, king_occupied = board_copy.bitBoards[Kb].Peek()
	}

	if !king_occupied {
		slog.Debug("King is not occupied.")
		return false
	}
	slog.Debug("Checking if king is attacked.", "king_sq", king_sq.ToStr(), "attacking_color", moving_color^1)
	return !board_copy.isSqAttacked(king_sq, moving_color^1)
}

func (b *Board) getColorOccupancy(c Color) BitBoard {
	if c == White {
		return b.whiteOccupancy()
	} else {
		return b.blackOccupancy()
	}
}

func (b *Board) whiteOccupancy() BitBoard {
	var occ BitBoard
	for i := range 6 {
		occ |= b.bitBoards[i]
	}
	return occ
}

func (b *Board) blackOccupancy() BitBoard {
	var occ BitBoard
	for i := 6; i < 12; i++ {
		occ |= b.bitBoards[i]
	}
	return occ
}

func (b *Board) setWhiteOO() {
	b.castlingFlags |= 1
}

func (b *Board) unsetWhiteOO() {
	b.castlingFlags &= ^uint8(1)
}

func (b *Board) CanWhiteOO() bool {
	return b.castlingFlags&1 > 0
}

func (b *Board) setWhiteOOO() {
	b.castlingFlags |= (1 << 1)
}

func (b *Board) unsetWhiteOOO() {
	b.castlingFlags &= ^uint8(1 << 1)
}

func (b *Board) CanWhiteOOO() bool {
	return b.castlingFlags&(1<<1) > 0
}

func (b *Board) setBlackOO() {
	b.castlingFlags |= (1 << 2)
}

func (b *Board) unsetBlackOO() {
	b.castlingFlags &= ^uint8(1 << 2)
}

func (b *Board) CanBlackOO() bool {
	return b.castlingFlags&(1<<2) > 0
}

func (b *Board) setBlackOOO() {
	b.castlingFlags |= (1 << 3)
}

func (b *Board) unsetBlackOOO() {
	b.castlingFlags &= ^uint8(1 << 3)
}

func (b *Board) CanBlackOOO() bool {
	return b.castlingFlags&(1<<3) > 0
}

func (b *Board) whitePawnMoves(pos Square) BitBoard {
	var quietMoves BitBoard
	friendly := b.getColorOccupancy(White)
	enemy := b.getColorOccupancy(Black)
	if (1<<pos)&EighthRank == 0 {
		// not in eighth rank
		quietMoves |= BitBoard(1<<(pos+8)) & ^friendly & ^enemy
	}
	if (1<<pos)&SecondRank > 0 {
		// in second rank
		quietMoves |= BitBoard(1<<(pos+16)) & ^friendly & ^enemy
	}

	attacks := PawnAtkTable[White][pos] & enemy

	ept, eptExists := b.epTarget.get()
	if eptExists && PawnAtkTable[White][pos]&(1<<ept) > 0 {
		attacks |= (1 << ept)
	}

	return quietMoves | attacks
}

func (b *Board) blackPawnMoves(pos Square) BitBoard {
	var quietMoves BitBoard
	friendly := b.getColorOccupancy(Black)
	enemy := b.getColorOccupancy(White)
	if (1<<pos)&FirstRank == 0 {
		// not in first rank
		quietMoves |= BitBoard(1<<(pos-8)) & ^friendly & ^enemy
	}
	if (1<<pos)&SeventhRank > 0 {
		// in seventh rank
		quietMoves |= BitBoard(1<<(pos-16)) & ^friendly & ^enemy
	}

	attacks := PawnAtkTable[Black][pos] & enemy

	ept, eptExists := b.epTarget.get()
	if eptExists && PawnAtkTable[Black][pos]&(1<<ept) > 0 {
		attacks |= (1 << ept)
	}

	return quietMoves | attacks
}

func (b *Board) knightMoves(pos Square) BitBoard {
	piece, _ := b.GetAtSq(pos)
	friendlyColor := piece.GetColor()
	friendly := b.getColorOccupancy(friendlyColor)
	return KnightAtkTable[pos] & ^friendly
}

func (b *Board) rookMoves(pos Square) BitBoard {
	piece, _ := b.GetAtSq(pos)
	friendlyColor := piece.GetColor()
	enemyColor := White
	if enemyColor == friendlyColor {
		enemyColor = Black
	}
	friendly := b.getColorOccupancy(friendlyColor)
	enemy := b.getColorOccupancy(enemyColor)
	return GetRookMoves(pos, friendly|enemy) & ^friendly
}

func (b *Board) bishopMoves(pos Square) BitBoard {
	piece, _ := b.GetAtSq(pos)
	friendlyColor := piece.GetColor()
	enemyColor := White
	if enemyColor == friendlyColor {
		enemyColor = Black
	}
	friendly := b.getColorOccupancy(friendlyColor)
	enemy := b.getColorOccupancy(enemyColor)
	return GetBishopMoves(pos, friendly|enemy) & ^friendly
}

func (b *Board) queenMoves(pos Square) BitBoard {
	return b.bishopMoves(pos) | b.rookMoves(pos)
}

func (b *Board) whiteKingMoves(pos Square) BitBoard {
	piece, _ := b.GetAtSq(pos)
	friendlyColor := piece.GetColor()
	enemyColor := White
	if enemyColor == friendlyColor {
		enemyColor = Black
	}
	friendly := b.getColorOccupancy(friendlyColor)
	enemy := b.getColorOccupancy(enemyColor)

	moves := KingAtkTable[pos] & ^friendly
	allOccupancy := friendly & enemy
	if pos == E1 && b.CanWhiteOO() {
		if (allOccupancy&(1<<F1) == 0) && (allOccupancy&(1<<G1) == 0) && !b.isSqAttacked(G1, Black) {
			moves = moves.Set(G1)
		}
	}
	if pos == E1 && b.CanWhiteOOO() {
		if (allOccupancy&(1<<D1) == 0) && (allOccupancy&(1<<C1) == 0) && (allOccupancy&(1<<B1) == 0) && !b.isSqAttacked(C1, Black) {
			moves = moves.Set(C1)
		}
	}
	return moves
}

func (b *Board) blackKingMoves(pos Square) BitBoard {
	piece, _ := b.GetAtSq(pos)
	friendlyColor := piece.GetColor()
	enemyColor := White
	if enemyColor == friendlyColor {
		enemyColor = Black
	}
	friendly := b.getColorOccupancy(friendlyColor)
	enemy := b.getColorOccupancy(enemyColor)

	moves := KingAtkTable[pos] & ^friendly
	allOccupancy := friendly & enemy
	if pos == E8 && b.CanBlackOO() {
		if (allOccupancy&(1<<F8) == 0) && (allOccupancy&(1<<G8) == 0) && !b.isSqAttacked(G8, White) {
			moves = moves.Set(G8)
		}
	}
	if pos == E8 && b.CanBlackOOO() {
		if (allOccupancy&(1<<D8) == 0) && (allOccupancy&(1<<C8) == 0) && (allOccupancy&(1<<B8) == 0) && !b.isSqAttacked(C8, White) {
			moves = moves.Set(C8)
		}
	}
	return moves
}

func (b *Board) isSqAttacked(sq Square, attackColor Color) bool {
	if attackColor == White {
		if (PawnAtkTable[Black][sq] & b.bitBoards[Pw]) > 0 {
			slog.Debug("Square attacked by white Pawns", "square", sq.ToStr())
			return true
		} else if (KnightAtkTable[sq] & b.bitBoards[Nw]) > 0 {
			slog.Debug("Square attacked by white Knights", "square", sq.ToStr())
			return true
		} else if (GetBishopMoves(sq, b.blackOccupancy()|b.whiteOccupancy()) & (b.bitBoards[Bw] | b.bitBoards[Qw])) > 0 {
			slog.Debug("Square attacked by white Bishops", "square", sq.ToStr())
			return true
		} else if (GetRookMoves(sq, b.blackOccupancy()|b.whiteOccupancy()) & (b.bitBoards[Rw] | b.bitBoards[Qw])) > 0 {
			slog.Debug("Square attacked by white Rooks", "square", sq.ToStr())
			return true
		} else {
			return false
		}
	} else {
		if (PawnAtkTable[White][sq] & b.bitBoards[Pb]) > 0 {
			slog.Debug("Square attacked by black Pawns", "square", sq.ToStr())
			return true
		} else if (KnightAtkTable[sq] & b.bitBoards[Nb]) > 0 {
			slog.Debug("Square attacked by black Knights", "square", sq.ToStr())
			return true
		} else if (GetBishopMoves(sq, b.blackOccupancy()|b.whiteOccupancy()) & (b.bitBoards[Bb] | b.bitBoards[Qb])) > 0 {
			slog.Debug("Square attacked by black Bishops", "square", sq.ToStr())
			return true
		} else if (GetRookMoves(sq, b.blackOccupancy()|b.whiteOccupancy()) & (b.bitBoards[Rb] | b.bitBoards[Qb])) > 0 {
			slog.Debug("Square attacked by black Rooks", "square", sq.ToStr())
			return true
		} else {
			return false
		}

	}
}

// getPieceMovesBB returns a BitBoard of all possible *pseudolegal* moves for a piece at the given square
func (b *Board) getPieceMovesBB(sq Square) (BitBoard, bool) {
	piece, ok := b.GetAtSq(sq)
	if !ok {
		return 0, false
	}

	switch piece {
	case Pw:
		return b.whitePawnMoves(sq), true
	case Pb:
		return b.blackPawnMoves(sq), true
	case Kw:
		return b.whiteKingMoves(sq), true
	case Kb:
		return b.blackKingMoves(sq), true
	case Nw, Nb:
		return b.knightMoves(sq), true
	case Rw, Rb:
		return b.rookMoves(sq), true
	case Bw, Bb:
		return b.bishopMoves(sq), true
	case Qw, Qb:
		return b.queenMoves(sq), true
	default:
		panic("Invalid piece.")
	}
}

// getLegalMoves returns a BitBoard comprising of *legal* moves for a piece at the given square.
// TODO: Should this function return a "ok" bool? Such functions should probably use a ctx
func (b *Board) getLegalMoves(sq Square) (MoveList, bool) {
	moves_bb, ok := b.getPieceMovesBB(sq)
	move_list := []Move{}
	if !ok {
		return []Move{}, false
	}

	var to Square
	for true {
		moves_bb, to, ok = moves_bb.PopSq()
		if !ok {
			break
		}
		move, ok := b.inferMove(sq, to)
		if !ok {
			continue
		}
		slog.Debug("Checking Move", "move", move.ToStr())
		if b.isMoveLegal(move) {
			slog.Debug("Move legal", "move", move.ToStr())
			move_list = append(move_list, move)
		}
	}
	return move_list, true
}

// getAllLegalMoves returns a BitBoard of all the leagal mvoes for a side.
func (b *Board) getAllLegalMoves(side Color) MoveList {
	move_list := MoveList{}
	for piece := range 12 {
		if (side == White && piece < 6) || (side == Black && piece >= 6) {
			piece_occupancy := b.bitBoards[piece]
			sq := Square(0)
			ok := false
			for true {
				piece_occupancy, sq, ok = piece_occupancy.PopSq()
				if !ok {
					break
				}
				piece_move_list, ok := b.getLegalMoves(sq)
				if ok {
					move_list = append(move_list, piece_move_list...)
				}
			}
		}
	}
	return move_list
}

func (b *Board) calculateHash() uint64 {
	key := uint64(0)
	for square := range 64 {
		if piece, occupied := b.GetAtSq(Square(square)); occupied {
			key ^= ZobPieceKeys[piece][square]
		}
	}
	if ep_sq, exists := b.epTarget.get(); exists {
		key ^= ZobEpKeys[ep_sq]
	}
	key ^= ZobCastleKeys[b.castlingFlags]
	if b.activeColor == Black {
		key ^= ZobBlackToMoveKey
	}
	return key
}
