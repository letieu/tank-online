package main

import (
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
	"time"
)

const frameRate = 60
const frameTime = time.Second / frameRate

func main() {
	render := render.NewRender()
	windowWidth, windowHeight := render.Screen.Size()

	game := game.NewGame(windowWidth, windowHeight)

	go game.ListenKeys(render.Screen)

	for {
		now := time.Now()

		render.ClearScreen()
		render.DrawBackground()
		game.Tick()
		render.DrawTanks(&game)
		render.DrawBullets(&game)
		render.ShowScreen()

		waitForFrame(now)
	}
}

func waitForFrame(startTime time.Time) {
	elapsed := time.Since(startTime)
	if elapsed < frameTime {
		time.Sleep(frameTime - elapsed)
	}
}
