package core

import (
	"errors"
	"fmt"
	u "unicode"
)

type Square uint8

func SquareFromXY(x int, y int) Square {
	return Square(y*8 + x)
}

func StrToSq(s string) (Square, error) {
	if len(s) > 2 {
		return 0, errors.New("Invalid square. 1")
	}
	var x uint8
	var y uint8

	if u.IsLetter(rune(s[0])) {
		x = uint8(rune(s[0]) - 'a')
	} else {
		return 0, errors.New("Invalid square.")
	}

	if u.IsDigit(rune(s[1])) {
		y = uint8(rune(s[1]) - '1')
	} else {
		return 0, errors.New("Invalid square.")
	}

	sq := Square(y*8 + x)
	if sq >= 64 {
		return 0, errors.New("Invalid square.")
	}
	return sq, nil
}

func (sq Square) ToStr() string {
	x, y := sq.ToXY()
	rank := y + 1
	file := rune(int('a') + int(x))

	return fmt.Sprintf("%c%d", file, rank)
}

func (sq Square) ToUint() uint8 {
	return uint8(sq)
}

func (sq Square) ToXY() (uint8, uint8) {
	x := sq % 8
	y := sq / 8
	return uint8(x), uint8(y)
}

const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

var MirrorSquare = [...]Square{
	A8, B8, C8, D8, E8, F8, G8, H8,
	A7, B7, C7, D7, E7, F7, G7, H7,
	A6, B6, C6, D6, E6, F6, G6, H6,
	A5, B5, C5, D5, E5, F5, G5, H5,
	A4, B4, C4, D4, E4, F4, G4, H4,
	A3, B3, C3, D3, E3, F3, G3, H3,
	A2, B2, C2, D2, E2, F2, G2, H2,
	A1, B1, C1, D1, E1, F1, G1, H1,
}
