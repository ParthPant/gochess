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
	SecondRank   BitBoard = FirstRank << 8
	ThirdRank    BitBoard = SecondRank << 8
	FourthRank   BitBoard = ThirdRank << 8
	FifthRank    BitBoard = FourthRank << 8
	SixthRank    BitBoard = FifthRank << 8
	SeventhRank  BitBoard = SixthRank << 8
	EighthRank   BitBoard = SeventhRank << 8
	WhiteSquares BitBoard = 0x55AA55AA55AA55AA
	BlackSquares BitBoard = 0xAA55AA55AA55AA55
)

func (b BitBoard) Print() {
	var buf bytes.Buffer

	for y := 7; y >= 0; y-- {
		buf.WriteString(fmt.Sprintf("%d", y+1))
		for x := 0; x <= 7; x++ {
			sq := Square(y*8 + x)
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

func (b BitBoard) Set(sq Square) BitBoard {
	return b | (1 << sq)
}

func (b BitBoard) UnSet(sq Square) BitBoard {
	return b & ^(1 << sq)
}

func (b BitBoard) IsSet(sq Square) bool {
	return b&(1<<sq) != 0
}

func (b BitBoard) Move(sq1 Square, sq2 Square) BitBoard {
	return b.UnSet(sq1).Set(sq2)
}

func (b BitBoard) PopSq() (BitBoard, Square, bool) {
	tzs := bits.TrailingZeros64(uint64(b))
	if tzs >= 64 {
		return 0, 0, false
	}
	return b.UnSet(Square(tzs)), Square(tzs), true
}

func (b BitBoard) Peek() (Square, bool) {
	tzs := bits.TrailingZeros64(uint64(b))
	if tzs >= 64 {
		return 0, false
	}
	return Square(tzs), true
}
