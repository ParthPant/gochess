package core

import "fmt"

const quietMove uint8 = 0b0000
const doublePawnPush uint8 = 0b0001
const kingCastle uint8 = 0b0010
const queenCastle uint8 = 0b0011
const epCapture uint8 = 0b0101
const captureMask uint8 = 0b0100
const promotionMask uint8 = 0b1000

type Move struct {
	flags uint8
	from  Square
	to    Square
}

type MoveList []Move

func (list *MoveList) ToBB() BitBoard {
	moves_bb := BitBoard(0)
	for _, move := range *list {
		moves_bb = moves_bb.Set(move.to)
	}
	return moves_bb
}

func NewQuietMove(from Square, to Square) Move {
	return Move{
		flags: 0,
		from:  from,
		to:    to,
	}
}

func NewCaptureMove(from Square, to Square) Move {
	return Move{
		flags: captureMask,
		from:  from,
		to:    to,
	}
}

func NewEpMove(from Square, to Square) Move {
	return Move{
		flags: epCapture,
		from:  from,
		to:    to,
	}
}

func NewKingCastle(from Square, to Square) Move {
	return Move{
		flags: kingCastle,
		from:  from,
		to:    to,
	}
}

func NewQueenCastle(from Square, to Square) Move {
	return Move{
		flags: queenCastle,
		from:  from,
		to:    to,
	}
}

func NewDoubelPawnPush(from Square, to Square) Move {
	return Move{
		flags: doublePawnPush,
		from:  from,
		to:    to,
	}
}

func NewPromotionMove(from Square, to Square, prom promotedPiece) Move {
	flags := promotionMask | uint8(prom)
	return Move{
		flags,
		from,
		to,
	}
}

func NewPromotionCapture(from Square, to Square, prom promotedPiece) Move {
	flags := captureMask | promotionMask | uint8(prom)
	return Move{
		flags,
		from,
		to,
	}
}

func (m Move) IsPromotion() bool {
	return m.flags&promotionMask > 0
}

func (m Move) IsCapture() bool {
	return m.flags&captureMask > 0
}

func (m Move) IsKingCastle() bool {
	return m.flags == kingCastle
}

func (m Move) IsQueenCastle() bool {
	return m.flags == queenCastle
}

func (m Move) IsEp() bool {
	return m.flags == epCapture
}

func (m Move) IsQuiet() bool {
	return m.flags == quietMove
}

func (m Move) IsDoublePawnPush() bool {
	return m.flags == doublePawnPush
}

func (m Move) GetPromPiece() promotedPiece {
	return promotedPiece(m.flags & 0b0011)
}

func (m *Move) SetPromPiece(p promotedPiece) {
	m.flags |= uint8(p)
}

func (m *Move) ToStr() string {
	return fmt.Sprintf("%s%s", m.from.ToStr(), m.to.ToStr())
}
