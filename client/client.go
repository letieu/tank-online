package client

import (
	"fmt"
	"tieu/learn/tank/game"

	"github.com/go-redis/redis"
)

type Client struct {
	game          *game.Game
	redisClient   *redis.Client
	sendStateCh   chan game.SyncState
	newTankStates map[string]game.SyncState
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
        newTankStates: make(map[string]game.SyncState),
	}
}

func (c *Client) Join() error {
	pubsub := c.redisClient.Subscribe("default")
	ch := pubsub.Channel()

	go func() {
		for {
			state, ok := <-c.sendStateCh
			if ok == true {
				e := c.redisClient.Publish("default", state).Err()
				if e != nil {
					fmt.Println(e)
				}
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

			c.newTankStates[state.Id] = state
		}
	}()

	return c.redisClient.Ping().Err()
}

func (c *Client) Leave() {
	err := c.redisClient.Publish("default", game.SyncState{
		Dead: true,
	}).Err()

	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) SendState() {
	state := c.game.GetSyncState()
	c.sendStateCh <- state
}

func (c *Client) UpdateState() {
	for _, v := range c.newTankStates {
		c.game.HandleRemoteState(v)
	}
}
