package main

import (
	"fmt"
	"os"
	"tieu/learn/tank/game"
	"tieu/learn/tank/render"
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

func (c *Client) join() {
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
	render := render.NewRender()

	g := game.NewGame(20, 20)
	c := NewClient(g)

	go g.ListenKeys(render.Screen)
	c.join()

	for {
		now := time.Now()

		c.sendStateCh <- g.GetSyncState()

		render.ClearScreen()
		if g.Dead {
			render.DrawEnd(g)
			render.ShowScreen()
			break
		}

		if g.Quit {
			os.Exit(0)
			break
		}

		render.DrawBackground(g)
		g.Tick()

		render.DrawTanks(g)
		render.DrawBullets(g)
		render.DrawScores(g)
		render.ShowScreen()

		waitForFrame(now)
	}

	c.leave()
}

func waitForFrame(startTime time.Time) {
	elapsed := time.Since(startTime)
	if elapsed < game.FrameTime {
		time.Sleep(game.FrameTime - elapsed)
	}
}
