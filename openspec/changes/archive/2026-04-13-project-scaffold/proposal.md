## Why

The terminal world explorer needs a bootstrapped Go project with the correct module layout, dependency graph, and a runnable BubbleTea skeleton before any feature work can begin. Without a working entry point and typed model foundation, all subsequent changes have nowhere to land.

## What Changes

- Initialize `go.mod` with module path `karma_vine` and Go 1.22
- Pull in direct dependencies: `bubbletea`, `lipgloss`, `opensimplex-go`
- Create `main.go` entry point that starts the BubbleTea program
- Define the top-level `Model` struct with all fields from the brief (world tier, local tier, player, UI)
- Define supporting value types: `WorldCoord`, `LocalCoord`, `ChunkCoord`, `Mode` (ModeWorld / ModeLocal)
- Provide stub `Init`, `Update`, and `View` methods that compile and produce a placeholder render

## Capabilities

### New Capabilities
- `project-scaffold`: Go module initialization, entry point, core model struct and type definitions, stub BubbleTea lifecycle methods

### Modified Capabilities
<!-- none — this is greenfield -->

## Impact

- Creates the repository's Go source tree from scratch
- All future changes depend on the types defined here
- No external APIs affected; no breaking changes possible at this stage
