package render

import (
	"tieu/learn/tank/game"

	"github.com/gdamore/tcell"
)

type Render struct {
	Screen tcell.Screen
	styles map[string]tcell.Style
}

func NewRender() *Render {
	styles := map[string]tcell.Style{
		"background": tcell.StyleDefault.Background(tcell.Color17).Foreground(tcell.Color17),
		"my_tank":    tcell.StyleDefault.Background(tcell.Color190).Foreground(tcell.Color190),
        "bullet":     tcell.StyleDefault.Background(tcell.Color122).Foreground(tcell.Color122),
	}

	screen, err := tcell.NewScreen()

	if err != nil {
		panic(err)
	}

	if err := screen.Init(); err != nil {
		panic(err)
	}

	return &Render{Screen: screen, styles: styles}
}

func (r *Render) DrawBackground() {
	r.DrawBox(0, 0, 100, 100, r.styles["background"])
}

func (r *Render) DrawTanks(game *game.Game) {
	r.Screen.SetContent(game.MyTank.Pos.X, game.MyTank.Pos.Y, ' ', nil, r.styles["my_tank"])
}

func (r *Render) DrawBullets(game *game.Game) {
	for _, bullet := range game.Bullets {
		r.Screen.SetContent(bullet.Pos.X, bullet.Pos.Y, ' ', nil, r.styles["bullet"])
	}
}

func (r *Render) DrawBox(x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			r.Screen.SetContent(col, row, ' ', nil, style)
		}
	}
}

func (r *Render) ClearScreen() {
	r.Screen.Clear()
}

func (r *Render) ShowScreen() {
	r.Screen.Show()
}
