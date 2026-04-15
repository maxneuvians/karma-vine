package game

import (
	"strings"

	"charm.land/lipgloss/v2"
)

// portraitCell represents a single pixel in a 20×40 portrait grid.
type portraitCell struct {
	r     rune
	color string
}

// c is a shorthand constructor for portraitCell.
func c(r rune, color string) portraitCell {
	return portraitCell{r: r, color: color}
}

// bg returns an empty (space) portrait cell with no visible color.
func bg() portraitCell {
	return portraitCell{r: ' ', color: ""}
}

// Portrait dimensions.
const (
	portraitRows = 20
	portraitCols = 40
)

// row pads or truncates a slice of portraitCell to exactly portraitCols.
func row(cells ...portraitCell) [portraitCols]portraitCell {
	var r [portraitCols]portraitCell
	for i := range r {
		if i < len(cells) {
			r[i] = cells[i]
		} else {
			r[i] = bg()
		}
	}
	return r
}

// ── Player portrait: heroic humanoid silhouette ──────────────────────────────

var playerPortrait = buildPlayerPortrait()

func buildPlayerPortrait() [][portraitCols]portraitCell {
	// Colors
	skin := "#d4a574"
	hair := "#3d2b1f"
	armor := "#708090"
	armorH := "#8899aa"
	armorD := "#556677"
	boot := "#5c3317"
	cape := "#8b0000"
	capeD := "#660000"
	sword := "#c0c0c0"
	swordH := "#e0e0e0"
	bg := ""

	_ = bg
	B := func() portraitCell { return portraitCell{r: ' ', color: ""} }
	H := func(r rune) portraitCell { return c(r, hair) }
	S := func(r rune) portraitCell { return c(r, skin) }
	A := func(r rune) portraitCell { return c(r, armor) }
	AH := func(r rune) portraitCell { return c(r, armorH) }
	AD := func(r rune) portraitCell { return c(r, armorD) }
	BT := func(r rune) portraitCell { return c(r, boot) }
	C := func(r rune) portraitCell { return c(r, cape) }
	CD := func(r rune) portraitCell { return c(r, capeD) }
	SW := func(r rune) portraitCell { return c(r, sword) }
	SH := func(r rune) portraitCell { return c(r, swordH) }

	p := make([][portraitCols]portraitCell, portraitRows)

	// Row 0: top of head / hair
	p[0] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), H('░'), H('▒'), H('▓'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('▓'), H('▒'), H('░'))

	// Row 1: hair
	p[1] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), H('▒'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('█'), H('▒'))

	// Row 2: forehead + hair sides
	p[2] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), H('▒'), H('█'), H('█'), S('▓'), S('▓'), S('▓'), S('█'), S('█'), S('▓'), S('▓'), S('▓'), H('█'), H('█'), H('▒'))

	// Row 3: eyes
	p[3] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), H('▒'), H('█'), S('█'), S('█'), c('●', "#222222"), S('█'), S('█'), S('█'), S('█'), c('●', "#222222"), S('█'), S('█'), H('█'), H('▒'))

	// Row 4: nose/mouth
	p[4] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), H('▒'), S('█'), S('█'), S('█'), S('▓'), S('▒'), S('▓'), S('█'), S('█'), S('█'), S('█'), H('▒'))

	// Row 5: chin
	p[5] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▒'), S('▓'), S('█'), S('█'), S('█'), S('█'), S('█'), S('▓'), S('▒'))

	// Row 6: neck
	p[6] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▓'), S('█'), S('█'), S('█'), S('▓'))

	// Row 7: shoulders with cape and armor
	p[7] = row(B(), B(), B(), B(), B(), B(), B(), B(), C('░'), C('▒'), C('▓'), C('█'), A('▓'), A('█'), AH('█'), A('█'), A('█'), A('█'), A('█'), A('█'), A('█'), A('█'), A('█'), A('█'), AH('█'), A('█'), A('▓'), B(), B(), B(), B(), SW('│'), SH('│'))

	// Row 8: upper chest + arms
	p[8] = row(B(), B(), B(), B(), B(), B(), B(), C('▒'), C('▓'), C('█'), C('█'), A('▒'), A('▓'), A('█'), AH('█'), AH('█'), A('█'), A('█'), A('█'), A('█'), A('█'), AH('█'), AH('█'), A('█'), A('▓'), A('▒'), B(), B(), B(), B(), B(), SW('│'), SH('│'))

	// Row 9: chest plate
	p[9] = row(B(), B(), B(), B(), B(), B(), B(), C('▓'), C('█'), C('█'), AD('▒'), A('▓'), A('█'), A('█'), AH('█'), c('◆', "#ffcc00"), A('█'), A('█'), A('█'), A('█'), c('◆', "#ffcc00"), AH('█'), A('█'), A('█'), A('▓'), AD('▒'), B(), B(), B(), B(), B(), SW('│'), SH('│'))

	// Row 10: belt area
	p[10] = row(B(), B(), B(), B(), B(), B(), B(), B(), C('▓'), C('█'), AD('▓'), A('█'), A('█'), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), A('█'), A('█'), AD('▓'))

	// Row 11: upper legs with cape draping
	p[11] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), C('█'), CD('▓'), A('▓'), A('█'), A('█'), AD('▓'), AD('░'), B(), B(), AD('░'), AD('▓'), A('█'), A('█'), A('▓'), CD('▓'))

	// Row 12: mid legs
	p[12] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), C('▓'), CD('▒'), A('▓'), A('█'), AD('▓'), B(), B(), B(), B(), AD('▓'), A('█'), A('▓'), CD('▒'))

	// Row 13: lower legs
	p[13] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), CD('░'), A('▒'), A('█'), AD('▒'), B(), B(), B(), B(), AD('▒'), A('█'), A('▒'))

	// Row 14: knees
	p[14] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), A('█'), A('▓'), B(), B(), B(), B(), A('▓'), A('█'), A('▒'))

	// Row 15: calves
	p[15] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('░'), A('▓'), A('█'), B(), B(), B(), B(), A('█'), A('▓'), A('░'))

	// Row 16: shins
	p[16] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), A('▓'), B(), B(), B(), B(), A('▓'), A('▒'))

	// Row 17: ankles
	p[17] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BT('▒'), BT('▓'), B(), B(), B(), B(), BT('▓'), BT('▒'))

	// Row 18: boots
	p[18] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BT('▒'), BT('▓'), BT('█'), BT('█'), B(), B(), BT('█'), BT('█'), BT('▓'), BT('▒'))

	// Row 19: boot soles
	p[19] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BT('░'), BT('▓'), BT('█'), BT('█'), BT('█'), B(), B(), BT('█'), BT('█'), BT('█'), BT('▓'), BT('░'))

	return p
}

