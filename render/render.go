package render

import (
	"tieu/learn/tank/game"

	"github.com/gdamore/tcell"
)

var tankSprites = map[int][3][3]rune{
	game.Up: {
		{' ', '🀫', '🀫'},
		{'🀫', '🀫', '🀫'},
		{' ', '🀫', '🀫'},
	},
	game.Down: {
		{'🀫', '🀫', ' '},
		{'🀫', '🀫', '🀫'},
		{'🀫', '🀫', ' '},
	},
	game.Left: {
		{' ', '🀫', ' '},
		{'🀫', '🀫', '🀫'},
		{'🀫', '🀫', '🀫'},
	},
	game.Right: {
		{'🀫', '🀫', '🀫'},
		{'🀫', '🀫', '🀫'},
		{' ', '🀫', ' '},
	},
}

type Render struct {
	Screen tcell.Screen
	styles map[string]tcell.Style
}

func NewRender() *Render {
	styles := map[string]tcell.Style{
		"background": tcell.StyleDefault.Background(tcell.Color17).Foreground(tcell.Color17),
		"my_tank":    tcell.StyleDefault.Background(tcell.Color17).Foreground(tcell.Color190),
		"bullet":     tcell.StyleDefault.Background(tcell.Color17).Foreground(tcell.Color122),
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

func (r *Render) DrawTanks(g *game.Game) {
	r.drawTank(g.MyTank)

	for _, tank := range g.EnemyTanks {
		r.drawTank(tank)
	}
}

func (r *Render) drawTank(t *game.Tank) {
	var tankSprite [3][3]rune = tankSprites[t.Direction]

	x, y := t.Pos.X, t.Pos.Y
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r.Screen.SetContent(x+i-1, y+j-1, tankSprite[i][j], nil, r.styles["my_tank"])
		}
	}

}

func (r *Render) DrawBullets(g *game.Game) {
	for _, bullet := range g.Bullets {
		r.Screen.SetContent(bullet.Pos.X, bullet.Pos.Y, '⬤', nil, r.styles["bullet"])
	}

	for _, bullet := range g.EnemyBullets {
		r.Screen.SetContent(bullet.Pos.X, bullet.Pos.Y, '⬤', nil, r.styles["bullet"])
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
