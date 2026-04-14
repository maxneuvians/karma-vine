## MODIFIED Requirements

### Requirement: Animals move one step per tick with random direction
On each `TickMsg`, **if `m.paused == false`**, the system SHALL move each animal in `localMap.Animals` by one step in a randomly chosen direction from the 8 cardinal and diagonal directions `{(-1,-1), (-1,0), (-1,1), (0,-1), (0,1), (1,-1), (1,0), (1,1)}`. The new position SHALL be clamped to `[0, 41] × [0, 17]`. If the candidate cell contains a blocking `Object`, the move SHALL be skipped and the animal stays put. **When `m.paused == true`**, `moveAnimals` SHALL NOT be called and all animal positions SHALL remain unchanged.

#### Scenario: Animal position changes after a tick
- **WHEN** a `TickMsg` is dispatched, `m.paused == false`, and the animal's candidate cell is unblocked and within bounds
- **THEN** the animal's `X` or `Y` (or both) has changed from its previous value

#### Scenario: Animal does not move outside map bounds
- **WHEN** an animal is at position `{0, 0}` and the random direction is `(-1, -1)` and `m.paused == false`
- **THEN** the animal remains at `{0, 0}` after the tick

#### Scenario: Animals do not move while paused
- **WHEN** `m.paused == true`, `m.mode == ModeLocal`, and a `TickMsg` is dispatched
- **THEN** all animal positions in `m.localMap.Animals` are unchanged after the update
