package models

import (
	"davidmultiplayersnake/api/config"
)

const playerSpeed = 40
const playerRadius = 4

// PlayerDiameter diameter of each player segment
const PlayerDiameter = playerRadius * 2
const minPlayerPositions = 4

const distanceAdjustment = 64

//const distanceAdjustment = 16

// Player represents a player in the game
type Player struct {
	ID                string `json:"player_id"`
	Name              string `json:"name"`
	PlayerLength      int
	PlayerTotalLength int
	Color             string     `json:"color"`
	Positions         []Position `json:"positions"`
	Direction         Direction  `json:"direction"`
	Score             int        `json:"score,omitempty"`
}

// NewPlayer returns a *Player
func NewPlayer(name string, head Position, direction Direction) (p *Player) {
	positions := make([]Position, minPlayerPositions, 8)

	p = &Player{
		Positions:         positions,
		Direction:         direction,
		PlayerLength:      minPlayerPositions,
		PlayerTotalLength: minPlayerPositions * PlayerDiameter,
		Name:              name,
	}

	p.ChangeHead(head)

	return
}

// Move controls player movement
func (p *Player) Move() {
	for i := len(p.Positions) - 1; i > 0; i-- {
		p.Positions[i].X = p.Positions[i-1].X
		p.Positions[i].Y = p.Positions[i-1].Y
	}

	switch p.Direction {
	case DirectionDown:
		p.Positions[0].Y += playerSpeed
	case DirectionLeft:
		p.Positions[0].X -= playerSpeed
	case DirectionUp:
		p.Positions[0].Y -= playerSpeed
	case DirectionRight:
		p.Positions[0].X += playerSpeed
	}

	if p.Positions[0].X > config.PlayfieldWidth {
		p.Positions[0].X = config.MinPlayfieldWidth
	} else if p.Positions[0].X < config.MinPlayfieldWidth {
		p.Positions[0].X = config.PlayfieldWidth
	}

	if p.Positions[0].Y > config.PlayfieldHeight {
		p.Positions[0].Y = config.MinPlayfieldHeight
	} else if p.Positions[0].Y < config.MinPlayfieldHeight {
		p.Positions[0].Y = config.PlayfieldHeight
	}

}

// CheckCollision checks whether the head of another player has collided with the current player and also at which position
func (p *Player) CheckCollision(player *Player) (bool, int) {
	if distance(p.Positions[0], player.Positions[0]) < p.PlayerTotalLength+distanceAdjustment {
		for i := 0; i < p.PlayerLength; i++ {
			dist := distance(p.Positions[i], player.Positions[0])

			if dist <= PlayerDiameter+distanceAdjustment {
				return true, i
			}

		}
	}
	return false, -1
}

// ChangeHead changes the position of the head segment and rearranges the rest according to player's direction
func (p *Player) ChangeHead(head Position) {
	p.Positions[0] = head

	if p.Direction == DirectionLeft || p.Direction == DirectionRight {
		for i := len(p.Positions) - 1; i > 0; i-- {
			p.Positions[i].X = p.Positions[0].X + playerRadius*2*i + playerSpeed
			p.Positions[i].Y = p.Positions[0].Y
		}
	} else {
		for i := len(p.Positions) - 1; i > 0; i-- {
			p.Positions[i].X = p.Positions[0].X
			p.Positions[i].Y = p.Positions[0].Y + playerRadius*2*i + playerSpeed
		}
	}
}
