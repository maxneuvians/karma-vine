## ADDED Requirements

### Requirement: Mouse support is declared in View()
The system SHALL set `v.MouseMode = tea.MouseModeCellMotion` on the `tea.View` struct returned by `View()` in `model.go`. This causes BubbleTea v2 to request mouse cell-motion tracking from the terminal, enabling `tea.MouseClickMsg` and `tea.MouseWheelMsg` events to be delivered to `Update()`. `tea.WithMouseCellMotion()` does not exist in v2 and SHALL NOT be used.

#### Scenario: Left-click delivers a MouseClickMsg
- **WHEN** mouse support is enabled and the user left-clicks on the terminal
- **THEN** `Update()` receives a `tea.MouseClickMsg` with `Button == tea.MouseLeft` and `Mouse().X`, `Mouse().Y` populated

#### Scenario: Scroll-wheel delivers a MouseWheelMsg
- **WHEN** the user scrolls the mouse wheel up
- **THEN** `Update()` receives a `tea.MouseWheelMsg` with `Button == tea.MouseWheelUp`

### Requirement: handleMouseClick and handleMouseWheel dispatch mouse events
The system SHALL implement `handleMouseClick(msg tea.MouseClickMsg, m Model) (Model, tea.Cmd)` and `handleMouseWheel(msg tea.MouseWheelMsg, m Model) (Model, tea.Cmd)` in `input.go`. `Update()` SHALL route `tea.MouseClickMsg` events to `handleMouseClick` and `tea.MouseWheelMsg` events to `handleMouseWheel` via separate `case` branches. All other message types continue to route as before.

#### Scenario: Update routes MouseClickMsg to handleMouseClick
- **WHEN** `Update()` receives a `tea.MouseClickMsg`
- **THEN** `handleMouseClick` is called and its returned model is used

#### Scenario: Update routes MouseWheelMsg to handleMouseWheel
- **WHEN** `Update()` receives a `tea.MouseWheelMsg`
- **THEN** `handleMouseWheel` is called and its returned model is used

### Requirement: Left-click in ScreenNormal moves player one step toward click
When `m.screenMode == ScreenNormal` and `m.mode == ModeLocal` or `m.mode == ModeDungeon`, a left-click (`tea.MouseClickMsg` with `Button == tea.MouseLeft`) at terminal position `(msg.Mouse().X, msg.Mouse().Y)` SHALL cause the player to take one cardinal step toward the clicked map cell, using the same collision logic as `applyDelta`. The click SHALL be ignored if the sidebar, map picker, or inventory is open (as coordinate mapping would be incorrect).

The step direction SHALL be determined by: compute `dx = clickMapX - playerPos.X`, `dy = clickMapY - playerPos.Y`; if `|dx| >= |dy|` move horizontally (`dx > 0` → +1, `dx < 0` → -1`); else move vertically. If the click is on the player's own cell, it is a no-op.

#### Scenario: Click to the right moves player right one step
- **WHEN** `screenMode == ScreenNormal`, `mode == ModeLocal`, no panels open, and the player left-clicks a cell to the right of the player
- **THEN** `playerPos.X` increases by 1 (subject to collision)

#### Scenario: Click on own cell is a no-op
- **WHEN** the player left-clicks the cell they are already on
- **THEN** `playerPos` is unchanged

#### Scenario: Click ignored when sidebar is open
- **WHEN** `showSidebar == true` and the player left-clicks
- **THEN** `playerPos` is unchanged

#### Scenario: Click ignored on world map
- **WHEN** `m.mode == ModeWorld` and the player left-clicks
- **THEN** `worldPos` is unchanged

### Requirement: Scroll-wheel in ScreenInventory moves inventory cursor
When `m.screenMode == ScreenInventory`, `tea.MouseWheelMsg` with `Button == tea.MouseWheelUp` SHALL decrement `inventoryCursor` (clamped at 0) and `Button == tea.MouseWheelDown` SHALL increment it (clamped at `len(inventory.Items)-1`).

#### Scenario: Scroll up moves cursor up
- **WHEN** `screenMode == ScreenInventory` and the user scrolls up
- **THEN** `inventoryCursor` decrements (clamped at 0)

#### Scenario: Scroll down moves cursor down
- **WHEN** `screenMode == ScreenInventory` and the user scrolls down
- **THEN** `inventoryCursor` increments (clamped at `len(inventory.Items)-1`)

### Requirement: Left-click on inventory row selects that row
When `m.screenMode == ScreenInventory`, a left-click at terminal row `my` SHALL map the click to an inventory row index and set `inventoryCursor` to that index, clamped to `[0, len(inventory.Items)-1]`. Clicks outside the item list area (header rows, ragdoll column) are no-ops.

#### Scenario: Click on item row sets cursor to that row
- **WHEN** `screenMode == ScreenInventory` and the player clicks on the second item row
- **THEN** `inventoryCursor == 1`

#### Scenario: Click outside item list is a no-op
- **WHEN** `screenMode == ScreenInventory` and the player clicks on the header
- **THEN** `inventoryCursor` is unchanged
