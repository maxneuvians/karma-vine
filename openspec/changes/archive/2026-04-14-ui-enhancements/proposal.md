## Why

The current UI panels (inventory, sidebar) are narrow fixed-width columns appended to the map view, which limits information density and makes future features like a ragdoll equipment model impractical. A review of BubbleTea v2's capabilities reveals a cleaner mouse API (split `MouseClickMsg`/`MouseWheelMsg` types), a declarative `tea.View` return type that consolidates terminal feature flags, and a stable vanity-domain import path. Upgrading now — before adding mouse support and fullscreen overlays — means implementing those features against the final v2 API rather than v1 and migrating later.

## What Changes

- Add **mouse support** to the program: left-click on the map moves the player; left-click on an inventory row selects that row; scroll-wheel scrolls the inventory cursor.
- Replace the current narrow inventory column with a **fullscreen inventory overlay** that takes over the terminal using alt-screen and renders a centred panel with more room for future ragdoll equipment slots, tooltips, and item details.
- Introduce a **`ScreenMode`** concept (`ScreenNormal` / `ScreenInventory`) so the main render path cleanly separates the map game view from full-screen overlay views.
- Lay the groundwork for a **ragdoll equipment panel**: the fullscreen inventory shall reserve a right-hand column for an ASCII body-outline with named slot labels (Head, Chest, Legs, Feet, Left Hand, Right Hand); the slots will be empty/placeholder in this change.

## Capabilities

### New Capabilities
- `bubbletea-v2`: Upgrade from `github.com/charmbracelet/bubbletea v1.3.10` to `charm.land/bubbletea/v2`; update lipgloss to `charm.land/lipgloss/v2`; migrate `View() string` → `View() tea.View`; migrate `tea.KeyMsg` → `tea.KeyPressMsg`; adopt split mouse message types.
- `mouse-support`: Handling `tea.MouseClickMsg` and `tea.MouseWheelMsg` events (v2 API); click-to-move on local/dungeon maps, click-to-select on inventory rows, scroll-wheel navigation.
- `fullscreen-inventory`: Full terminal inventory overlay using alt-screen render path (declared via `tea.View.AltScreen`); layout with item list on left and ragdoll equipment outline on right; accessible via `i` from all modes.

### Modified Capabilities
- `input-navigation`: Inventory cursor and player movement extended with mouse click handling; `ScreenMode` field added to `Model` controls which input branch is active.
- `rendering-system`: `buildView` extended to dispatch to `renderFullscreenInventory` when `ScreenMode == ScreenInventory`; mouse mode enabled in `Init()`.

## Impact

- `go.mod`: `github.com/charmbracelet/bubbletea` → `charm.land/bubbletea/v2`; `github.com/charmbracelet/lipgloss` → `charm.land/lipgloss/v2`.
- `cmd/karma-vine/main.go`: `tea.NewProgram` options removed; alt-screen and mouse mode declared in `View()` instead.
- `internal/game/model.go`: `View() string` → `View() tea.View`; `case tea.KeyMsg:` → `case tea.KeyPressMsg:`.
- `internal/game/types.go`: New `ScreenMode` type and constants.
- `internal/game/model.go`: `screenMode ScreenMode` field; `showInventory bool` field replaced by `screenMode`.
- `internal/game/input.go`: `handleKey` accepts `tea.KeyPressMsg`; new `handleMouseClick(tea.MouseClickMsg, Model)` and `handleMouseWheel(tea.MouseWheelMsg, Model)` functions; `ScreenInventory` mouse callbacks.
- `internal/game/render.go`: New `renderFullscreenInventory` function; `buildView` dispatches on `screenMode`; old `inventoryPanelW` side-panel path removed.
