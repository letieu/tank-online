package game

type Pos struct {
	X, Y int
}

type Tank struct {
	Pos       *Pos
	Direction int
}

func (t *Tank) move() {
	t.Pos.X++
}

type Game struct {
	MyTank *Tank
	Width  int
	Height int
}

func (g *Game) Tick() {
	g.MyTank.move()
}

func NewGame() Game {
	myTank := &Tank{Pos: &Pos{X: 5, Y: 5}}
	return Game{MyTank: myTank}
}
