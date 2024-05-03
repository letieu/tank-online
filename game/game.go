package game

import (
	"math/rand"
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

func (t *Tank) fire() *Bullet {
	t.Fire = false
	return &Bullet{Pos: &Pos{X: t.Pos.X, Y: t.Pos.Y}, Direction: t.Direction, Speed: t.FireSpeed}
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
	MyTank  *Tank
	Bullets []*Bullet
	Dead    bool

	Width  int
	Height int

	EnemyTanks   []*Tank
	EnemyBullets []*Bullet
}

func (g *Game) Tick() {
	g.MyTank.move(g.Width, g.Height)
	g.handleMyTankFire()

	remainEnemyTanks := g.handleEnemyTanks()
	g.EnemyTanks = remainEnemyTanks

	remainBullet := g.handleBullets()
	g.Bullets = remainBullet

	remainEnemyBullet := g.handleEnemyBullets()
	g.EnemyBullets = remainEnemyBullet

	if len(g.EnemyTanks) == 0 {
		g.EnemyTanks = append(g.EnemyTanks, g.newEnemyTank())
	}
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

	enemyTank1 := &Tank{Pos: &Pos{X: 10, Y: 10}, Direction: Down}
	enemyTank1.Speed = 10
	enemyTank1.FireSpeed = 30

	enemyTank2 := &Tank{Pos: &Pos{X: 20, Y: 30}, Direction: Right}
	enemyTank2.Speed = 13
	enemyTank2.FireSpeed = 30

	game := Game{MyTank: myTank, Width: width, Height: height}

	game.EnemyTanks = append(game.EnemyTanks, enemyTank1)
	game.EnemyTanks = append(game.EnemyTanks, enemyTank2)

	return game
}

func (g *Game) handleMyTankFire() {
	if g.MyTank.Fire {
		g.Bullets = append(g.Bullets, g.MyTank.fire())
	}
}

func (g *Game) handleEnemyTanks() []*Tank {
	remainEnemyTanks := make([]*Tank, 0)
	for _, enemyTank := range g.EnemyTanks {
		if rand.Intn(5) == 0 {
			enemyTank.Fire = true
		}

		enemyTank.move(g.Width, g.Height)

		if enemyTank.Fire {
			g.EnemyBullets = append(g.EnemyBullets, enemyTank.fire())
		}

		isHit := false
		for _, bullet := range g.Bullets {
			if enemyTank.isHit(bullet) {
				isHit = true
				break
			}
		}

		if !isHit {
			remainEnemyTanks = append(remainEnemyTanks, enemyTank)
		}
	}
	return remainEnemyTanks
}

func (g *Game) handleBullets() []*Bullet {
	remainBullet := make([]*Bullet, 0)
	for _, bullet := range g.Bullets {
		bullet.move()

		if !bullet.Pos.isOutOfScreen(g.Width, g.Height) {
			remainBullet = append(remainBullet, bullet)
		}
	}
	return remainBullet
}

func (g *Game) handleEnemyBullets() []*Bullet {
	remainEnemyBullet := make([]*Bullet, 0)
	for _, bullet := range g.EnemyBullets {
		bullet.move()

		if !bullet.Pos.isOutOfScreen(g.Width, g.Height) {
			remainEnemyBullet = append(remainEnemyBullet, bullet)
		}

		if g.MyTank.isHit(bullet) {
			g.Dead = true
		}
	}
	return remainEnemyBullet
}

func (g *Game) newEnemyTank() *Tank {
	return &Tank{Pos: &Pos{X: 10, Y: 10}, Direction: Right, Speed: 10, FireSpeed: 30}
}