// ── Enemy portraits ──────────────────────────────────────────────────────────

var humanoidPortrait = buildHumanoidPortrait()

func buildHumanoidPortrait() [][portraitCols]portraitCell {
	skin := "#7a8b5c"
	armor := "#5a4a3a"
	armorH := "#6b5b4b"
	eye := "#ff3300"
	weapon := "#888888"

	B := func() portraitCell { return portraitCell{r: ' ', color: ""} }
	S := func(r rune) portraitCell { return c(r, skin) }
	A := func(r rune) portraitCell { return c(r, armor) }
	AH := func(r rune) portraitCell { return c(r, armorH) }
	E := func(r rune) portraitCell { return c(r, eye) }
	W := func(r rune) portraitCell { return c(r, weapon) }

	p := make([][portraitCols]portraitCell, portraitRows)

	p[0] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('░'), S('▒'), S('▓'), S('█'), S('█'), S('▓'), S('▒'), S('░'))
	p[1] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▒'), S('█'), S('█'), S('█'), S('█'), S('█'), S('█'), S('█'), S('▒'))
	p[2] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▒'), S('█'), S('█'), E('●'), S('█'), S('█'), S('█'), E('●'), S('█'), S('█'), S('▒'))
	p[3] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▓'), S('█'), S('█'), S('▓'), S('▒'), S('▓'), S('█'), S('█'), S('▓'))
	p[4] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▒'), S('▓'), S('█'), S('█'), S('▓'), S('▒'))
	p[5] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), S('▓'), S('▓'))
	p[6] = row(B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), A('▓'), A('█'), A('█'), A('█'), A('█'), AH('█'), A('█'), A('█'), A('█'), A('█'), A('█'), A('█'), AH('█'), A('█'), A('▓'), A('▒'))
	p[7] = row(B(), B(), B(), B(), B(), B(), B(), A('░'), A('▓'), A('█'), A('█'), A('█'), A('█'), A('█'), AH('█'), AH('█'), A('█'), A('█'), A('█'), A('█'), AH('█'), AH('█'), A('█'), A('█'), A('▓'), A('░'))
	p[8] = row(B(), B(), B(), B(), B(), B(), S('▒'), S('▓'), A('█'), A('█'), A('█'), AH('█'), A('█'), A('█'), A('█'), A('█'), A('█'), A('█'), AH('█'), A('█'), A('█'), A('█'), S('▓'), S('▒'))
	p[9] = row(B(), B(), B(), B(), B(), B(), B(), S('▓'), A('▓'), A('█'), A('█'), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), A('█'), A('█'), A('▓'), S('▓'), B(), B(), B(), B(), B(), B(), B(), W('╱'))
	p[10] = row(B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), A('█'), A('█'), A('▓'), A('░'), B(), B(), A('░'), A('▓'), A('█'), A('█'), A('▒'), B(), B(), B(), B(), B(), B(), B(), W('│'))
	p[11] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▓'), A('█'), A('▒'), B(), B(), B(), B(), A('▒'), A('█'), A('▓'), B(), B(), B(), B(), B(), B(), B(), W('│'))
	p[12] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), A('█'), A('░'), B(), B(), B(), B(), A('░'), A('█'), A('▒'))
	p[13] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▓'), A('█'), B(), B(), B(), B(), A('█'), A('▓'))
	p[14] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), A('▓'), B(), B(), B(), B(), A('▓'), A('▒'))
	p[15] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('░'), A('▓'), B(), B(), B(), B(), A('▓'), A('░'))
	p[16] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), A('▒'), B(), B(), B(), B(), A('▒'))
	p[17] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), c('▒', "#5c3317"), c('▓', "#5c3317"), B(), B(), B(), B(), c('▓', "#5c3317"), c('▒', "#5c3317"))
	p[18] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), c('▒', "#5c3317"), c('▓', "#5c3317"), c('█', "#5c3317"), B(), B(), B(), B(), c('█', "#5c3317"), c('▓', "#5c3317"), c('▒', "#5c3317"))
	p[19] = row(B(), B(), B(), B(), B(), B(), B(), B(), c('░', "#5c3317"), c('▓', "#5c3317"), c('█', "#5c3317"), c('█', "#5c3317"), B(), B(), B(), B(), c('█', "#5c3317"), c('█', "#5c3317"), c('▓', "#5c3317"), c('░', "#5c3317"))

	return p
}

