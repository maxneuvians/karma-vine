## Context

The fullscreen inventory screen (`ScreenInventory`) has two columns: a left item list (driven by `inventory.Items []Item` and `inventoryCursor int`) and a right ragdoll panel that currently hardcodes `[ Empty ]` for all six slots. The `equipSlots` string slice and the ragdoll ASCII art are already defined in `render.go`. The `Item` struct has `Char`, `Color`, `Name`, `Count` but no slot eligibility information.

There is no current mechanism to distinguish wearable items (clothing, weapons) from stackable consumables (torches), nor any storage for what the player has equipped.

## Goals / Non-Goals

**Goals:**
- Define a `BodySlot` type so items can declare which slots they fit
- Store equipped items in `Inventory.Equipped` (one item per slot)
- Let the player equip from the inventory list and unequip from the ragdoll panel, all within `ScreenInventory`
- Render the ragdoll column with equipped item names
- Start with a default outfit so the slots aren't all empty from the beginning

**Non-Goals:**
- Stat bonuses or damage modifiers from equipped items (no combat system yet)
- Weight or encumbrance
- Two-handed weapon logic (occupying both hand slots simultaneously)
- Durability or item degradation
- Drag-and-drop mouse equip (keyboard-only for now)

## Decisions

### 1. `Inventory.Equipped` as a fixed-size array, not a map

**Decision:** `type Equipped [NumBodySlots]Item` (or a named array type) where the index is a `BodySlot` constant. An empty slot is represented by `Item{}` (zero value, `Name == ""`).

**Rationale:** A fixed array of size 6 is safe to index directly by `BodySlot` constant without nil checks. The zero-value sentinel (`Name == ""`) avoids a separate "slot is empty" bool. A `map[BodySlot]Item` would require a nil-check every render and adds allocation overhead for no benefit.

**Alternative considered:** `map[BodySlot]*Item`. Rejected — pointer indirection, GC pressure, more complex zero-check.

### 2. `Item.Slots []BodySlot` for eligibility

**Decision:** Each `Item` carries a `Slots []BodySlot` slice listing which body slots it can occupy. Empty slice = not equippable.

**Rationale:** This keeps all item metadata on the item itself, making it trivial to define new equippable items anywhere item construction happens (pickup, world gen, default outfit). No external lookup table required.

**Alternative considered:** A `map[string][]BodySlot` keyed by item name. Rejected — requires central registration and becomes a maintenance burden as more item types are added.

### 3. Two-cursor focus system: `equipFocused bool` + `equipCursor int`

**Decision:** Add `equipFocused bool` and `equipCursor int` to `Model`. In `ScreenInventory`, `Tab` toggles `equipFocused`. When `equipFocused == false`, `up`/`down` drive `inventoryCursor`; when `equipFocused == true`, they drive `equipCursor` (clamped `[0, NumBodySlots-1]`). The `e` key meaning is context-sensitive: equip when left-focused, unequip when right-focused.

**Rationale:** Separating focus state is cleaner than a single combined cursor that switches column at the edges. It also makes the UI intent explicit in the render (active header highlight) and avoids surprise cursor jumps.

**Alternative considered:** Single cursor that overflows into ragdoll slots (items 0–7 are inventory, slots 8–13 are ragdoll). Rejected — confusing UX and harder to render cleanly.

### 4. Equip action: best-slot selection, swap on conflict

**Decision:** When the player presses `e` with `equipFocused == false`:
1. Take `item = inventory.Items[inventoryCursor]`
2. If `len(item.Slots) == 0`, no-op (not equippable)
3. Iterate `item.Slots` in order; find the first slot where `Equipped[slot].Name == ""` (empty) → equip there
4. If no empty slot exists among `item.Slots`, use `item.Slots[0]` → swap: current equipped goes back to inventory, new item takes the slot
5. Remove (decrement count / remove slot) from `inventory.Items`

When `equipFocused == true` and the player presses `e`:
1. Take `slot = BodySlot(equipCursor)`
2. If `Equipped[slot].Name == ""`, no-op
3. Move `Equipped[slot]` back to inventory (stack if same name, or add new slot)
4. Clear `Equipped[slot]` to `Item{}`

**Rationale:** Preferring empty slots avoids unintended swaps. Falling back to slot 0 on conflict is predictable. The unequip path is symmetric and easy to test.

### 5. Equipped items are separate from the carry count

**Decision:** `inventory.Items` holds unequipped items only. The pickup cap (`InventoryMaxSlots = 8`) applies only to `Items`; equipped items do not count toward it.

**Rationale:** Treating equipped gear as occupying inventory slots would make the 8-slot limit feel punishing and confuse the inventory/ragdoll mental model. Equipped = worn, not carried.

### 6. Default outfit: three items pre-equipped in `NewModel()`

**Decision:** `NewModel()` calls a `defaultOutfit()` helper that returns an `Equipped` array with:
- `Chest`: `Item{Char: '♦', Color: "#a0a0a0", Name: "Cloth Tunic", Slots: []BodySlot{SlotChest}}`
- `Legs`: `Item{Char: '‖', Color: "#a0a0a0", Name: "Cloth Pants", Slots: []BodySlot{SlotLegs}}`
- `Feet`: `Item{Char: '∩', Color: "#8B4513", Name: "Leather Boots", Slots: []BodySlot{SlotFeet}}`

All other slots start empty.

**Rationale:** Three items give a sense of a dressed character without overwhelming the player. Head and hands are left open as interesting equip targets early in the game. A dedicated helper keeps `NewModel()` readable.

## Risks / Trade-offs

- **Existing render tests**: `render_test.go` currently asserts `[ Empty ]` for all ragdoll slots (or does not test equipped rendering). Tests will need updating for the default outfit. → Update affected tests as part of this change.
- **`e` key collision**: `e` is currently unbound globally. If future changes bind `e` for another action, this will conflict. → Document the binding.
- **`Tab` in terminal**: Some terminal emulators capture `Tab` before BubbleTea sees it. BubbleTea v2 delivers `tea.KeyPressMsg{Code: tea.KeyTab}` reliably in alt-screen mode, which this game already uses. → Low risk; verify in manual testing.
