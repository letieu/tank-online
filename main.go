package main

import (
	"fmt"
	"os"
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
	"tieu/learn/tank/viewport"
	"time"

	"github.com/go-redis/redis"
)

type Client struct {
	game        *game.Game
	redisClient *redis.Client
	sendStateCh chan game.SyncState
}

func NewClient(g *game.Game) *Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "secret",
	})

	sendStateCh := make(chan game.SyncState)

	return &Client{
		game:        g,
		redisClient: redisClient,
		sendStateCh: sendStateCh,
	}
}

func (c *Client) join() error {
	pubsub := c.redisClient.Subscribe("default")
	ch := pubsub.Channel()

	go func() {
		for {
			e := c.redisClient.Publish("default", <-c.sendStateCh).Err()
			if e != nil {
				fmt.Println(e)
			}
		}
	}()

	go func() {
		for msg := range ch {
			var state game.SyncState
			state.UnmarshalBinary([]byte(msg.Payload))

			if state.Id == c.game.MyTank {
				continue
			}

			c.game.HandleRemoteState(state)
		}
	}()

    return c.redisClient.Ping().Err()
}

func (c *Client) leave() {
	err := c.redisClient.Publish("default", game.SyncState{
		Dead: true,
	}).Err()

	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	gameState := game.NewGame(40, 20)

	client := NewClient(gameState)
    err := client.join()
    if err != nil {
        panic(err)
    }

	drawler := render.NewRender()
	go gameState.ListenKeys(drawler.Screen)

	viewPort := viewport.ViewPort{
		Width:  50,
		Height: 30,
	}

	for {
		now := time.Now()

		client.sendStateCh <- gameState.GetSyncState()

		if gameState.Dead {
			drawler.DrawEnd(gameState)
			drawler.ShowScreen()
			break
		}

		if gameState.Quit {
			os.Exit(0)
			break
		}

		gameState.Tick()
		viewPort.Move(gameState)

		drawler.ClearScreen()
		drawler.DrawBackground(gameState, &viewPort)
		drawler.DrawTanks(gameState, &viewPort)
		drawler.DrawBullets(gameState, &viewPort)
		drawler.DrawScores(gameState)
		drawler.ShowScreen()

		waitForFrame(now)
	}

	client.leave()
}

func waitForFrame(startTime time.Time) {
	elapsed := time.Since(startTime)
	if elapsed < game.FrameTime {
		time.Sleep(game.FrameTime - elapsed)
	}
}