var beastPortrait = buildBeastPortrait()

func buildBeastPortrait() [][portraitCols]portraitCell {
	fur := "#8b6914"
	furD := "#6b4f10"
	furL := "#a0801a"
	eye := "#ffcc00"
	nose := "#333333"
	claw := "#cccccc"

	B := func() portraitCell { return portraitCell{r: ' ', color: ""} }
	F := func(r rune) portraitCell { return c(r, fur) }
	FD := func(r rune) portraitCell { return c(r, furD) }
	FL := func(r rune) portraitCell { return c(r, furL) }
	E := func(r rune) portraitCell { return c(r, eye) }
	N := func(r rune) portraitCell { return c(r, nose) }
	CL := func(r rune) portraitCell { return c(r, claw) }

	p := make([][portraitCols]portraitCell, portraitRows)

	// Four-legged beast (wolf/bear silhouette)
	p[0] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('▓'), F('█'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('░'), F('▒'))
	p[1] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('█'), F('█'), F('█'), F('▓'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('▓'), F('█'))
	p[2] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('█'), F('█'), E('●'), F('█'), F('█'), F('▓'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▓'), F('█'), F('█'), F('█'))
	p[3] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▓'), F('█'), F('█'), F('█'), F('█'), N('▓'), N('█'), N('▒'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('█'), F('█'), F('█'), F('█'), F('▒'))
	p[4] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▓'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('▓'), F('▒'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('▓'), F('█'), F('█'), F('█'), F('▓'))
	p[5] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▓'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('▓'), F('▒'), F('▒'), F('▓'), F('▓'), F('▒'), F('▒'), F('▓'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('▓'))
	p[6] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FL('▒'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), FL('▒'))
	p[7] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FL('▓'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), FL('▓'))
	p[8] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('▓'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), F('█'), FD('▓'))
	p[9] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('▒'), F('█'), F('█'), F('█'), F('█'), FD('▓'), FD('░'), B(), B(), FD('░'), FD('▓'), F('█'), F('█'), F('█'), F('█'), F('█'), FD('▒'))
	p[10] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▓'), F('█'), F('▓'), B(), B(), B(), B(), B(), B(), B(), B(), F('▓'), F('█'), F('█'), F('▓'))
	p[11] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('█'), F('▒'), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('█'), F('█'), F('▒'))
	p[12] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('░'), F('▓'), F('░'), B(), B(), B(), B(), B(), B(), B(), B(), F('░'), F('▓'), F('▓'), F('░'))
	p[13] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('▒'))
	p[14] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), F('▒'), F('▒'))
	p[15] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('▒'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('▒'), FD('▒'))
	p[16] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('░'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('░'), FD('░'))
	p[17] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('░'), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), FD('░'), FD('░'))
	p[18] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), CL('▒'), CL('▓'), CL('▒'), B(), B(), B(), B(), B(), B(), B(), B(), CL('▒'), CL('▓'), CL('▓'), CL('▒'))
	p[19] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), CL('░'), CL('▓'), CL('█'), CL('▓'), CL('░'), B(), B(), B(), B(), B(), B(), CL('░'), CL('▓'), CL('█'), CL('█'), CL('▓'), CL('░'))

	return p
}

