# Karma Vine

A terminal-based world exploration game written in Go. Navigate a procedurally generated world, descend into local areas, and explore multi-level dungeons — all rendered in your terminal using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework.

## Features

- **Infinite procedural world** — OpenSimplex noise generates biomes, elevation, temperature zones, and day/night lighting
- **Three-tier exploration** — World map → local area → dungeon levels (up to 5–10 floors deep)
- **Dungeon system** — BSP-generated rooms and corridors, fog-of-war visibility, torches and braziers providing light
- **Map overlay modes** — Toggle between default, temperature, elevation, and political/contour views (`m` key)
- **Day/night cycle** — Time advances in real time, dimming tile colours at night; fires provide local illumination
- **Animals** — Biome-appropriate wildlife roams local maps

## Prerequisites

- Go 1.25+
- A terminal with Unicode and 256-colour support (most modern terminals)
- [OpenSpec CLI](https://openspec.dev) — required only if you want to add features using the spec-driven workflow

## Getting Started

```bash
# Clone the repo
git clone <repo-url> karma_vine
cd karma_vine

# Build and run
make run
```

## Controls

| Key | Action |
|-----|--------|
| `↑ ↓ ← →` or `w a s d` | Move |
| `enter` / `>` | Descend (world → local → dungeon) |
| `esc` / `<` | Ascend (dungeon → local → world) |
| `m` | Toggle map mode picker (world map only) |
| `?` | Toggle info sidebar |
| `]` / `[` | Speed up / slow down time |
| `+` / `-` | Zoom world map in/out |
| `q` or `ctrl+c` | Quit |

## Development

```bash
make build          # Compile to ./build/karma-vine
make test           # Run all tests
make test-coverage  # Run tests with coverage report
make lint           # Vet all packages
make clean          # Remove build artefacts
```

The project targets ≥ 90% test coverage. All new features should include unit tests.

---

## Adding Features with OpenSpec

This project uses [OpenSpec](https://openspec.dev) — a spec-driven development workflow — to plan, track, and archive changes. Every shipped feature has a corresponding set of artefacts (proposal, design, specs, tasks) that live in `openspec/`.

```
openspec/
├── config.yaml          # Project context and schema settings
├── specs/               # Living specifications (one folder per capability)
│   ├── dungeon-generation/spec.md
│   ├── rendering-system/spec.md
│   └── ...
└── changes/
    ├── <active-changes>/  # In-flight work
    └── archive/           # Completed changes with full history
```

### Workflow Overview

```
propose → apply → archive
```

1. **Propose** — Describe what you want to build. OpenSpec generates a proposal, design doc, capability specs, and a task checklist.
2. **Apply** — Implement the tasks. Mark each `- [ ]` as `- [x]` as you go.
3. **Archive** — When all tasks are done, archive the change. OpenSpec syncs the delta specs into the main `openspec/specs/` directory.

### Step 1 — Propose a change

```bash
openspec new change "my-feature-name"
```

Or use the AI-assisted workflow with Claude Code:

```
/opsx:propose I want to add <description of your feature>
```

This creates `openspec/changes/my-feature-name/` and walks you through generating:

| Artefact | Purpose |
|----------|---------|
| `proposal.md` | The *why* — motivation, what changes, which capabilities are affected |
| `design.md` | The *how* — architectural decisions, risks, trade-offs |
| `specs/<capability>/spec.md` | The *what* — testable requirements with WHEN/THEN scenarios |
| `tasks.md` | Implementation checklist (`- [ ] 1.1 ...`) |

Check the current state of all artefacts at any time:

```bash
openspec status --change "my-feature-name"
```

### Step 2 — Implement

Work through the tasks in `tasks.md`. Mark each task complete as you finish it:

```markdown
- [x] 1.1 Add new type to types.go
- [ ] 1.2 Implement generator function   ← currently working on
```

With Claude Code:

```
/opsx:apply my-feature-name
```

The AI reads the context files and implements remaining tasks one by one, marking each checkbox as it goes.

### Step 3 — Archive

Once all tasks are done:

```bash
openspec archive my-feature-name
```

Or with Claude Code:

```
/opsx:archive my-feature-name
```

This will:
1. Check all artefacts and tasks are complete (warns if not)
2. Offer to sync your delta specs into the main `openspec/specs/` directory
3. Move the change folder to `openspec/changes/archive/YYYY-MM-DD-my-feature-name/`

### Writing Good Specs

Specs live in `openspec/specs/<capability>/spec.md`. They define *what the system should do*, not *how*. Each requirement needs at least one WHEN/THEN scenario — these map directly to unit tests.

```markdown
### Requirement: Player cannot move into a wall
The system SHALL block movement into `CellWall` cells in `ModeDungeon`.

#### Scenario: Movement blocked by wall
- **WHEN** the player presses `right` and the target cell is `CellWall`
- **THEN** `playerPos` remains unchanged
```

**Rules:**
- Use `SHALL`/`MUST` for normative requirements
- Scenarios use exactly `####` (4 hashes) — not 3, not bullets
- Every requirement must have at least one scenario
- New capabilities get a new `specs/<name>/spec.md`
- Changes to existing capabilities use ADDED/MODIFIED/REMOVED sections in a delta spec under `changes/<name>/specs/<capability>/spec.md`

### Project Context

The `openspec/config.yaml` provides project context to all AI-generated artefacts:

```yaml
schema: spec-driven
context: |
  Tech stack: Go
  We use conventional commits
  We use test coverage and want to maintain at least 90% coverage
  Modular codebase with clear separation of concerns
  Domain: Game
```

Adjust this file if the project's conventions change.

### Browsing Existing Specs

To see what's already specified:

```bash
openspec list --specs        # list all capability specs
openspec show dungeon-generation  # view a specific spec
openspec view                # interactive dashboard
```

All shipped capabilities are documented under `openspec/specs/`. When in doubt about how something should behave, check there first.
