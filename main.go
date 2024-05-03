package main

import (
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
	"time"
)

const (
	Up    = 1
	Down  = -1
	Left  = 2
	Right = -2
)

func main() {
	render := render.NewRender()
	game := game.NewGame()

	for {
		render.ClearScreen()
		render.DrawBackground()
		game.Tick()
		render.DrawTanks(&game)
		render.ShowScreen()

		time.Sleep(time.Millisecond * 100) // TODO: if online, this not good
	}
}
