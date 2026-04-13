package game

import "math/rand"

// directions is the set of 8 cardinal and diagonal unit steps.
var directions = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// randomStep returns a randomly chosen direction from the 8 candidates.
func randomStep() (int, int) {
	d := directions[rand.Intn(8)]
	return d[0], d[1]
}

// fleeStep returns the direction (dx, dy) from the 8 candidates that
// maximises the resulting Manhattan distance from (px, py).
func fleeStep(ax, ay, px, py int) (int, int) {
	bestDx, bestDy := 0, 0
	bestDist := -1
	for _, d := range directions {
		nx, ny := ax+d[0], ay+d[1]
		dist := abs(nx-px) + abs(ny-py)
		if dist > bestDist {
			bestDist = dist
			bestDx, bestDy = d[0], d[1]
		}
	}
	return bestDx, bestDy
}

// inBounds reports whether (x, y) is within the local map grid.
func inBounds(x, y int) bool {
	return x >= 0 && x < LocalMapW && y >= 0 && y < LocalMapH
}

// isBlocking reports whether the cell at (x, y) has a blocking object.
func isBlocking(m *Model, x, y int) bool {
	obj := m.localMap.Objects[x][y]
	return obj != nil && obj.Blocking
}

// moveAnimals advances every animal in m.localMap by one step.
// Flee animals within Manhattan distance 3 of the player move away greedily;
// all others take a random step. Moves that leave the map or land on a
// blocking object are skipped.
func moveAnimals(m *Model) {
	px, py := m.playerPos.X, m.playerPos.Y
	for _, a := range m.localMap.Animals {
		if a.Flee && abs(a.X-px)+abs(a.Y-py) <= 3 {
			// Try all 8 directions, picking the valid one furthest from player.
			bestDx, bestDy := 0, 0
			bestDist := -1
			for _, d := range directions {
				nx, ny := a.X+d[0], a.Y+d[1]
				if !inBounds(nx, ny) || isBlocking(m, nx, ny) {
					continue
				}
				dist := abs(nx-px) + abs(ny-py)
				if dist > bestDist {
					bestDist = dist
					bestDx, bestDy = d[0], d[1]
				}
			}
			if bestDist >= 0 {
				a.X += bestDx
				a.Y += bestDy
			}
			continue
		}
		// Random walk.
		dx, dy := randomStep()
		nx, ny := a.X+dx, a.Y+dy
		if inBounds(nx, ny) && !isBlocking(m, nx, ny) {
			a.X, a.Y = nx, ny
		}
	}
}
