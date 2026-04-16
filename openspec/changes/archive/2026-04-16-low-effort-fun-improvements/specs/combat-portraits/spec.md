## MODIFIED Requirements

### Requirement: Enemy portrait is selected by enemy name
The system SHALL implement `enemyPortraitByName(name string) portrait` that maps enemy names to archetypes:
- Humanoid: "Goblin", "Bandit", "Jungle Troll", "Frost Giant", "Stone Golem" → `humanoidPortrait`
- Beast: "Cave Crustacean", "Cave Rat" → `beastPortrait`
- Undead: "Sand Wraith", "Ice Wraith" → `undeadPortrait`
- Any unrecognised name → `fallbackPortrait`

`renderEnemyPanel` SHALL call `enemyPortraitByName(cs.Enemy.Name)` instead of `enemyPortrait(enemyChar)`.

The existing `enemyPortrait(char rune)` function MAY be retained for backward compatibility with existing tests, or tests SHALL be updated to use `enemyPortraitByName`.

#### Scenario: Goblin name returns humanoid portrait
- **WHEN** `enemyPortraitByName("Goblin")` is called
- **THEN** the returned portrait equals `humanoidPortrait`

#### Scenario: Sand Wraith name returns undead portrait
- **WHEN** `enemyPortraitByName("Sand Wraith")` is called
- **THEN** the returned portrait equals `undeadPortrait`

#### Scenario: Cave Crustacean name returns beast portrait
- **WHEN** `enemyPortraitByName("Cave Crustacean")` is called
- **THEN** the returned portrait equals `beastPortrait`

#### Scenario: Unknown name returns fallback portrait
- **WHEN** `enemyPortraitByName("Unknown Monster")` is called
- **THEN** the returned portrait equals `fallbackPortrait`

#### Scenario: renderEnemyPanel uses name-based portrait selection
- **WHEN** `m.combatDungeonEnemy.Template.Name == "Goblin"` and `renderEnemyPanel` is called
- **THEN** the rendered output uses the humanoid portrait art (non-empty, contains block characters)
