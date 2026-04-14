package game

import (
	"math/rand"
)

// ── BSP helpers ─────────────────────────────────────────────────────────────

type bspRect struct {
	x, y, w, h int
}

type bspNode struct {
	rect        bspRect
	left, right *bspNode
	room        *bspRect // non-nil only for leaf nodes
}

const bspMinLeaf = 6

// bspSplit recursively splits a rectangle into leaf nodes.
func bspSplit(rect bspRect, rng *rand.Rand) *bspNode {
	node := &bspNode{rect: rect}

	// Stop splitting if below minimum size for both axes.
	if rect.w < bspMinLeaf*2 && rect.h < bspMinLeaf*2 {
		return node
	}

	// Decide split direction. Prefer splitting the longer axis.
	horizontal := rng.Intn(2) == 0
	if rect.w > rect.h*2 {
		horizontal = false
	} else if rect.h > rect.w*2 {
		horizontal = true
	}

	if horizontal {
		if rect.h < bspMinLeaf*2 {
			return node // can't split horizontally
		}
		split := bspMinLeaf + rng.Intn(rect.h-bspMinLeaf*2+1)
		node.left = bspSplit(bspRect{rect.x, rect.y, rect.w, split}, rng)
		node.right = bspSplit(bspRect{rect.x, rect.y + split, rect.w, rect.h - split}, rng)
	} else {
		if rect.w < bspMinLeaf*2 {
			return node // can't split vertically
		}
		split := bspMinLeaf + rng.Intn(rect.w-bspMinLeaf*2+1)
		node.left = bspSplit(bspRect{rect.x, rect.y, split, rect.h}, rng)
		node.right = bspSplit(bspRect{rect.x + split, rect.y, rect.w - split, rect.h}, rng)
	}

	return node
}

// bspLeaves collects all leaf nodes from the BSP tree.
func bspLeaves(node *bspNode) []*bspNode {
	if node.left == nil && node.right == nil {
		return []*bspNode{node}
	}
	var leaves []*bspNode
	if node.left != nil {
		leaves = append(leaves, bspLeaves(node.left)...)
	}
	if node.right != nil {
		leaves = append(leaves, bspLeaves(node.right)...)
	}
	return leaves
}

// ── Dungeon level generation ────────────────────────────────────────────────

