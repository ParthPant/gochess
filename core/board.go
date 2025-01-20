package core

import (
	"bytes"
	"fmt"
)

type Board struct {
	bitBoards     [12]BitBoard
	halfMoveClock uint
	fullMoveClock uint
	activeColor   Color
	epTarget      *Square
	castlingFlags uint8
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
	epTarget := "N/A"
	if b.epTarget != nil {
		epTarget = b.epTarget.ToStr()
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

func (b *Board) getColorOccupancy(c Color) BitBoard {
	if c == White {
		return b.whiteOccupancy()
	} else {
		return b.blackOccupancy()
	}
}

func (b *Board) whiteOccupancy() BitBoard {
	var occ BitBoard
	for i := 0; i < 6; i++ {
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
	if (1<<pos)&EightRank == 0 {
		// not in eighth rank
		quietMoves |= BitBoard(1<<(pos+8)) & ^friendly & ^enemy
	}
	if (1<<pos)&SecondRank > 0 {
		// in second rank
		quietMoves |= BitBoard(1<<(pos+16)) & ^friendly & ^enemy
	}

	attacks := PawnAtkTable[White][pos] & enemy

	if b.epTarget != nil && attacks&(1<<*b.epTarget) > 0 {
		attacks |= (1 << *b.epTarget)
	}

	return quietMoves | attacks
}

func (b *Board) blackPawnMoves(pos Square) BitBoard {
	var quietMoves BitBoard
	friendly := b.getColorOccupancy(White)
	enemy := b.getColorOccupancy(Black)
	if (1<<pos)&FirstRank == 0 {
		// not in first rank
		quietMoves |= BitBoard(1<<(pos-8)) & ^friendly & ^enemy
	}
	if (1<<pos)&SeventhRank > 0 {
		// in seventh rank
		quietMoves |= BitBoard(1<<(pos-16)) & ^friendly & ^enemy
	}

	attacks := PawnAtkTable[Black][pos] & enemy

	if b.epTarget != nil && attacks&(1<<*b.epTarget) > 0 {
		attacks |= (1 << *b.epTarget)
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
		if (allOccupancy&(1<<D1) == 0) && (allOccupancy&(1<<C1) == 0) && (allOccupancy&(1<<B1) == 0) && !b.isSqAttacked(G1, Black) {
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
		if (allOccupancy&(1<<D8) == 0) && (allOccupancy&(1<<C8) == 0) && (allOccupancy&(1<<B8) == 0) && !b.isSqAttacked(G8, White) {
			moves = moves.Set(C8)
		}
	}
	return moves
}

func (b *Board) isSqAttacked(sq Square, attackColor Color) bool {
	if attackColor == White {
		if (PawnAtkTable[Black][sq] & b.bitBoards[Pb]) > 0 {
			return true
		} else if (KnightAtkTable[sq] & b.bitBoards[Nb]) > 0 {
			return true
		} else if (GetBishopMoves(sq, b.blackOccupancy()&b.whiteOccupancy()) & (b.bitBoards[Bb] | b.bitBoards[Qb])) > 0 {
			return true
		} else if (GetRookMoves(sq, b.blackOccupancy()&b.whiteOccupancy()) & (b.bitBoards[Rb] | b.bitBoards[Qb])) > 0 {
			return true
		} else {
			return false
		}
	} else {
		if (PawnAtkTable[White][sq] & b.bitBoards[Pw]) > 0 {
			return true
		} else if (KnightAtkTable[sq] & b.bitBoards[Nw]) > 0 {
			return true
		} else if (GetBishopMoves(sq, b.blackOccupancy()&b.whiteOccupancy()) & (b.bitBoards[Bw] | b.bitBoards[Qw])) > 0 {
			return true
		} else if (GetRookMoves(sq, b.blackOccupancy()&b.whiteOccupancy()) & (b.bitBoards[Rw] | b.bitBoards[Qw])) > 0 {
			return true
		} else {
			return false
		}
	}
}

func (b *Board) getPieceMoves(sq Square) (BitBoard, bool) {
	piece, ok := b.GetAtSq(sq)
	if !ok {
		return 0, false
	}

	friendlyColor := piece.GetColor()
	enemyColor := White
	if enemyColor == friendlyColor {
		enemyColor = Black
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