var undeadPortrait = buildUndeadPortrait()

func buildUndeadPortrait() [][portraitCols]portraitCell {
	bone := "#d4c9a8"
	boneD := "#a89878"
	glow := "#44ff88"
	robe := "#2a2a3a"
	robeD := "#1a1a2a"
	robeL := "#3a3a4a"

	B := func() portraitCell { return portraitCell{r: ' ', color: ""} }
	BN := func(r rune) portraitCell { return c(r, bone) }
	BD := func(r rune) portraitCell { return c(r, boneD) }
	G := func(r rune) portraitCell { return c(r, glow) }
	R := func(r rune) portraitCell { return c(r, robe) }
	RD := func(r rune) portraitCell { return c(r, robeD) }
	RL := func(r rune) portraitCell { return c(r, robeL) }

	p := make([][portraitCols]portraitCell, portraitRows)

	// Skeletal/spectral figure
	p[0] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('░'), BN('▒'), BN('▓'), BN('█'), BN('█'), BN('▓'), BN('▒'), BN('░'))
	p[1] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('▒'), BN('█'), BN('█'), BN('█'), BN('█'), BN('█'), BN('█'), BN('█'), BN('▒'))
	p[2] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('▒'), BN('█'), BD('▓'), G('●'), BD('░'), BD('░'), BD('░'), G('●'), BD('▓'), BN('█'), BN('▒'))
	p[3] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('▓'), BN('█'), BD('▒'), BD('▓'), BD('█'), BD('▓'), BD('▒'), BN('█'), BN('▓'))
	p[4] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('▒'), BN('▓'), BD('█'), BD('█'), BN('▓'), BN('▒'))
	p[5] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BD('▒'), BD('▒'))
	p[6] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[7] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), RL('█'), RL('█'), R('█'), R('█'), RL('█'), RL('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[8] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('▒'), BD('▓'), R('█'), R('█'), R('█'), RL('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), RL('█'), R('█'), R('█'), BD('▓'), BN('▒'))
	p[9] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), BN('▓'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), BN('▓'))
	p[10] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), RD('▒'), R('█'), RD('▓'), RD('░'), B(), B(), RD('░'), RD('▓'), R('█'), RD('▒'))
	p[11] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('▓'), R('█'), R('▓'), R('▒'), R('▒'), R('▓'), R('█'), R('▓'))
	p[12] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), RD('▒'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), RD('▒'))
	p[13] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('░'))
	p[14] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[15] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[16] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[17] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[18] = row(B(), B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))
	p[19] = row(B(), B(), B(), B(), B(), B(), B(), R('░'), R('▒'), R('▓'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('█'), R('▓'), R('▒'), R('░'))

	return p
}

