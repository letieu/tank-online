package viewport

import (
	"tieu/learn/tank/game"
)

type ViewPort struct {
	Width, Height int
	game.Pos
}

func (v *ViewPort) Move(g *game.Game) {
    myTank := g.GetMyTank()
    _, _ = v.Translate(myTank.Pos.X, myTank.Pos.Y)
}

func (v *ViewPort) Translate(x, y int) (int, int) {
	viewX, viewY := x-v.X, y-v.Y
	return viewX, viewY
}
