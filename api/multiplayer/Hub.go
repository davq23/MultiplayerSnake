package multiplayer

import (
	"davidmultiplayersnake/api/models"
	"davidmultiplayersnake/utils"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Hub center of game operations
type Hub struct {
	Name       string              `json:"name"`
	PlayerNum  int64               `json:"player_num"`
	Register   chan *Client        `json:"-"`
	Move       chan models.Message `json:"-"`
	Tracking   chan models.Message `json:"-"`
	Unregister chan *Client        `json:"-"`
	clients    map[*Client]bool
	players    map[string]*Client
	logger     *utils.Logger
	lock       sync.Mutex
	manager    *HubManager
}

// NewHub returns a *Hub
func NewHub(manager *HubManager, logger *utils.Logger) *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Move:       make(chan models.Message),
		Tracking:   make(chan models.Message),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		players:    make(map[string]*Client),
		logger:     logger,
		manager:    manager,
	}
}

func (h *Hub) checkCollision(client *Client, msg *models.Message) {

	if client.Player.ID != msg.Player.ID && !h.players[msg.Player.ID].close {
		collides, pos := msg.Player.CheckCollision(client.Player)

		if collides {
			if pos != 0 {
				lost := msg.Player.Positions[pos:]

				client.Player.Positions = append(client.Player.Positions, lost...)
				msg.Player.Positions = msg.Player.Positions[:pos]
				client.Player.Score += len(lost)

				if len(msg.Player.Positions) == 1 {
					msg.Type = models.MessageUnregister
					h.players[msg.Player.ID].close = true
				}
			} else {
				switch msg.Player.Direction {
				case models.DirectionDown:
					msg.Player.Direction = models.DirectionUp
				case models.DirectionUp:
					msg.Player.Direction = models.DirectionDown
				case models.DirectionLeft:
					msg.Player.Direction = models.DirectionRight
				case models.DirectionRight:
					msg.Player.Direction = models.DirectionLeft
				}

				msg.Player.Move()
			}
		}
	}

}

func (h *Hub) broadcast(msg models.Message) {
	for client := range h.clients {
		if len(client.Send) < maxMessages {
			select {
			case client.Send <- msg:
			default:
				h.logger.LogChan <- "1"
				close(client.Send)
				client.close = true
			}
		}
	}
}

// Run starts the hub
func (h *Hub) Run() {

	defer func() {
		h.manager.End <- h
	}()

	h.logger.LogChan <- "Starting hub " + h.Name

	tickerDead := time.NewTicker(time.Duration(5 * time.Minute))

	for {
		select {
		case msg := <-h.Tracking:
			msg.Type = models.MessageMove

			for client := range h.clients {
				h.checkCollision(client, &msg)

				if len(client.Send) < 98 {
					select {
					case client.Send <- msg:
					default:
						h.logger.LogChan <- "1"
						close(client.Send)
						client.close = true
					}
				}
			}
		case c := <-h.Register:
			rand.Seed(time.Now().Unix())

			var message models.Message

			h.lock.Lock()

			message.Type = models.MessageRegister

			id, err := uuid.NewRandom()

			if err != nil {
				continue
			}

			c.Player.ID = id.String()

			var r, g, b = utils.GenerateRandomColors(rand.Float64(), 0.99, 0.99)
			c.Player.Color = utils.RGB2Hex(r, g, b)

			message.Player = c.Player

			h.clients[c] = true
			h.players[c.Player.ID] = c

			h.PlayerNum = int64(len(h.players))

			tickerDead.Stop()
			tickerDead = time.NewTicker(time.Duration(5 * time.Minute))
			h.lock.Unlock()

			h.broadcast(message)
		case message := <-h.Move:
			message.Type = models.MessageMove

			h.broadcast(message)

		case c := <-h.Unregister:
			h.lock.Lock()

			delete(h.clients, c)
			delete(h.players, c.Player.ID)

			var message models.Message
			message.Player = c.Player
			message.Type = models.MessageUnregister

			h.PlayerNum = int64(len(h.players))
			tickerDead.Stop()
			tickerDead = time.NewTicker(time.Duration(5 * time.Minute))

			h.lock.Unlock()

			h.broadcast(message)

		case <-tickerDead.C:
			h.lock.Lock()

			h.logger.LogChan <- "Dead ticker"
			h.logger.LogChan <- len(h.players)

			if len(h.players) < 2 {
				h.logger.LogChan <- "Finishing hub " + h.Name

				for k, client := range h.players {
					client.connection.Close()
					delete(h.players, k)
					delete(h.clients, client)
				}

				return

			}
			h.lock.Unlock()

		}
	}

}
