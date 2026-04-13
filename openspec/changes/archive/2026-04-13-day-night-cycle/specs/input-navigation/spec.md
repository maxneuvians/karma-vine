## REMOVED Requirements

### Requirement: n key toggles night mode
**Reason**: Night mode toggle is replaced by automatic time-of-day dimming. Manual control is now time-speed only.
**Migration**: Remove `case "n"` from `Update()`. The `nightMode` field is removed from `Model`.

## ADDED Requirements

### Requirement: [ and ] keys adjust time speed
The system SHALL handle `[` and `]` key messages at all times (both `ModeWorld` and `ModeLocal`). `]` SHALL increase `timeScale` to the next value in `{1, 2, 5, 10}`, clamped at 10. `[` SHALL decrease to the previous value, clamped at 1.

#### Scenario: ] advances time scale from 1 to 2
- **WHEN** `timeScale == 1` and the player presses `]`
- **THEN** `timeScale == 2`

#### Scenario: [ retreats time scale from 10 to 5
- **WHEN** `timeScale == 10` and the player presses `[`
- **THEN** `timeScale == 5`

#### Scenario: ] is ignored at maximum scale
- **WHEN** `timeScale == 10` and the player presses `]`
- **THEN** `timeScale` remains `10`
