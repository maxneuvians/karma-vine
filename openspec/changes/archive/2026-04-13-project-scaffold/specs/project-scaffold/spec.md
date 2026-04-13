## ADDED Requirements

### Requirement: Go module is initialized with correct dependencies
The project SHALL have a valid `go.mod` declaring module `karma_vine`, Go 1.22 minimum, and direct dependencies on `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, and `github.com/ojrac/opensimplex-go`.

#### Scenario: Module compiles with all dependencies present
- **WHEN** `go build ./...` is run from the project root
- **THEN** the command exits 0 with no errors and produces a runnable binary

#### Scenario: Dependencies are reproducible
- **WHEN** `go mod verify` is run
- **THEN** all module hashes match `go.sum` and the command exits 0

### Requirement: Shared value types are defined
The project SHALL define the following types in `types.go`: `WorldCoord{X, Y int}`, `LocalCoord{X, Y int}`, `ChunkCoord{X, Y int}`, `Mode` (an `int` with constants `ModeWorld` and `ModeLocal`).

#### Scenario: Types are accessible from any file in the package
- **WHEN** another `.go` file in the `main` package references `WorldCoord`, `LocalCoord`, `ChunkCoord`, or `ModeWorld`/`ModeLocal`
- **THEN** the file compiles without an "undefined" error

### Requirement: Model struct matches the brief specification
The project SHALL define a `Model` struct containing: `globalSeed int`, `worldPos WorldCoord`, `chunks map[ChunkCoord]*Chunk`, `localMap *LocalMap`, `localCache map[WorldCoord]*LocalMap`, `playerPos LocalCoord`, `viewportW int`, `viewportH int`, `mode Mode`. Placeholder structs `Chunk` and `LocalMap` SHALL be defined with at least one unexported field so they compile.

#### Scenario: Model zero-value is safe to construct
- **WHEN** `Model{}` is instantiated
- **THEN** no nil-pointer panic occurs before any map operations are performed

### Requirement: BubbleTea program starts and exits cleanly
The `main` package SHALL wire `tea.NewProgram(model, tea.WithAltScreen())` and call `.Run()`. The program SHALL exit when `q` or `ctrl+c` is pressed.

#### Scenario: Program starts without panic
- **WHEN** the compiled binary is executed in a terminal
- **THEN** the alt-screen activates and the stub view renders without a panic

#### Scenario: Program exits on quit key
- **WHEN** the user presses `q` or `ctrl+c`
- **THEN** the program exits with code 0 and restores the normal terminal screen
