package main

import (
	"os"
	"tieu/learn/tank/client"
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
	"tieu/learn/tank/viewport"
	"time"
)

func main() {
	gameState := game.NewGame(100, 30)

	client := client.NewClient(gameState)
	err := client.Join()
	if err != nil {
		panic(err)
	}
	defer client.Leave()

	drawler := render.NewRender()

	viewPort := viewport.NewViewPort(drawler.Screen)

	go gameState.ListenKeys(drawler.Screen)

	for {
		drawler.ClearScreen()
		now := time.Now()

		client.SendState()
        client.UpdateState()

		if gameState.Dead {
			drawler.DrawEnd(gameState)
			drawler.ShowScreen()
			break
		}

		if gameState.Quit {
			drawler.Screen.Fini()
			os.Exit(0)
			break
		}

		gameState.Tick()
		viewPort.Move(gameState)

		drawler.DrawBackground(gameState, viewPort)
		drawler.DrawTanks(gameState, viewPort)
		drawler.DrawBullets(gameState, viewPort)
		drawler.DrawScores(gameState)
		drawler.ShowScreen()

		waitForFrame(now)
	}
}

func waitForFrame(startTime time.Time) {
	elapsed := time.Since(startTime)
	if elapsed < game.FrameTime {
		time.Sleep(game.FrameTime - elapsed)
	}
}
