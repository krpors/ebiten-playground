package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Radians float32
}

func (g *Game) Update() error {
	g.Radians += 1 / float32(ebiten.TPS())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(320, 240)
	vector.FillRect(img, 0, 0, 80, 80, color.RGBA{0xff, 0xff, 0xcc, 0x22}, true)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-40, -40)
	op.GeoM.Rotate(float64(g.Radians))
	op.GeoM.Translate(40, 40)
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