// GenerateDungeonLevel produces a deterministic dungeon level using BSP.
func GenerateDungeonLevel(globalSeed, wx, wy, depth, maxDepth int, biome Biome) *DungeonLevel {
	seed := int64(globalSeed ^ wx*31337 ^ wy*1619 ^ depth*7919)
	rng := rand.New(rand.NewSource(seed))

	level := &DungeonLevel{}

	// All cells start as CellWall (zero value).

	// BSP split the dungeon area (leaving a 1-cell border).
	root := bspSplit(bspRect{1, 1, DungeonW - 2, DungeonH - 2}, rng)
	leaves := bspLeaves(root)

	// Carve rooms inside each leaf.
	rooms := make([]roomInfo, 0, len(leaves))

	for _, leaf := range leaves {
		r := leaf.rect
		// Room must be at least 3×3 interior, leave 1-cell wall padding inside the leaf.
		maxRW := r.w - 2
		maxRH := r.h - 2
		if maxRW < 3 {
			maxRW = 3
		}
		if maxRH < 3 {
			maxRH = 3
		}
		rw := 3 + rng.Intn(maxRW-2)
		if rw > r.w-2 {
			rw = r.w - 2
		}
		rh := 3 + rng.Intn(maxRH-2)
		if rh > r.h-2 {
			rh = r.h - 2
		}
		rx := r.x + 1 + rng.Intn(r.w-rw-1)
		if rx+rw > r.x+r.w-1 {
			rx = r.x + r.w - 1 - rw
		}
		ry := r.y + 1 + rng.Intn(r.h-rh-1)
		if ry+rh > r.y+r.h-1 {
			ry = r.y + r.h - 1 - rh
		}

		// Bounds safety.
		if rx < 1 {
			rx = 1
		}
		if ry < 1 {
			ry = 1
		}
		if rx+rw >= DungeonW {
			rw = DungeonW - rx - 1
		}
		if ry+rh >= DungeonH {
			rh = DungeonH - ry - 1
		}

		for x := rx; x < rx+rw; x++ {
			for y := ry; y < ry+rh; y++ {
				level.Cells[x][y].Kind = CellFloor
			}
		}

		rm := roomInfo{x: rx, y: ry, w: rw, h: rh, cx: rx + rw/2, cy: ry + rh/2}
		leaf.room = &bspRect{rx, ry, rw, rh}
		rooms = append(rooms, rm)
	}

	// Connect adjacent rooms with L-shaped corridors.
	for i := 1; i < len(rooms); i++ {
		carveCorridor(level, rooms[i-1].cx, rooms[i-1].cy, rooms[i].cx, rooms[i].cy)
	}

	// Place up-staircase in the first room.
	upX := rooms[0].cx
	upY := rooms[0].cy
	level.UpStair = LocalCoord{X: upX, Y: upY}
	level.Cells[upX][upY].Object = &Object{Char: '<', Color: "#e8c96a", Blocking: false, Name: "Staircase Up"}

	// Place down-staircase in a different room if not final level.
	if depth < maxDepth && len(rooms) > 1 {
		downRoom := rooms[len(rooms)-1]
		downX := downRoom.cx
		downY := downRoom.cy
		level.DownStair = LocalCoord{X: downX, Y: downY}
		level.HasDownStair = true
		level.Cells[downX][downY].Object = &Object{Char: '>', Color: "#e8c96a", Blocking: false, Name: "Staircase Down"}
	}

	// Place torches on wall cells adjacent to floor cells (~1 per 5 rooms).
	torchCount := len(rooms) / 5
	if torchCount < 1 && len(rooms) > 0 {
		torchCount = 1
	}
	wallAdjacentFloor := collectWallAdjacentFloor(level)
	rng.Shuffle(len(wallAdjacentFloor), func(i, j int) {
		wallAdjacentFloor[i], wallAdjacentFloor[j] = wallAdjacentFloor[j], wallAdjacentFloor[i]
	})
	for i := 0; i < torchCount && i < len(wallAdjacentFloor); i++ {
		p := wallAdjacentFloor[i]
		level.Cells[p.X][p.Y].Object = &Object{Char: '†', Color: "#e8c96a", Blocking: true, Name: "Torch", Pickupable: true}
	}

	// Place braziers on floor cells inside rooms (~1 per 6 rooms).
	brazierCount := len(rooms) / 6
	if brazierCount < 1 && len(rooms) > 0 {
		brazierCount = 1
	}
	floorCells := collectFloorCellsInRooms(level, rooms)
	rng.Shuffle(len(floorCells), func(i, j int) {
		floorCells[i], floorCells[j] = floorCells[j], floorCells[i]
	})
	placed := 0
	for _, p := range floorCells {
		if placed >= brazierCount {
			break
		}
		// Don't place on staircases or existing objects.
		if level.Cells[p.X][p.Y].Object != nil {
			continue
		}
		level.Cells[p.X][p.Y].Object = &Object{Char: 'Ω', Color: "#e07030", Blocking: false, Name: "Brazier", Pickupable: true}
		placed++
	}

	// Spawn enemies based on biome and depth.
	tmpl, ok := dungeonEnemyRoster[biome]
	if !ok {
		tmpl = dungeonEnemyFallback
	}
	// Collect all floor positions not occupied by objects or staircases.
	var enemyPositions []LocalCoord
	for x := 0; x < DungeonW; x++ {
		for y := 0; y < DungeonH; y++ {
			if level.Cells[x][y].Kind == CellFloor && level.Cells[x][y].Object == nil {
				enemyPositions = append(enemyPositions, LocalCoord{X: x, Y: y})
			}
		}
	}
	rng.Shuffle(len(enemyPositions), func(i, j int) {
		enemyPositions[i], enemyPositions[j] = enemyPositions[j], enemyPositions[i]
	})
	enemyCount := depth
	if enemyCount > len(enemyPositions) {
		enemyCount = len(enemyPositions)
	}
	for i := 0; i < enemyCount; i++ {
		p := enemyPositions[i]
		level.Enemies = append(level.Enemies, spawnEnemy(tmpl, p.X, p.Y, depth, maxDepth))
	}

	return level
}

