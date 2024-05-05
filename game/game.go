package game

import (
	"os"
	"strconv"
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

func (t *Tank) move(width, height int) {
	if t.FramesUntilNextMove > 0 {
		t.FramesUntilNextMove--
		return
	}

	t.Pos.move(t.Direction)

	if t.Pos.isOutOfScreen(width, height) {
		switch t.Direction {
		case Up:
			t.Pos.Y = height - 1
		case Down:
			t.Pos.Y = 0
		case Left:
			t.Pos.X = width - 1
		case Right:
			t.Pos.X = 0
		}
	}

	t.FramesUntilNextMove = FrameRate / t.Speed
}

func (t *Tank) fire() Bullet {
	t.Fire = false
	return Bullet{Pos: &Pos{X: t.Pos.X, Y: t.Pos.Y}, Direction: t.Direction, Speed: t.FireSpeed}
}

func (t *Tank) isHit(b *Bullet) bool {
	tankBounds := [3][3]Pos{
		{Pos{X: t.Pos.X - 1, Y: t.Pos.Y - 1}, Pos{X: t.Pos.X, Y: t.Pos.Y - 1}, Pos{X: t.Pos.X + 1, Y: t.Pos.Y - 1}},
		{Pos{X: t.Pos.X - 1, Y: t.Pos.Y}, Pos{X: t.Pos.X, Y: t.Pos.Y}, Pos{X: t.Pos.X + 1, Y: t.Pos.Y}},
		{Pos{X: t.Pos.X - 1, Y: t.Pos.Y + 1}, Pos{X: t.Pos.X, Y: t.Pos.Y + 1}, Pos{X: t.Pos.X + 1, Y: t.Pos.Y + 1}},
	}

	for _, bound := range tankBounds {
		for _, pos := range bound {
			if pos.X == b.Pos.X && pos.Y == b.Pos.Y {
				return true
			}
		}
	}

	return false
}

type Game struct {
	Dead   bool
	Kills  int
	Width  int
	Height int

	Tanks  map[string]*Tank
	MyTank string

	Bullets []Bullet
}

func (g *Game) Tick() {
	g.handleMove()
	g.handleFire()
	g.handleBullets()
}

func (g *Game) ListenKeys(screen tcell.Screen) {
	myTank := g.Tanks[g.MyTank]

	for {
		direction := myTank.Direction
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

		myTank.Direction = direction
		myTank.Fire = fire

		g.Tanks[g.MyTank] = myTank
	}
}

func NewGame(width, height int) *Game {
	myTank := &Tank{Pos: &Pos{X: 5, Y: 5}, Direction: Up}
	myTank.Speed = 30
	myTank.FireSpeed = 40
	myId := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	tanks := map[string]*Tank{
		myId: myTank,
	}

	game := Game{MyTank: myId, Width: width, Height: height, Tanks: tanks}

	return &game
}

func (g *Game) handleMove() {
	for _, tank := range g.Tanks {
		tank.move(g.Width, g.Height)
	}
}

func (g *Game) handleFire() {
	for _, tank := range g.Tanks {
		if tank.Fire {
			g.Bullets = append(g.Bullets, tank.fire())
		}
	}
}

func (g *Game) handleBullets() []Bullet {
	remainBullet := make([]Bullet, 0)
	for _, bullet := range g.Bullets {
		bullet.move()

		if !bullet.Pos.isOutOfScreen(g.Width, g.Height) {
			remainBullet = append(remainBullet, bullet)
		}
	}
	return remainBullet
}

type SyncState struct {
	Pos                 Pos
	Direction           int
	Fire                bool
	Speed               int
	FramesUntilNextMove int
	FireSpeed           int
}

func (g *Game) GetSyncState() SyncState {
	myTank := g.Tanks[g.MyTank]

	syncState := SyncState{
		Pos:                 *myTank.Pos,
		Direction:           myTank.Direction,
		Fire:                myTank.Fire,
		Speed:               myTank.Speed,
		FramesUntilNextMove: myTank.FramesUntilNextMove,
		FireSpeed:           myTank.FireSpeed,
	}

	return syncState
}
