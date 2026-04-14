## 0. BubbleTea v2 Migration

- [x] 0.1 Update `go.mod`: run `go get charm.land/bubbletea/v2@latest charm.land/lipgloss/v2@latest` and `go mod tidy` to switch from the GitHub paths to the vanity domain
- [x] 0.2 Update all import paths in `*.go` source files: `github.com/charmbracelet/bubbletea` → `charm.land/bubbletea/v2`; `github.com/charmbracelet/lipgloss` → `charm.land/lipgloss/v2`
- [x] 0.3 Change `View() string` to `View() tea.View` in `model.go`; return `tea.NewView(buildView(m))` with `AltScreen = true` (leaves `buildView` signature unchanged)
- [x] 0.4 Remove `tea.WithAltScreen()` from `tea.NewProgram(...)` in `main.go` (option no longer exists in v2)
- [x] 0.5 Change `case tea.KeyMsg:` to `case tea.KeyPressMsg:` in `Update()` in `model.go`; update `handleKey` signature in `input.go` from `tea.KeyMsg` to `tea.KeyPressMsg`
- [x] 0.6 Update all test key constructions: `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}` → `tea.KeyPressMsg{Code: 'x', Text: "x"}`; `tea.KeyMsg{Type: tea.KeyUp}` → `tea.KeyPressMsg{Code: tea.KeyUp}` (same for Down, Left, Right, Enter)
- [x] 0.7 Run `go build ./...` and `go test ./...` and confirm they pass before proceeding to later groups

## 1. Types

- [x] 1.1 Add `ScreenMode int` type with `ScreenNormal` and `ScreenInventory` constants to `types.go`

## 2. Model

- [x] 2.1 Replace `showInventory bool` with `screenMode ScreenMode` on `Model` in `model.go`
- [x] 2.2 Remove `showInventory` initialisation from `NewModel()` (zero value of `screenMode` is already `ScreenNormal`)

## 3. Mouse Support — Program Init

- [x] 3.1 Set `v.MouseMode = tea.MouseModeCellMotion` on the `tea.View` returned by `View()` in `model.go` (extend the `tea.NewView` call added in task 0.3)

## 4. Mouse Support — Input Handling

- [x] 4.1 Add `case tea.MouseClickMsg:` and `case tea.MouseWheelMsg:` branches to `Update()` in `model.go`, routing to `handleMouseClick` and `handleMouseWheel` respectively
- [x] 4.2 Implement `handleMouseClick(msg tea.MouseClickMsg, m Model) (Model, tea.Cmd)` and `handleMouseWheel(msg tea.MouseWheelMsg, m Model) (Model, tea.Cmd)` in `input.go`
- [x] 4.3 In `handleMouseWheel`: when `screenMode == ScreenInventory`, handle `msg.Button == tea.MouseWheelUp` (decrement `inventoryCursor`, clamp at 0) and `tea.MouseWheelDown` (increment, clamp at `len(inventory.Items)-1`)
- [x] 4.4 In `handleMouseClick`: when `screenMode == ScreenInventory` and `msg.Button == tea.MouseLeft`, map `msg.Mouse().Y` to an inventory row index and set `inventoryCursor`; clicks outside the item list are no-ops
- [x] 4.5 In `handleMouseClick`: when `screenMode == ScreenNormal` and `mode == ModeLocal` or `ModeDungeon` and no panels open, handle `msg.Button == tea.MouseLeft` using `msg.Mouse().X`, `msg.Mouse().Y` to take one cardinal step toward clicked cell via `applyDelta`

## 5. Input — Keyboard Migration

- [x] 5.1 Update `i` key handler: set `m.screenMode = ScreenInventory` (was `m.showInventory = true`); toggle to `ScreenNormal` when already `ScreenInventory`
- [x] 5.2 Update `esc` key handler: when `screenMode == ScreenInventory`, set `screenMode = ScreenNormal` before existing dungeon/local ascent logic
- [x] 5.3 Update inventory cursor keys (`up`/`w`, `down`/`s`): check `m.screenMode == ScreenInventory` instead of `m.showInventory`
- [x] 5.4 Update `d` (drop) and `u` (use) handlers: check `m.screenMode == ScreenInventory` instead of `m.showInventory`

## 6. Rendering — Fullscreen Inventory

- [x] 6.1 Implement `renderFullscreenInventory(m Model) string` in `render.go` with two-column layout
- [x] 6.2 Left column: title "Inventory", separator, item rows with glyph+name+count, cursor highlight, "Empty" placeholder, hint row at bottom
- [x] 6.3 Right column: ASCII ragdoll body outline (`~O~`, `|H|`, etc.) centred vertically with slot labels: Head, Chest, Left Hand, Right Hand, Legs, Feet — each showing `[ Empty ]`
- [x] 6.4 Use `lipgloss.JoinHorizontal` to combine columns to `m.viewportW` width and `m.viewportH` height

## 7. Rendering — buildView Migration

- [x] 7.1 Add `if m.screenMode == ScreenInventory { return renderFullscreenInventory(m) }` at the top of `buildView`
- [x] 7.2 Remove the old `if m.showInventory` side-panel block from `buildView`
- [x] 7.3 Remove the `inventoryPanelW` constant and associated panel style variables that are no longer used

## 8. Tests

- [x] 8.1 Update all existing tests that set/check `m.showInventory` to use `m.screenMode`; update all `tea.KeyMsg{...}` constructions to `tea.KeyPressMsg{...}` (covered by group 0 task 0.6 but verify no stragglers)
- [x] 8.2 Test: `NewModel()` sets `screenMode == ScreenNormal`
- [x] 8.3 Test: `i` key sets `screenMode` to `ScreenInventory`; second `i` returns to `ScreenNormal`
- [x] 8.4 Test: `esc` closes inventory (sets `screenMode = ScreenNormal`) when `screenMode == ScreenInventory`
- [x] 8.5 Test: `tea.MouseWheelMsg{Button: tea.MouseWheelUp}` decrements `inventoryCursor` when `screenMode == ScreenInventory`; `MouseWheelDown` increments it
- [x] 8.6 Test: `tea.MouseClickMsg{Button: tea.MouseLeft, ...}` on item row sets `inventoryCursor` when `screenMode == ScreenInventory`
- [x] 8.7 Test: `tea.MouseClickMsg{Button: tea.MouseLeft, ...}` in `ScreenNormal`/`ModeLocal` moves player one step toward click
- [x] 8.8 Test: `tea.MouseClickMsg` ignored (player does not move) when sidebar open
- [x] 8.9 Test: `renderFullscreenInventory` output contains "Inventory", "Head", "Chest", "Left Hand"
- [x] 8.10 Test: `renderFullscreenInventory` shows "Empty" when no items
- [x] 8.11 Test: `buildView` returns fullscreen inventory content when `screenMode == ScreenInventory`
- [x] 8.12 Test: `buildView` returns normal map when `screenMode == ScreenNormal`
