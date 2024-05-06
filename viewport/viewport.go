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

    // Check tank is outside viewport
    if myTank.Pos.X > v.X + v.Width {
        v.X += v.Width
    } else if myTank.Pos.X < v.X {
        v.X -= v.Width
    }

    if myTank.Pos.Y > v.Y + v.Height {
        v.Y += v.Height
    } else if myTank.Pos.Y < v.Y {
        v.Y -= v.Height
    }
}

// translate Game X, Y to view port X, Y
func (v *ViewPort) Translate(x, y int) (int, int) {
	viewX, viewY := x-v.X, y-v.Y
	return viewX, viewY
}
