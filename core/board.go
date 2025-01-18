package core

import (
	"bytes"
	"fmt"
)

type Board struct {
	bitBoards     [12]BitBoard
	halfMoveClock uint
	fullMoveClock uint
	activeColor   color
	epTarget      square
	castlingFlags uint8
	epPossible    bool
}

func (b *Board) Print() {
	var buf bytes.Buffer

	for y := 7; y >= 0; y-- {
		buf.WriteString(fmt.Sprintf("%d", y+1))
		for x := 0; x <= 7; x++ {
			sq := square(y*8 + x)
			piece, ok := b.getAtSq(sq)
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
	if b.epPossible {
		epTarget = b.epTarget.ToStr()
	}
	buf.WriteString(fmt.Sprintf("\nActive Color: %s", b.activeColor.ToStr()))
	buf.WriteString(fmt.Sprintf("\nEn-Passant Target: %s", epTarget))
	buf.WriteString(fmt.Sprintf("\nHalf Move: %d\tFull Move: %d", b.halfMoveClock, b.fullMoveClock))
	buf.WriteString(fmt.Sprintf("\nCastling Flags: %04b", b.castlingFlags))
	fmt.Println(buf.String())
}

func (b *Board) getAtSq(sq square) (piece, bool) {
	for _, piece := range BoardPieces {
		if b.bitBoards[piece].IsSet(sq) {
			return piece, true
		}
	}
	return 0, false
}

func (b *Board) SetWhiteOO() {
	b.castlingFlags |= 1
}

func (b *Board) UnsetWhiteOO() {
	b.castlingFlags &= ^uint8(1)
}

func (b *Board) CanWhiteOO() bool {
	return b.castlingFlags&1 > 0
}

func (b *Board) SetWhiteOOO() {
	b.castlingFlags |= (1 << 1)
}

func (b *Board) UnsetWhiteOOO() {
	b.castlingFlags &= ^uint8(1 << 1)
}

func (b *Board) CanWhiteOOO() bool {
	return b.castlingFlags&(1<<1) > 0
}

func (b *Board) SetBlackOO() {
	b.castlingFlags |= (1 << 2)
}

func (b *Board) UnsetBlackOO() {
	b.castlingFlags &= ^uint8(1 << 2)
}

func (b *Board) CanBlackOO() bool {
	return b.castlingFlags&(1<<2) > 0
}

func (b *Board) SetBlackOOO() {
	b.castlingFlags |= (1 << 3)
}

func (b *Board) UnsetBlackOOO() {
	b.castlingFlags &= ^uint8(1 << 3)
}

func (b *Board) CanBlackOOO() bool {
	return b.castlingFlags&(1<<3) > 0
}