// ── Fallback portrait ────────────────────────────────────────────────────────

var fallbackPortrait = buildFallbackPortrait()

func buildFallbackPortrait() [][portraitCols]portraitCell {
	body := "#888888"
	bodyD := "#666666"
	eye := "#ff4444"

	B := func() portraitCell { return portraitCell{r: ' ', color: ""} }
	M := func(r rune) portraitCell { return c(r, body) }
	MD := func(r rune) portraitCell { return c(r, bodyD) }
	E := func(r rune) portraitCell { return c(r, eye) }

	p := make([][portraitCols]portraitCell, portraitRows)

	// Generic blob creature
	p[0] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[1] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[2] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[3] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), E('●'), M('█'), M('█'), M('█'), E('●'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[4] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[5] = row(B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[6] = row(B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[7] = row(B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[8] = row(B(), B(), B(), B(), B(), B(), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'))
	p[9] = row(B(), B(), B(), B(), B(), B(), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'))
	p[10] = row(B(), B(), B(), B(), B(), B(), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'))
	p[11] = row(B(), B(), B(), B(), B(), B(), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'))
	p[12] = row(B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[13] = row(B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[14] = row(B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), M('█'), M('█'), MD('▓'), MD('▒'), MD('░'), MD('░'), MD('▒'), MD('▓'), M('█'), M('█'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[15] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), M('▓'), M('█'), MD('▓'), B(), B(), B(), B(), B(), B(), MD('▓'), M('█'), M('█'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[16] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), M('▒'), MD('▓'), B(), B(), B(), B(), B(), B(), B(), B(), MD('▓'), M('█'), M('█'), M('▓'), M('▒'), M('░'))
	p[17] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), M('░'), MD('▒'), B(), B(), B(), B(), B(), B(), B(), B(), MD('▒'), M('▓'), M('▓'), M('▒'), M('░'))
	p[18] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), MD('░'), B(), B(), B(), B(), B(), B(), B(), B(), B(), MD('▒'), MD('▒'), MD('░'))
	p[19] = row(B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), B(), MD('░'), MD('░'))

	return p
}

// ── Portrait selection ───────────────────────────────────────────────────────

// enemyPortrait returns the portrait grid for the given enemy character rune.
func enemyPortrait(char rune) [][portraitCols]portraitCell {
	switch char {
	// Humanoid archetypes
	case 'H', 'K', 'G', 'T', 'O':
		return humanoidPortrait
	// Beast archetypes
	case 'W', 'B', 'S', 'R', 'D':
		return beastPortrait
	// Undead archetypes
	case 'Z', 'V', 'L', 'M':
		return undeadPortrait
	default:
		return fallbackPortrait
	}
}

// ── Portrait rendering ───────────────────────────────────────────────────────

// renderPortrait renders a portrait grid as a styled string.
// Each cell is rendered with its foreground colour. Rows are clipped to panelWidth.
func renderPortrait(p [][portraitCols]portraitCell, panelWidth int) string {
	if panelWidth <= 0 {
		return ""
	}
	cols := portraitCols
	if panelWidth < cols {
		cols = panelWidth
	}

	rows := make([]string, len(p))
	for i, pRow := range p {
		var b strings.Builder
		for j := 0; j < cols; j++ {
			cell := pRow[j]
			if cell.color == "" {
				b.WriteRune(cell.r)
			} else {
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(cell.color)).Render(string(cell.r)))
			}
		}
		rows[i] = b.String()
	}
	return strings.Join(rows, "\n")
}
