## Requirements

### Requirement: Import paths use the charm.land vanity domain
The system SHALL update `go.mod` and all Go source files to use the v2 module paths:
- `github.com/charmbracelet/bubbletea` → `charm.land/bubbletea/v2`
- `github.com/charmbracelet/lipgloss` → `charm.land/lipgloss/v2`

`go mod tidy` SHALL be run after updating the import paths to remove the v1 modules and resolve indirect dependencies.

#### Scenario: Build succeeds after path migration
- **WHEN** all import paths are updated and `go mod tidy` is run
- **THEN** `go build ./...` completes with no errors

### Requirement: View() returns tea.View
The `View()` method on `Model` SHALL return `tea.View` instead of `string`. The returned `tea.View` SHALL set:
- Content via `tea.NewView(buildView(m))`
- `AltScreen = true` (replaces `tea.WithAltScreen()` program option)
- `MouseMode = tea.MouseModeCellMotion` (replaces `tea.WithMouseCellMotion()` program option)

`tea.WithAltScreen()` and `tea.WithMouseCellMotion()` SHALL be removed from `tea.NewProgram(...)` in `main.go` as these options no longer exist in v2; terminal feature flags live in the `tea.View` struct instead.

#### Scenario: View returns AltScreen and mouse mode via View struct
- **WHEN** `View()` is called
- **THEN** the returned `tea.View` has `AltScreen == true` and `MouseMode == tea.MouseModeCellMotion`

#### Scenario: main.go NewProgram call has no options
- **WHEN** the program is built
- **THEN** `tea.NewProgram(game.NewModel())` takes no option arguments

### Requirement: Key press events use tea.KeyPressMsg
`Update()` SHALL match on `case tea.KeyPressMsg:` instead of `case tea.KeyMsg:`. The `handleKey` function in `input.go` SHALL accept `tea.KeyPressMsg`. All test code constructing key messages SHALL use `tea.KeyPressMsg`:

| v1 | v2 |
|---|---|
| `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")}` | `tea.KeyPressMsg{Code: 'x', Text: "x"}` |
| `tea.KeyMsg{Type: tea.KeyUp}` | `tea.KeyPressMsg{Code: tea.KeyUp}` |
| `tea.KeyMsg{Type: tea.KeyDown}` | `tea.KeyPressMsg{Code: tea.KeyDown}` |
| `tea.KeyMsg{Type: tea.KeyLeft}` | `tea.KeyPressMsg{Code: tea.KeyLeft}` |
| `tea.KeyMsg{Type: tea.KeyRight}` | `tea.KeyPressMsg{Code: tea.KeyRight}` |
| `tea.KeyMsg{Type: tea.KeyEnter}` | `tea.KeyPressMsg{Code: tea.KeyEnter}` |

The `msg.String()` method continues to work as before — all existing key string matches (`"q"`, `"ctrl+c"`, `"up"`, `"esc"`, `"i"`, etc.) remain valid in v2.

#### Scenario: Key press is delivered via KeyPressMsg
- **WHEN** the user presses a key
- **THEN** `Update()` receives a `tea.KeyPressMsg` and routes to `handleKey`

#### Scenario: Test key construction uses tea.KeyPressMsg
- **WHEN** test files construct `tea.KeyPressMsg{Code: tea.KeyUp}` to simulate arrow key presses
- **THEN** the test compiles and `handleKey` behaves identically to before the migration

### Requirement: Mouse events use split message types
In v2, mouse events are split into distinct message types. `Update()` SHALL match on the specific types:
- `tea.MouseClickMsg` — left/right/middle button presses; coordinates via `msg.Mouse().X`, `msg.Mouse().Y`; button via `msg.Button` (`tea.MouseLeft`, `tea.MouseRight`, `tea.MouseMiddle`)
- `tea.MouseWheelMsg` — scroll events; direction via `msg.Button` (`tea.MouseWheelUp`, `tea.MouseWheelDown`)
- `tea.MouseReleaseMsg` and `tea.MouseMotionMsg` — no `case` needed; not handled by this change

The monolithic `tea.MouseMsg` struct from v1 is replaced by these distinct types. Coordinates in all v2 mouse events are accessed via the `msg.Mouse()` method returning a `tea.Mouse` struct with `X` and `Y` fields.

#### Scenario: Left click delivers MouseClickMsg
- **WHEN** the user left-clicks on the terminal
- **THEN** `Update()` receives `tea.MouseClickMsg` with `Button == tea.MouseLeft` and `Mouse().X`, `Mouse().Y` populated

#### Scenario: Scroll delivers MouseWheelMsg
- **WHEN** the user scrolls the mouse wheel up
- **THEN** `Update()` receives `tea.MouseWheelMsg` with `Button == tea.MouseWheelUp`
