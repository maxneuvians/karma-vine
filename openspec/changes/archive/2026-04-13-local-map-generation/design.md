## Context

Local maps must be 42×18 (matching the brief) so they fill the terminal viewport. Each cell is a stack: ground first, then optional object, then optional animal. Two noise passes give independent control over terrain shape and object density. The `hash` function from the brief converts integer coordinates and a seed into a stable float64 for local-seed derivation — this is the glue that ties deterministic local generation to world position.

## Goals / Non-Goals

**Goals:**
- `hash(x, y, seed int) float64` matches the exact bit-manipulation formula in the brief
- `Ground` has `Char rune`, `Color string`, `Passable bool`
- `Object` has `Char rune`, `Color string`, `Blocking bool`
- `Animal` has `X, Y int`, `Char rune`, `Color string`, `Flee bool` (movement logic deferred to `animal-system`)
- `LocalMap` has `[42][18]Ground`, `[42][18]*Object`, `[]*Animal`
- Object placement controlled by a per-biome `objectThreshold float64` constant
- Animals placed at random passable positions up to a per-biome maximum count

**Non-Goals:**
- Animal movement / tick (that is `animal-system`)
- Rendering local maps (that is `rendering-system`)
- Interaction / inventory (v1 out of scope)

## Decisions

**Separate noise objects for terrain and objects** — Same rationale as world generation: two independent seeds produce uncorrelated terrain-shape and object-density fields. Seeds are `localSeed` and `localSeed + 0.5` (float offset on the same noise object is also valid, but separate objects are cleaner).

**Object threshold as a per-biome constant** — Different biomes have different object densities (forest is cluttered; desert is sparse). A single global threshold would produce uniform density. Alternative: per-tile random roll — rejected because it bypasses the noise field and loses spatial coherence.

**Animals placed once at generation; movement handled elsewhere** — Keeping generation pure (no time dependency) makes the function testable and ensures the cached map has stable initial state. The `animal-system` change will add tick-based movement on top.

## Risks / Trade-offs

- **Fixed map size 42×18** → If the viewport is smaller the map will be clipped. Rendering change must handle this. Mitigation: the rendering change reads `viewportW`/`viewportH` and crops accordingly.
- **No object overlap check** → Two objects could be placed on adjacent cells in the same chunk — fine for v1.
- **`hash` uses signed integer arithmetic** → Go's integer overflow behaviour is defined (wraps), so the formula is safe and reproducible across platforms.
