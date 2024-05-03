package game

import (
	"os"
	"time"

	"github.com/gdamore/tcell"
)

const (
	Up    = 1
	Down  = -1
	Left  = 2
	Right = -2
)

const FrameRate = 60
const FrameTime = time.Second / FrameRate

type Pos struct {
	X, Y int
}

func (p *Pos) isOutOfScreen(width, height int) bool {
	return p.X < 0 || p.X >= width || p.Y < 0 || p.Y >= height
}

func (p *Pos) move(direction int) {
	switch direction {
	case Up:
		p.Y--
	case Down:
		p.Y++
	case Left:
		p.X--
	case Right:
		p.X++
	}
}

type Bullet struct {
	Pos                 *Pos
	Direction           int
	Speed               int
	FramesUntilNextMove int
}

func (b *Bullet) move() {
	b.Pos.move(b.Direction)

	if b.FramesUntilNextMove > 0 {
		b.FramesUntilNextMove--
		return
	}
}

type Tank struct {
	Pos                 *Pos
	Direction           int
	Fire                bool
	Speed               int
	FramesUntilNextMove int
	FireSpeed           int
}

func (t *Tank) move() {
	if t.FramesUntilNextMove > 0 {
		t.FramesUntilNextMove--
		return
	}

	t.Pos.move(t.Direction)

	t.FramesUntilNextMove = FrameRate / t.Speed
}

func (t *Tank) fire() *Bullet {
	return &Bullet{Pos: &Pos{X: t.Pos.X, Y: t.Pos.Y}, Direction: t.Direction, Speed: t.FireSpeed}
}

type Game struct {
	MyTank *Tank
	Width  int
	Height int

	Bullets []*Bullet
}

func (g *Game) Tick() {
	g.MyTank.move()

	if g.MyTank.Fire {
		g.Bullets = append(g.Bullets, g.MyTank.fire())
		g.MyTank.Fire = false
	}

	remainBullet := make([]*Bullet, 0)
	for _, bullet := range g.Bullets {
		bullet.move()

		if !bullet.Pos.isOutOfScreen(g.Width, g.Height) {
			remainBullet = append(remainBullet, bullet)
		}
	}

	g.Bullets = remainBullet
}

func (g *Game) ListenKeys(screen tcell.Screen) {
	for {
		direction := g.MyTank.Direction
		fire := false

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			// arrow keys
			switch ev.Key() {
			case tcell.KeyUp:
				direction = Up
			case tcell.KeyDown:
				direction = Down
			case tcell.KeyLeft:
				direction = Left
			case tcell.KeyRight:
				direction = Right

			case tcell.KeyCtrlC:
				screen.Fini()
				os.Exit(0)
				return
			}

			// vim keys
			switch ev.Rune() {
			case 'h':
				direction = Left
			case 'j':
				direction = Down
			case 'k':
				direction = Up
			case 'l':
				direction = Right
			// fire
			case ' ':
				fire = true
			}
		}

		g.MyTank.Direction = direction
		g.MyTank.Fire = fire
	}
}

func NewGame(width, height int) Game {
	myTank := &Tank{Pos: &Pos{X: 5, Y: 5}, Direction: Up}
	myTank.Speed = 30
	myTank.FireSpeed = 40
	return Game{MyTank: myTank, Width: width, Height: height}
}
