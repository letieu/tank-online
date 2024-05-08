package render

import (
	"github.com/gdamore/tcell/v2"
	"tieu/learn/tank/game"
	"tieu/learn/tank/viewport"
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
	Width  int
	Height int
}

func NewRender() *Render {
	background := tcell.Color18

	styles := map[string]tcell.Style{
		"background": tcell.StyleDefault.Background(background).Foreground(tcell.Color17),
		"my_tank":    tcell.StyleDefault.Background(background).Foreground(tcell.Color190),
		"enemy_tank": tcell.StyleDefault.Background(background).Foreground(tcell.Color50),
		"bullet":     tcell.StyleDefault.Background(background).Foreground(tcell.Color122),
		"score":      tcell.StyleDefault.Background(background).Foreground(tcell.ColorRed),
		"view_port":  tcell.StyleDefault.Background(tcell.Color88).Foreground(tcell.Color106),
	}

	screen, err := tcell.NewScreen()
	screen.Clear()

	if err != nil {
		panic(err)
	}

	if err := screen.Init(); err != nil {
		panic(err)
	}

	width, height := screen.Size()

	return &Render{Screen: screen, styles: styles, Width: width, Height: height}
}

func (r *Render) DrawBackground(g *game.Game, vp *viewport.ViewPort) {
	x1, y1 := vp.Translate(0, 0)
	x2, y2 := vp.Translate(g.Width, g.Height)

	if x2 > vp.Width {
		x2 = vp.Width
	}

	if y2 > vp.Height {
		y2 = vp.Height
	}

	r.DrawBox(0, 0, vp.Width, vp.Height, r.styles["view_port"])
	r.DrawBox(x1, y1, x2, y2, r.styles["background"])
}

func (r *Render) DrawTanks(g *game.Game, vp *viewport.ViewPort) {
	for id, tank := range g.Tanks {
		if id == g.MyTank {
			r.drawTank(tank, "my_tank", vp)
		} else {
			r.drawTank(tank, "enemy_tank", vp)
		}
	}
}

func (r *Render) DrawScores(g *game.Game) {
}

func (r *Render) DrawEnd(g *game.Game) {
	r.DrawText(r.Width/2-5, r.Height/2, "You are dead!", r.styles["score"])
}

func (r *Render) drawTank(t *game.Tank, style string, vp *viewport.ViewPort) {
	var tankSprite [3][3]rune = tankSprites[t.Direction]

	x, y := vp.Translate(t.Pos.X, t.Pos.Y)

	if x > vp.Width || y > vp.Height {
		return
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r.Screen.SetContent(x+i-1, y+j-1, tankSprite[i][j], nil, r.styles[style])
		}
	}

}

func (r *Render) DrawBullets(g *game.Game, vp *viewport.ViewPort) {
	for _, bullet := range g.Bullets {
		x, y := vp.Translate(bullet.Pos.X, bullet.Pos.Y)
		if x > vp.Width || y > vp.Height {
			continue
		}
		r.Screen.SetContent(x, y, '⬤', nil, r.styles["bullet"])
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
			r.Screen.SetContent(col, row, '-', nil, style)
		}
	}
}

func (r *Render) DrawText(x, y int, text string, style tcell.Style) {
	for i, c := range text {
		r.Screen.SetContent(x+i, y, c, nil, style)
	}
}

func (r *Render) ClearScreen() {
	r.Screen.Clear()
}

func (r *Render) ShowScreen() {
	r.Screen.Show()
}
