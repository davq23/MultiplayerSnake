package multiplayer

import (
	"davidmultiplayersnake/api/config"
	"davidmultiplayersnake/api/models"
	"davidmultiplayersnake/api/security"
	"time"

	"github.com/gorilla/websocket"
)

// Client multiplayer game client
type Client struct {
	connection *websocket.Conn
	hub        *Hub
	Player     *models.Player
	Send       chan models.Message
	close      bool
}

// NewClient returns a *Client
func NewClient(connection *websocket.Conn, hub *Hub, player *models.Player) *Client {
	return &Client{
		connection: connection,
		hub:        hub,
		Player:     player,
		Send:       make(chan models.Message, maxMessages),
		close:      false,
	}
}

// ReadPump begins concurrent reads
func (c *Client) ReadPump() {
	defer func() {
		if c.Player != nil {
			c.hub.logger.LogChan <- "Finishing read goroutine of " + c.Player.ID

		}

		c.hub.Unregister <- c
		c.connection.Close()
		c.close = true
	}()

	for {
		var msg models.Message
		var err error

		if !c.close {
			err = c.connection.ReadJSON(&msg)

			if err != nil {
				c.hub.logger.LogChan <- err

				return
			}
		} else {
			return
		}

		switch msg.Type {
		case models.MessageMove:
			//msg.Player.ID = c.Player.ID
			c.Player.Direction = msg.Player.Direction

			/*msg.Player = c.Player

			c.hub.Move <- msg*/

		case models.MessageRegister:

			c.hub.Register <- c

		case models.MessageUnregister:

			c.hub.Unregister <- c

		case models.MessageGetPlayers:
			players := make(map[string]models.Player)

			for pid := range c.hub.players {
				if pid != c.Player.ID {
					players[pid] = *c.hub.players[pid].Player
				}
			}

			msg.Players = &players

			c.Send <- msg
			/*
				case models.MessageTracking:
					msg.Player = c.Player*/

		}

	}
}

// WritePump begins concurrent writes
func (c *Client) WritePump() {
	defer func() {
		c.hub.Unregister <- c
		if c.Player != nil {
			c.hub.logger.LogChan <- "Finishing write goroutine of " + c.Player.ID
		}
		c.connection.Close()
		c.close = true
	}()

	tickerPing := time.NewTicker(time.Duration(trackingTick) * time.Millisecond)
	tickerRefresh := time.NewTicker(time.Duration(1) * time.Minute)

	c.connection.SetPongHandler(func(string) error {
		c.connection.SetWriteDeadline(time.Now().Add(time.Duration(time.Minute)))
		return nil
	})

	for {

		select {
		case message, ok := <-c.Send:
			if ok {
				message.ReceivedAt = time.Now()

				if !c.close {
					c.connection.WriteJSON(message)

				}

			} else {
				c.hub.logger.LogChan <- message
				return
			}
		case <-tickerPing.C:
			if c.Player != nil {
				if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}

				var msg models.Message
				msg.Type = models.MessageTracking
				msg.Player = c.Player

				if len(c.Send) < maxMessages {
					c.hub.Tracking <- msg
				}

			}
		case <-tickerRefresh.C:
			if c.Player != nil {
				var msg models.Message
				msg.Type = models.MessageRefresh

				newTokenChannel := make(chan string)

				go func(ch chan<- string) {
					str, err := security.GetToken(c.Player.Name, c.hub.Name, config.JWTSecret, c.Player.Score, int64(time.Duration(time.Minute*15)))

					if err != nil {
						str = ""
					}

					ch <- str
				}(newTokenChannel)

				msg.NewToken = <-newTokenChannel

				if msg.NewToken != "" {
					c.connection.WriteJSON(msg)
				}

			}
		}

	}

}
