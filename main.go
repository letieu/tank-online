package main

import (
	"flag"
	"tieu/learn/tank/client"
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
	"tieu/learn/tank/viewport"
	"time"
)

var name string

func main() {
	flag.StringVar(&name, "name", "Tank", "Name of your tank")
    flag.Parse()

	gameState := game.NewGame(100, 30, name)

	client := client.NewClient(gameState)
	err := client.Join()
	if err != nil {
		panic(err)
	}
	defer client.Leave()

	drawler := render.NewRender()

	screenW, screenH := drawler.Screen.Size()
	viewPort := viewport.NewViewPort(screenW, screenH)

	go gameState.ListenKeys(drawler.Screen)

	for {
		drawler.ClearScreen()
		now := time.Now()

		client.SendState()
		client.UpdateState()

		if gameState.Dead || gameState.Quit {
			break
		}

		gameState.Tick()
		viewPort.Move(gameState)
		drawler.DrawGame(gameState, viewPort)

		waitForFrame(now)
		drawler.ShowScreen()
	}

	drawler.DrawEnd(gameState)
	drawler.ShowScreen()
}

func waitForFrame(startTime time.Time) {
	elapsed := time.Since(startTime)
	if elapsed < game.FrameTime {
		time.Sleep(game.FrameTime - elapsed)
	}
}
