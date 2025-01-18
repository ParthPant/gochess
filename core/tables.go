package core

import (
	"errors"
	"log/slog"
	"math/bits"
	"math/rand"
	"sync"
)

var PawnAtkTable [2][64]BitBoard
var KnightAtkTable [64]BitBoard
var KingAtkTable [64]BitBoard

type MagicEntry struct {
	mask      BitBoard
	magic     uint64
	indexBits uint8
}

type relevantOccupancyFunc func(int) BitBoard
type attackFunc func(int, BitBoard) BitBoard

var RookMagics [64]MagicEntry
var BishopMagics [64]MagicEntry

var BishopMoves [64][]BitBoard
var RookMoves [64][]BitBoard

func init() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		computePawnAtkTable()
		computeKnightAtkTable()
		computeKingAtkTable()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		computeMagics(rookRelevantOccupancy, rookAttack, &RookMagics, &RookMoves)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		computeMagics(bishopRelevantOccupancy, bishopAttack, &BishopMagics, &BishopMoves)
	}()
	wg.Wait()
	slog.Info("Magic Tables have been constructed.")
}

func GetBishopMoves(sq square, blockers BitBoard) BitBoard {
	magic := BishopMagics[sq]
	moves := BishopMoves[sq]
	return moves[magic.magic_index(blockers)]
}

func GetRookMoves(sq square, blockers BitBoard) BitBoard {
	magic := RookMagics[sq]
	moves := RookMoves[sq]
	return moves[magic.magic_index(blockers)]
}

func computePawnAtkTable() {
	// << up ; >> down
	for i := 0; i < 64; i++ {
		var whitePawn BitBoard = 1 << i
		whitePawn = (whitePawn<<7) & ^HFile | (whitePawn<<9) & ^AFile
		PawnAtkTable[White][i] = whitePawn

		var blackPawn BitBoard = 1 << i
		blackPawn = (blackPawn>>7) & ^AFile | (blackPawn>>9) & ^HFile
		PawnAtkTable[Black][i] = blackPawn
	}
}

func computeKnightAtkTable() {
	// << up ; >> down
	for i := 0; i < 64; i++ {
		var sq BitBoard = 1 << i
		var knight BitBoard = 0
		knight |= sq << 6
		knight |= sq << 15
		knight |= sq << 17
		knight |= sq << 10

		knight |= sq >> 6
		knight |= sq >> 15
		knight |= sq >> 17
		knight |= sq >> 10

		if sq&AFile > 0 || sq&BFile > 0 {
			knight &= ^(GFile | HFile)
		}

		if sq&GFile > 0 || sq&HFile > 0 {
			knight &= ^(AFile | BFile)
		}

		KnightAtkTable[i] = knight
	}
}

func computeKingAtkTable() {
	// << up ; >> down
	for i := 0; i < 64; i++ {
		var sq BitBoard = 1 << i
		king := (sq << 8) | (sq >> 8) | (sq << 1) | (sq >> 1) | (sq << 9) | (sq >> 9) | (sq << 7) | (sq >> 7)

		if sq&AFile > 0 {
			king &= ^HFile
		}

		if sq&HFile > 0 {
			king &= ^AFile
		}

		KingAtkTable[i] = king
	}
}

func bishopRelevantOccupancy(sq int) BitBoard {
	var occ BitBoard

	// file and rank
	f, r := sq%8+1, sq/8+1
	for {
		if f > 6 || r > 6 {
			break
		}
		occ |= 1 << (r*8 + f)
		f++
		r++
	}

	f, r = sq%8-1, sq/8+1
	for {
		if f < 1 || r > 6 {
			break
		}
		occ |= 1 << (r*8 + f)
		f--
		r++
	}

	f, r = sq%8-1, sq/8-1
	for {
		if f < 1 || r < 1 {
			break
		}
		occ |= 1 << (r*8 + f)
		f--
		r--
	}

	f, r = sq%8+1, sq/8-1
	for {
		if f > 6 || r < 1 {
			break
		}
		occ |= 1 << (r*8 + f)
		f++
		r--
	}
	return occ
}

