package main

import (
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "secret",
	})

	render := render.NewRender()
	windowWidth, windowHeight := render.Screen.Size()

	g := game.NewGame(windowWidth, windowHeight)

	sendStateCh := make(chan game.SyncState)

	go g.ListenKeys(render.Screen)

	go func() {
		for {
			sendGameState(<-sendStateCh, rdb)
		}
	}()

	for {
		now := time.Now()

		render.ClearScreen()
		if g.Dead {
			render.DrawEnd(g)
			render.ShowScreen()
			break
		}

		render.DrawBackground()

		g.Tick()
		sendStateCh <- g.GetSyncState()

		render.DrawTanks(g)
		render.DrawBullets(g)
		render.DrawScores(g)
		render.ShowScreen()

		waitForFrame(now)
	}
}

func waitForFrame(startTime time.Time) {
	elapsed := time.Since(startTime)
	if elapsed < game.FrameTime {
		time.Sleep(game.FrameTime - elapsed)
	}
}

func sendGameState(state game.SyncState, redisClient *redis.Client) {
    redisClient.Set("state", state, time.Minute)
}
