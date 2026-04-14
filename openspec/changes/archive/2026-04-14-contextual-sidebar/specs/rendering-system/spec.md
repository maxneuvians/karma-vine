## MODIFIED Requirements

### Requirement: View composition includes map picker overlay when active
The system SHALL, in `buildView`, check `m.showMapPicker`. When true, it SHALL render `renderMapPicker(m, mapH)` on the right side of the viewport (using the same composition approach as `showSidebar`), reducing the available map width by the picker panel width. When both `showMapPicker` and `showSidebar` are false, layout is unchanged from the current implementation.

#### Scenario: Map picker reduces map width when open
- **WHEN** `showMapPicker == true` and the viewport is 80×24
- **THEN** the rendered map occupies fewer than 80 columns (picker panel takes the remainder)

#### Scenario: No layout change when picker is closed
- **WHEN** `showMapPicker == false` and `showSidebar == false`
- **THEN** the map renders at full viewport width (minus HUD rows), identical to current behavior

## ADDED Requirements

### Requirement: renderSidebar dispatches on active mode
The system SHALL update `renderSidebar` to switch on `m.mode`, rendering world content for `ModeWorld`, local content for `ModeLocal`, and dungeon content for `ModeDungeon`. The `localCharNames` lookup map SHALL be removed; object and animal names SHALL be read directly from `obj.Name` and `a.Name`.

#### Scenario: renderSidebar called in ModeDungeon returns dungeon content
- **WHEN** `renderSidebar` is called with `m.mode == ModeDungeon`
- **THEN** the returned string contains `Dungeon` and `Depth:`

#### Scenario: renderSidebar called in ModeLocal uses Name field
- **WHEN** `renderSidebar` is called with `m.mode == ModeLocal` and an object has `Name: "Dungeon Entrance"`
- **THEN** the returned string contains `Dungeon Entrance`
