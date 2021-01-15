package models

import (
	"math"
)

// Position is in-game object position
type Position struct {
	X int
	Y int
}

func distance(p1 Position, p2 Position) int {
	dx := math.Abs(float64(p1.X - p2.X))
	dy := math.Abs(float64(p1.Y - p2.Y))

	return int(math.Round(0.394*(dx+dy) + 0.554*math.Max(dx, dy)))
}

/*
func euclideanDistance(p1 Position, p2 Position) int {
	return int(math.)
}
*/
