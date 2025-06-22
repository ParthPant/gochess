package core

const MinScore int32 = -10000000
const MaxScore int32 = 10000000
const MatingScore int32 = -9999999

var PieceScore = [...]int32{300, 350, 500, 1000, 10000, 100, -300, -350, -500, -1000, -10000, 100}

var PawnPosScore = [...]int32{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, -10, -10, 0, 0, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	5, 5, 10, 20, 20, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 30, 20, 20,
	30, 30, 30, 40, 40, 30, 30, 30,
	90, 90, 90, 90, 90, 90, 90, 90,
}

var KnightScore = [...]int32{
	-5, -10, 0, 0, 0, 0, -10, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
	-5, 5, 20, 10, 10, 20, 5, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 10, 20, 30, 30, 20, 10, -5,
	-5, 5, 20, 20, 20, 20, 5, -5,
	-5, 0, 0, 10, 10, 0, 0, -5,
	-5, 0, 0, 0, 0, 0, 0, -5,
}

var BishopScore = [...]int32{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 30, 0, 0, 0, 0, 30, 0,
	0, 10, 0, 0, 0, 0, 10, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var RookScore = [...]int32{
	0, 0, 0, 20, 20, 0, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	50, 50, 50, 50, 50, 50, 50, 50,
	50, 50, 50, 50, 50, 50, 50, 50,
}

var KingScore = [...]int32{
	0, 0, 5, 0, -15, 0, 10, 0,
	0, 5, 5, -5, -5, 0, 5, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 10, 20, 20, 10, 5, 0,
	0, 5, 5, 10, 10, 5, 5, 0,
	0, 0, 5, 5, 5, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

func evaluateBoard(b *Board) int32 {
	var score int32 = 0
	for square := range 64 {
		if piece, occupied := b.GetAtSq(Square(square)); occupied {
			material_score := PieceScore[piece]

			var position_score int32
			var position Square
			switch piece.GetColor() {
			case White:
				position = Square(square)
			case Black:
				position = MirrorSquare[square]
			}
			switch piece {
			case Pw, Pb:
				position_score = PawnPosScore[position]
			case Nw, Nb:
				position_score = KnightScore[position]
			case Bw, Bb:
				position_score = BishopScore[position]
			case Rw, Rb:
				position_score = RookScore[position]
			case Kw, Kb:
				position_score = KingScore[position]
			default:
				position_score = 0
			}

			switch piece.GetColor() {
			case White:
				score += (material_score + position_score)
			case Black:
				score += (material_score - position_score)
			}
		}
	}
	return score
}
