package core

import (
	"bytes"
	"fmt"
	"math/bits"
)

type BitBoard uint64

const (
	AFile        BitBoard = 0x0101010101010101
	BFile        BitBoard = AFile << 1
	GFile        BitBoard = AFile << 6
	HFile        BitBoard = AFile << 7
	FirstRank    BitBoard = 0x00000000000000FF
	EightRank    BitBoard = 0xFF00000000000000
	WhiteSquares BitBoard = 0x55AA55AA55AA55AA
	BlckSquares  BitBoard = 0xAA55AA55AA55AA55
)

func (b BitBoard) Print() {
	var buf bytes.Buffer

	for y := 7; y >= 0; y-- {
		buf.WriteString(fmt.Sprintf("%d", y+1))
		for x := 0; x <= 7; x++ {
			sq := square(y*8 + x)
			if b.IsSet(sq) {
				buf.WriteString(fmt.Sprintf("  🞣"))
			} else {
				buf.WriteString("  •")
			}
		}
		buf.WriteRune('\n')
	}
	buf.WriteString(fmt.Sprintf("   A  B  C  D  E  F  G  H"))
	fmt.Println(buf.String())
}

func (b BitBoard) Set(sq square) BitBoard {
	return b | (1 << sq)
}

func (b BitBoard) UnSet(sq square) BitBoard {
	return b & ^(1 << sq)
}

func (b BitBoard) IsSet(sq square) bool {
	return b&(1<<sq) != 0
}

func (b BitBoard) Move(sq1 square, sq2 square) BitBoard {
	return b.UnSet(sq1).Set(sq2)
}

func (b BitBoard) PopSq() (BitBoard, square, bool) {
	tzs := bits.TrailingZeros64(uint64(b))
	if tzs >= 64 {
		return 0, 0, false
	}
	return b.UnSet(square(tzs)), square(tzs), true
}

func (b BitBoard) Peek() (square, bool) {
	tzs := bits.TrailingZeros64(uint64(b))
	if tzs >= 64 {
		return 0, false
	}
	return square(tzs), true
}
