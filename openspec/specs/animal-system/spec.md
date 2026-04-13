## Requirements

### Requirement: A TickMsg is dispatched every 500 ms while in local mode
The system SHALL define a `TickMsg` type. `Init()` SHALL return `tea.Every(500*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg{} })`. `Update()` SHALL re-schedule the next tick by returning the same command alongside the updated model whenever a `TickMsg` is handled.

#### Scenario: Tick fires after 500 ms
- **WHEN** the model is initialised and 500 ms elapses
- **THEN** `Update` receives a `TickMsg`

#### Scenario: Tick re-schedules itself
- **WHEN** `Update` handles a `TickMsg`
- **THEN** the returned `tea.Cmd` is non-nil (contains the next tick schedule)

### Requirement: Animals move one step per tick with random direction
On each `TickMsg`, the system SHALL move each animal in `localMap.Animals` by one step in a randomly chosen direction from the 8 cardinal and diagonal directions `{(-1,-1), (-1,0), (-1,1), (0,-1), (0,1), (1,-1), (1,0), (1,1)}`. The new position SHALL be clamped to `[0, 41] × [0, 17]`. If the candidate cell contains a blocking `Object`, the move SHALL be skipped and the animal stays put.

#### Scenario: Animal position changes after a tick
- **WHEN** a `TickMsg` is dispatched and the animal's candidate cell is unblocked and within bounds
- **THEN** the animal's `X` or `Y` (or both) has changed from its previous value

#### Scenario: Animal does not move outside map bounds
- **WHEN** an animal is at position `{0, 0}` and the random direction is `(-1, -1)`
- **THEN** the animal remains at `{0, 0}` after the tick

### Requirement: Animals with Flee:true move away from the player when within 3 tiles
For each animal where `Flee == true`, if the Manhattan distance to `playerPos` is ≤ 3, the system SHALL choose the direction (from the 8 candidates) that maximises the resulting Manhattan distance from `playerPos`. If the best direction is blocked or out-of-bounds, the next-best direction SHALL be tried; if none are valid, the animal stays put.

#### Scenario: Flee animal moves away when player is adjacent
- **WHEN** a flee animal is at `{5, 5}` and `playerPos` is `{5, 6}`
- **THEN** after the tick the animal is at a position further from `{5, 6}` than `{5, 5}` was

#### Scenario: Non-flee animal ignores player proximity
- **WHEN** an animal with `Flee == false` is at `{5, 5}` and `playerPos` is `{5, 6}`
- **THEN** the animal's movement is random (not consistently away from player)

### Requirement: Animal state persists in localCache on revisit
The system SHALL NOT reset animal positions when the player returns to a previously-visited world tile. Since animals are stored as pointers in `LocalMap` and `LocalMap` is cached in `Model.localCache`, their updated positions SHALL be retained automatically.

#### Scenario: Animals resume positions after leaving and returning
- **WHEN** the player descends to a tile, waits for animals to move, ascends, then descends again to the same tile
- **THEN** the animals are at the positions they occupied when the player ascended, not their initial generated positions