// carveCorridor carves an L-shaped corridor between two points.
func carveCorridor(level *DungeonLevel, x1, y1, x2, y2 int) {
	// Horizontal then vertical.
	x := x1
	for x != x2 {
		if x >= 0 && x < DungeonW && y1 >= 0 && y1 < DungeonH {
			level.Cells[x][y1].Kind = CellFloor
		}
		if x < x2 {
			x++
		} else {
			x--
		}
	}
	y := y1
	for y != y2 {
		if x2 >= 0 && x2 < DungeonW && y >= 0 && y < DungeonH {
			level.Cells[x2][y].Kind = CellFloor
		}
		if y < y2 {
			y++
		} else {
			y--
		}
	}
	// Carve the final cell.
	if x2 >= 0 && x2 < DungeonW && y2 >= 0 && y2 < DungeonH {
		level.Cells[x2][y2].Kind = CellFloor
	}
}

// collectWallAdjacentFloor returns wall cells that are adjacent to at least one floor cell.
func collectWallAdjacentFloor(level *DungeonLevel) []LocalCoord {
	var result []LocalCoord
	for x := 1; x < DungeonW-1; x++ {
		for y := 1; y < DungeonH-1; y++ {
			if level.Cells[x][y].Kind != CellWall || level.Cells[x][y].Object != nil {
				continue
			}
			// Check 4-neighbours for floor.
			if level.Cells[x-1][y].Kind == CellFloor ||
				level.Cells[x+1][y].Kind == CellFloor ||
				level.Cells[x][y-1].Kind == CellFloor ||
				level.Cells[x][y+1].Kind == CellFloor {
				result = append(result, LocalCoord{X: x, Y: y})
			}
		}
	}
	return result
}

type roomInfo struct {
	x, y, w, h int
	cx, cy      int
}

// collectFloorCellsInRooms returns floor cells that are inside rooms and don't have objects.
func collectFloorCellsInRooms(level *DungeonLevel, rooms []roomInfo) []LocalCoord {
	var result []LocalCoord
	for _, rm := range rooms {
		for x := rm.x; x < rm.x+rm.w; x++ {
			for y := rm.y; y < rm.y+rm.h; y++ {
				if x >= 0 && x < DungeonW && y >= 0 && y < DungeonH &&
					level.Cells[x][y].Kind == CellFloor {
					result = append(result, LocalCoord{X: x, Y: y})
				}
			}
		}
	}
	return result
}

// ── Cache accessors ─────────────────────────────────────────────────────────

// DungeonLevelFor returns the DungeonLevel for the given coordinates, generating
// and caching it on first access.
func DungeonLevelFor(wx, wy, depth int, m *Model) *DungeonLevel {
	key := dungeonKey{wx: wx, wy: wy, depth: depth}
	if dl, ok := m.dungeonCache[key]; ok {
		return dl
	}
	meta := DungeonMetaFor(wx, wy, m)
	dl := GenerateDungeonLevel(m.globalSeed, wx, wy, depth, meta.MaxDepth, meta.Biome)
	m.dungeonCache[key] = dl
	return dl
}

// DungeonMetaFor returns the DungeonMeta for the given world coordinate,
// creating one with a randomised MaxDepth on first access.
func DungeonMetaFor(wx, wy int, m *Model) DungeonMeta {
	key := WorldCoord{X: wx, Y: wy}
	if meta, ok := m.dungeonMeta[key]; ok {
		return meta
	}
	seed := int64(m.globalSeed ^ wx*31337 ^ wy*1619 ^ 9973)
	rng := rand.New(rand.NewSource(seed))
	meta := DungeonMeta{MaxDepth: 5 + rng.Intn(6)} // [5, 10]
	m.dungeonMeta[key] = meta
	return meta
}
