package core

import "errors"

type Piece uint8
type Color uint8

const (
	White Color = 0
	Black Color = 1
)

func (c Color) ToStr() string {
	switch c {
	case 0:
		return "Black"
	case 1:
		return "White"
	default:
		panic("Invalid Color.")
	}
}

type promotedPiece uint8

const (
	Knight promotedPiece = iota
	Bishop
	Rook
	Queen
)

const (
	Nw = iota
	Bw
	Rw
	Qw
	Kw
	Pw
	Nb
	Bb
	Rb
	Qb
	Kb
	Pb
)

var BlackPieces = [...]Piece{Pb, Nb, Bb, Rb, Qb, Kb}
var WhitePieces = [...]Piece{Pw, Nw, Bw, Rw, Qw, Kw}
var BoardPieces = [...]Piece{Pb, Nb, Bb, Rb, Qb, Kb, Pw, Nw, Bw, Rw, Qw, Kw}

func (p promotedPiece) WithColor(c Color) Piece {
	if c == White {
		return Piece(p)
	} else {
		return Piece(p + 6)
	}
}

func (p Piece) GetColor() Color {
	if p < 6 {
		return White
	} else {
		return Black
	}
}

func (p Piece) Char() rune {
	r, ok := map[Piece]rune{Pb: 'p', Nb: 'n', Bb: 'b', Rb: 'r', Qb: 'q', Kb: 'k',
		Pw: 'P', Nw: 'N', Bw: 'B', Rw: 'R', Qw: 'Q', Kw: 'K'}[p]
	if !ok {
		panic("Unrecognized Piece.")
	}
	return r
}

func (p Piece) UtfRune() rune {
	r, ok := map[Piece]rune{Pb: '♟', Nb: '♞', Bb: '♝', Rb: '♜', Qb: '♛', Kb: '♚',
		Pw: '♙', Nw: '♘', Bw: '♗', Rw: '♖', Qw: '♕', Kw: '♔'}[p]
	if !ok {
		panic("Unrecognized Piece.")
	}
	return r
}

func CharToPiece(c rune) (Piece, error) {
	switch c {
	case 'p':
		return Pb, nil
	case 'n':
		return Nb, nil
	case 'b':
		return Bb, nil
	case 'r':
		return Rb, nil
	case 'q':
		return Qb, nil
	case 'k':
		return Kb, nil

	case 'P':
		return Pw, nil
	case 'N':
		return Nw, nil
	case 'B':
		return Bw, nil
	case 'R':
		return Rw, nil
	case 'Q':
		return Qw, nil
	case 'K':
		return Kw, nil

	default:
		return Pb, errors.New("Invalid char input.")
	}
}
