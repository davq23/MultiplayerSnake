package multiplayer

import (
	"davidmultiplayersnake/api/models"
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
			msg.Player.ID = c.Player.ID
			c.Player.Direction = msg.Player.Direction

			msg.Player = c.Player

			c.hub.Move <- msg

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

		case models.MessageTracking:
			msg.Player = c.Player

			if len(c.Send) < maxMessages {
				c.hub.Tracking <- msg
			}
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

	tickerPing := time.NewTicker(time.Duration(100 * time.Millisecond))

	//frames := 0

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

				if !c.close {
					if err := c.connection.WriteJSON(msg); err != nil {
						c.hub.logger.LogChan <- err

						return
					}
				}

			}

		}

	}

}
