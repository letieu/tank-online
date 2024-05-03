package render

import (
	"tieu/learn/tank/game"

	"github.com/gdamore/tcell"
)

type Render struct {
	screen tcell.Screen
	styles map[string]tcell.Style
}

func NewRender() *Render {
	styles := map[string]tcell.Style{
		"background": tcell.StyleDefault.Background(tcell.Color17).Foreground(tcell.Color17),
		"my_tank":    tcell.StyleDefault.Background(tcell.Color190).Foreground(tcell.Color190),
	}

	screen, err := tcell.NewScreen()

	if err != nil {
		panic(err)
	}

	if err := screen.Init(); err != nil {
		panic(err)
	}

	return &Render{screen: screen, styles: styles}
}

func (r *Render) DrawBackground() {
	r.DrawBox(0, 0, 100, 100, r.styles["background"])
}

func (r *Render) DrawTanks(game *game.Game) {
	r.screen.SetContent(game.MyTank.Pos.X, game.MyTank.Pos.Y, ' ', nil, r.styles["my_tank"])
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
			r.screen.SetContent(col, row, ' ', nil, style)
		}
	}
}

func (r *Render) ClearScreen() {
	r.screen.Clear()
}

func (r *Render) ShowScreen() {
	r.screen.Show()
}
