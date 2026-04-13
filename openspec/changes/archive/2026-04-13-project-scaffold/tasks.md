## 1. Go Module Setup

- [x] 1.1 Run `go mod init karma_vine` to create `go.mod` with Go 1.22
- [x] 1.2 Run `go get github.com/charmbracelet/bubbletea` to add BubbleTea dependency
- [x] 1.3 Run `go get github.com/charmbracelet/lipgloss` to add Lipgloss dependency
- [x] 1.4 Run `go get github.com/ojrac/opensimplex-go` to add OpenSimplex noise dependency
- [x] 1.5 Run `go mod tidy` and verify `go.sum` is generated

## 2. Core Type Definitions

- [x] 2.1 Create `types.go` with `WorldCoord`, `LocalCoord`, and `ChunkCoord` structs (each with `X, Y int`)
- [x] 2.2 Add `Mode` type as `int` with `ModeWorld` and `ModeLocal` constants in `types.go`
- [x] 2.3 Create placeholder `Chunk` struct in `world.go` with at least one unexported field
- [x] 2.4 Create placeholder `LocalMap` struct in `local.go` with at least one unexported field

## 3. Model Struct

- [x] 3.1 Create `model.go` and define the `Model` struct with all fields from the brief: `globalSeed`, `worldPos`, `chunks`, `localMap`, `localCache`, `playerPos`, `viewportW`, `viewportH`, `mode`
- [x] 3.2 Initialise maps in a `NewModel()` constructor (`chunks` and `localCache` must not be nil)

## 4. BubbleTea Lifecycle Stubs

- [x] 4.1 Implement `Init() tea.Cmd` on `Model` returning `nil`
- [x] 4.2 Implement stub `Update(msg tea.Msg) (tea.Model, tea.Cmd)` that handles `tea.KeyMsg` for `q`/`ctrl+c` to call `tea.Quit`
- [x] 4.3 Implement stub `View() string` returning a placeholder string (e.g. `"World Explorer — loading..."`)

## 5. Entry Point

- [x] 5.1 Create `main.go` that constructs `NewModel()` and starts `tea.NewProgram(model, tea.WithAltScreen())`
- [x] 5.2 Handle the error return from `program.Run()` and exit non-zero if it fails
- [x] 5.3 Run `go build ./...` and confirm it exits 0
- [x] 5.4 Run the binary in a terminal and confirm alt-screen activates, placeholder text appears, and `q` exits cleanly
