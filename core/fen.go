package core

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	s "strings"
	u "unicode"
)

func BoardFromFen(fen string) (Board, error) {
	var board Board
	board.castlingFlags = 0
	board.epPossible = false
	board.activeColor = White
	board.halfMoveClock = 0
	board.fullMoveClock = 0

	fenParts := s.Split(fen, " ")
	piecesPart := fenParts[0]
	for i, row := range s.Split(piecesPart, "/") {
		j := 0
		for _, c := range row {
			sq := square((7-i)*8 + j)
			if u.IsDigit(c) {
				j += int(c - '0')
			} else {
				p, err := CharToPiece(c)
				if err != nil {
					return board, errors.New("Error while parsing fen.")
				}
				board.bitBoards[p] = board.bitBoards[p].Set(sq)
				j += 1
			}
		}
	}

	if len(fenParts) > 1 {
		activeColor := s.ToLower(fenParts[1])
		switch activeColor {
		case "w":
			board.activeColor = White
		case "b":
			board.activeColor = Black
		default:
			return board, errors.New(fmt.Sprintf("Invalid active color field: %s", activeColor))
		}
	}

	if len(fenParts) > 2 {
		castlingRights := fenParts[2]
		if castlingRights != "-" {
			for _, c := range castlingRights {
				switch c {
				case 'K':
					board.SetWhiteOO()
				case 'Q':
					board.SetWhiteOOO()
				case 'k':
					board.SetBlackOO()
				case 'q':
					board.SetBlackOOO()
				}
			}
		}
	}

	if len(fenParts) > 3 {
		epTarget := fenParts[3]
		if epTarget != "-" {
			sq, err := StrToSq(epTarget)
			if err != nil {
				slog.Error(fmt.Sprintf("%s", err))
				return board, errors.New(fmt.Sprintf("Invalid en-passant target: %s", epTarget))
			}
			board.epTarget = sq
			board.epPossible = true
		}
	}

	if len(fenParts) > 4 {
		if fenParts[4] != "-" {
			halfMoveClock, err := strconv.Atoi(fenParts[4])
			if err != nil {
				return board, errors.New(fmt.Sprintf("Invalid halfMoveclock number: %s", halfMoveClock))
			}
			board.halfMoveClock = uint(halfMoveClock)
		}
	}

	if len(fenParts) > 5 {
		if fenParts[5] != "-" {
			fullMoveclock, err := strconv.Atoi(fenParts[5])
			if err != nil {
				return board, errors.New(fmt.Sprintf("Invalid fullMoveclock number: %s", fullMoveclock))
			}
			board.fullMoveClock = uint(fullMoveclock)
		}
	}

	return board, nil
}
