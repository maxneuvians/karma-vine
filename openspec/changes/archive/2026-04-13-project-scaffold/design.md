## Context

This is a greenfield Go project. There is no existing codebase. The goal is to produce a compiling, runnable skeleton that every subsequent change can build on top of. The module layout, type definitions, and BubbleTea lifecycle stubs must be stable so that parallel feature work (world generation, rendering, input) can proceed without merge conflicts on core types.

## Goals / Non-Goals

**Goals:**
- Establish the `go.mod` / `go.sum` with pinned versions of all three direct dependencies
- Define all shared value types (`WorldCoord`, `LocalCoord`, `ChunkCoord`, `Mode`) in a single `types.go` file
- Define the top-level `Model` struct matching the brief exactly
- Provide minimal `Init`, `Update`, and `View` stubs that compile and show a placeholder screen
- Wire up `tea.Program` in `main.go` with `WithAltScreen()`

**Non-Goals:**
- Any actual world generation, rendering logic, or input handling (those are separate changes)
- Multiplayer / Wish SSH server (out of scope for v1 per the brief)
- CI configuration or Dockerfile

## Decisions

**Single package (`main`) for now** — The project is small enough that a single package avoids import cycles during early development. If the codebase grows, splitting into sub-packages (`world`, `render`, `input`) is straightforward. Alternative considered: `internal/` packages from the start — rejected because it adds friction before the shape of the code is known.

**Types in a dedicated `types.go`** — All value types that cross feature boundaries live in one file. This prevents the "who owns this type?" problem when multiple changes land simultaneously. Alternative: define types inline in each feature file — rejected because it forces import cycles or duplication.

**`opensimplex-go` pinned version** — `github.com/ojrac/opensimplex-go` v1.0.1 is the last stable pure-Go release. The noise API is called in both world and local generation so the version must be locked upfront.

## Risks / Trade-offs

- **Stub View returns empty string** → The terminal will show a blank alt-screen until the rendering change lands. Acceptable for development; not user-visible in production.
- **All fields in one Model struct** → Slightly large struct, but keeps the BubbleTea contract simple and avoids pointer indirection on hot paths. Trade-off: harder to unit-test sub-components in isolation.
