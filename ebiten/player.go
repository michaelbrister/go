package main

import "github.com/hajimehoshi/ebiten/v2"

func drawLayer(screen *ebiten.Image) error {
	playerOne.movePlayer()

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	playerOp := &ebiten.DrawImageOptions{}
	playerOp.GeoM.Translate(playerOne.xPos, playerOne.yPos)
	screen.DrawImage(playerOne.image, playerOp)

	return nil
}

func (p *playerClass) movePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.yPos -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.yPos += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.xPos -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.xPos += p.speed
	}
}
