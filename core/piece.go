package core

import "errors"

type piece uint8
type color uint8

const (
	White color = 0
	Black color = 1
)

func (c color) ToStr() string {
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

var BlackPieces = [...]piece{Pb, Nb, Bb, Rb, Qb, Kb}
var WhitePieces = [...]piece{Pw, Nw, Bw, Rw, Qw, Kw}
var BoardPieces = [...]piece{Pb, Nb, Bb, Rb, Qb, Kb, Pw, Nw, Bw, Rw, Qw, Kw}

func (p piece) GetColor() color {
	if p&(0b1000) > 0 {
		return Black
	} else {
		return White
	}
}

func (p piece) Char() rune {
	r, ok := map[piece]rune{Pb: 'p', Nb: 'n', Bb: 'b', Rb: 'r', Qb: 'q', Kb: 'k',
		Pw: 'P', Nw: 'N', Bw: 'B', Rw: 'R', Qw: 'Q', Kw: 'K'}[p]
	if !ok {
		panic("Unrecognized Piece.")
	}
	return r
}

func (p piece) UtfRune() rune {
	r, ok := map[piece]rune{Pb: '♟', Nb: '♞', Bb: '♝', Rb: '♜', Qb: '♛', Kb: '♚',
		Pw: '♙', Nw: '♘', Bw: '♗', Rw: '♖', Qw: '♕', Kw: '♔'}[p]
	if !ok {
		panic("Unrecognized Piece.")
	}
	return r
}

func CharToPiece(c rune) (piece, error) {
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
