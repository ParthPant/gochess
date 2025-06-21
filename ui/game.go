package ui

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ParthPant/gochess/core"
	"github.com/ParthPant/gochess/ui/resources"
)

var pieceImages [12]*ebiten.Image

func init() {
	for i, png := range resources.PiecePngs {
		img, _, err := image.Decode(bytes.NewReader(*png))
		if err != nil {
			log.Fatal(err)
		}
		pieceImages[i] = ebiten.NewImageFromImage(img)
	}
}

type ChessGui struct {
	chess            core.ChessGame
	boardSize        int
	whiteColor       color.Color
	blackColor       color.Color
	highlightColor   color.Color
	pickedPiece      *core.Piece
	pickedSquare     *core.Square
	pickedPieceMoves *core.BitBoard
}

func CreateGui(chess core.ChessGame, boardSize int) ChessGui {
	ebiten.SetWindowSize(boardSize, boardSize)
	ebiten.SetWindowTitle("Go Chess.")
	return ChessGui{
		chess:          chess,
		boardSize:      boardSize,
		whiteColor:     color.RGBA{0xe3, 0xc1, 0x6f, 0xff},
		blackColor:     color.RGBA{0xb8, 0x8b, 0x4a, 0xff},
		highlightColor: color.RGBA{0x3f, 0x7a, 0xd9, 0xff},
		pickedPiece:    nil,
		pickedSquare:   nil,
	}
}

func (g *ChessGui) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p, ok := g.pieceAt(x, y)
		if ok {
			g.pickedPiece = &p
			sq := g.squareAt(x, y)
			g.pickedSquare = &sq
			pickedPieceMoves := g.chess.GetLegalPieceMovesBB(sq)
			g.pickedPieceMoves = &pickedPieceMoves
		}
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if g.pickedPiece != nil {
			x, y := ebiten.CursorPosition()
			to := g.squareAt(x, y)
			if (*g.pickedPiece == core.Pw) && (core.SeventhRank.IsSet(to)) {
				slog.Debug("White Promotion.")
			} else if (*g.pickedPiece == core.Pb) && (core.SecondRank.IsSet(to)) {
				slog.Debug("Black Promotion.")
			}
			g.chess.MakeMove(*g.pickedSquare, to, core.Queen)
		}

		g.pickedPiece = nil
		g.pickedSquare = nil
		g.pickedPieceMoves = nil
	}
	return nil
}

func (g *ChessGui) Draw(screen *ebiten.Image) {
	// Draw checkboard
	checkBoard := ebiten.NewImage(g.boardSize, g.boardSize)
	tileSize := g.boardSize / 8
	square := ebiten.NewImage(tileSize, tileSize)
	hightlightSquare := ebiten.NewImage(tileSize, tileSize)
	hightlightSquare.Fill(g.highlightColor)

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			var color color.Color
			if y%2 == 0 {
				if x%2 == 0 {
					color = g.blackColor
				} else {
					color = g.whiteColor
				}
			} else {
				if x%2 == 0 {
					color = g.whiteColor
				} else {
					color = g.blackColor
				}
			}
			square.Fill(color)

			piece, hasPiece := g.chess.Board().GetAtSq(core.SquareFromXY(x, y))

			if g.pickedPieceMoves != nil {
				if g.pickedPieceMoves.IsSet(core.SquareFromXY(x, y)) {
					op := ebiten.DrawImageOptions{}
					op.ColorScale.ScaleAlpha(0.5)
					square.DrawImage(hightlightSquare, &op)
				}
			}

			if hasPiece {
				if g.pickedSquare == nil || *g.pickedSquare != core.SquareFromXY(x, y) {
					drawPiece(uint8(piece), square)
				}
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*tileSize), float64((7-y)*tileSize))
			checkBoard.DrawImage(square, op)

			square.Clear()
		}
	}

	if g.pickedPiece != nil {
		drawPickedPiece(uint8(*g.pickedPiece), checkBoard)
	}
	screen.DrawImage(checkBoard, &ebiten.DrawImageOptions{})
}

func (g *ChessGui) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (g *ChessGui) GameLoop() {
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *ChessGui) squareAt(x int, y int) core.Square {
	tileSize := g.boardSize / 8
	x, y = x/tileSize, y/tileSize
	x, y = x, 7-y
	return core.SquareFromXY(x, y)
}

func (g *ChessGui) pieceAt(x int, y int) (core.Piece, bool) {
	return g.chess.Board().GetAtSq(g.squareAt(x, y))
}

func drawPickedPiece(piece uint8, image *ebiten.Image) {
	pieceImage := pieceImages[piece]
	op := &ebiten.DrawImageOptions{}
	// piece width and height
	pw, ph := pieceImage.Bounds().Dx(), pieceImage.Bounds().Dy()
	// cursor position
	x, y := ebiten.CursorPosition()
	op.GeoM.Translate(float64(x-pw/2), float64(y-ph/2))
	image.DrawImage(pieceImage, op)
}

func drawPiece(piece uint8, sq *ebiten.Image) {
	pieceImage := pieceImages[piece]
	// piece width and height
	pw, ph := pieceImage.Bounds().Dx(), pieceImage.Bounds().Dy()
	// square width and height
	sw, sh := sq.Bounds().Dx(), sq.Bounds().Dy()
	// translations to center the pieces
	tx, ty := float64(sw-pw)/2.0, float64(sh-ph)/2.0
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(tx, ty)
	sq.DrawImage(pieceImage, op)
}