func bishopAttack(sq int, blockers BitBoard) BitBoard {
	var atk BitBoard

	// file and rank
	f, r := sq%8+1, sq/8+1
	for {
		if f > 6 || r > 6 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		f++
		r++
	}

	f, r = sq%8-1, sq/8+1
	for {
		if f < 1 || r > 6 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		f--
		r++
	}

	f, r = sq%8-1, sq/8-1
	for {
		if f < 1 || r < 1 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		f--
		r--
	}

	f, r = sq%8+1, sq/8-1
	for {
		if f > 6 || r < 1 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		f++
		r--
	}

	return atk
}

func rookRelevantOccupancy(sq int) BitBoard {
	var occ BitBoard

	// file and rank
	f, r := sq%8+1, sq/8
	for {
		if f > 6 {
			break
		}
		occ |= 1 << (r*8 + f)
		f++
	}

	f, r = sq%8, sq/8+1
	for {
		if r > 6 {
			break
		}
		occ |= 1 << (r*8 + f)
		r++
	}

	f, r = sq%8-1, sq/8
	for {
		if f < 1 {
			break
		}
		occ |= 1 << (r*8 + f)
		f--
	}

	f, r = sq%8, sq/8-1
	for {
		if r < 1 {
			break
		}
		occ |= 1 << (r*8 + f)
		r--
	}

	return occ
}

func rookAttack(sq int, blockers BitBoard) BitBoard {
	var atk BitBoard

	// file and rank
	f, r := sq%8+1, sq/8
	for {
		if f > 6 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		f++
	}

	f, r = sq%8, sq/8+1
	for {
		if r > 6 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		r++
	}

	f, r = sq%8-1, sq/8
	for {
		if f < 1 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		f--
	}

	f, r = sq%8, sq/8-1
	for {
		if r < 1 {
			break
		}
		atk |= 1 << (r*8 + f)
		if blockers.IsSet(square(r*8 + f)) {
			break
		}
		r--
	}

	return atk
}

func (e *MagicEntry) magic_index(blockers BitBoard) uint64 {
	blockers = blockers & e.mask
	hash := uint64(blockers) * e.magic
	index := hash >> (64 - e.indexBits)
	return index
}

func subsets(set BitBoard) chan BitBoard {
	var subset BitBoard
	ch := make(chan BitBoard)
	go func() {
		defer close(ch)
		for {
			ch <- subset
			subset = (subset - set) & set
			if subset == 0 {
				break
			}
		}
	}()
	return ch
}

func computeMagics(relevantOccupancyFn relevantOccupancyFunc, attackFn attackFunc, magicTable *[64]MagicEntry, movesTable *[64][]BitBoard) {
	for i := 0; i < 64; i++ {
		set := relevantOccupancyFn(i)
		indexBits := uint8(bits.OnesCount64(uint64(set)))
		for {
			magic := rand.Uint64() & rand.Uint64() & rand.Uint64()
			magicEntry := MagicEntry{mask: set, magic: magic, indexBits: indexBits}
			table, err := tryMakeTable(attackFn, magicEntry, square(i))
			if err == nil {
				(*magicTable)[i] = magicEntry
				(*movesTable)[i] = table
				break
			}
		}
	}
}

func tryMakeTable(attackFn attackFunc, m MagicEntry, sq square) ([]BitBoard, error) {
	table := make([]BitBoard, 1<<m.indexBits)
	for blockers := range subsets(m.mask) {
		moves := attackFn(int(sq), blockers)

		tableEntry := &table[m.magic_index(blockers)]
		if *tableEntry == 0 {
			*tableEntry = moves
		} else if *tableEntry != moves {
			return []BitBoard{}, errors.New("Magic entry collision.")
		}
	}
	return table, nil
}
