## MODIFIED Requirements

### Requirement: Time-of-day advances continuously each tick
The system SHALL maintain a `timeOfDay float64` field in `Model` in the range `[0, 1)`, where `0` is midnight, `0.25` is 6 AM, `0.5` is noon, and `0.75` is 6 PM. On each `TickMsg`, `timeOfDay` SHALL advance by `timeScale / (ticksPerSecond * secondsPerDay)` **unless `m.paused == true`**, in which case `timeOfDay` SHALL NOT change. One full cycle SHALL take 30 real seconds at `timeScale == 10` (of unpaused time). `timeOfDay` SHALL wrap around at 1.0.

#### Scenario: Time advances on each tick at 10× speed
- **WHEN** sixty `TickMsg` events are dispatched at `timeScale == 10` and `m.paused == false`
- **THEN** `timeOfDay` has advanced by approximately `1.0` (one full day)

#### Scenario: Time wraps at midnight
- **WHEN** `timeOfDay` is `0.99` and a tick advances it past `1.0` and `m.paused == false`
- **THEN** `timeOfDay` wraps to a value less than `1.0` (not ≥ 1.0)

#### Scenario: Time does not advance while paused
- **WHEN** `m.paused == true` and a `TickMsg` is received
- **THEN** `timeOfDay` is unchanged after the update
