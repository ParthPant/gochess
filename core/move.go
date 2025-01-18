package core

const quietMove uint8 = 0b0000
const doublePawnPush uint8 = 0b0001
const kingCastle uint8 = 0b0010
const queenCastle uint8 = 0b0011
const epCapture uint8 = 0b0101
const captureMask uint8 = 0b0100
const promotionMask uint8 = 0b1000

type Move struct {
	flags uint8
	from  square
	to    square
}

type promotedPiece uint8

const (
	Knight promotedPiece = iota
	Bishop
	Rook
	Queen
)

func NewQuietMove(from square, to square) Move {
	return Move{
		flags: 0,
		from:  from,
		to:    to,
	}
}

func NewCaptureMove(from square, to square) Move {
	return Move{
		flags: captureMask,
		from:  from,
		to:    to,
	}
}

func NewEpMove(from square, to square) Move {
	return Move{
		flags: epCapture,
		from:  from,
		to:    to,
	}
}

func NewKingCastle(from square, to square) Move {
	return Move{
		flags: kingCastle,
		from:  from,
		to:    to,
	}
}

func NewQueenCastle(from square, to square) Move {
	return Move{
		flags: queenCastle,
		from:  from,
		to:    to,
	}
}

func NewPromotionMove(from square, to square, prom promotedPiece) Move {
	flags := promotionMask | uint8(prom)
	return Move{
		flags,
		from,
		to,
	}
}

func NewPromotionCapture(from square, to square, prom promotedPiece) Move {
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

func (m Move) GetPromPiece() promotedPiece {
	return promotedPiece(m.flags & 0b1100)
}
