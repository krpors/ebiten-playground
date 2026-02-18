package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type StateTest struct {
	Radians float32
}

func (s *StateTest) Update() error {
	s.Radians += 0.5 / float32(ebiten.TPS())
	return nil
}

func (s *StateTest) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(320, 240)
	vector.FillRect(img, 0, 0, 80, 80, color.RGBA{0xff, 0xff, 0xcc, 0x22}, true)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-40, -40)
	op.GeoM.Rotate(float64(s.Radians))
	op.GeoM.Translate(40, 40)
	screen.DrawImage(img, op)

	tri := ebiten.NewImage(320, 240)
	path := vector.Path{}
	path.LineTo(0, 0)
	path.LineTo(40, 0)
	path.LineTo(20, 60)
	vector.FillPath(tri, &path, &vector.FillOptions{}, &vector.DrawPathOptions{})

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(path.Bounds().Dx())/2, -float64(path.Bounds().Dy())/2)
	op.GeoM.Rotate(-float64(s.Radians))
	op.GeoM.Translate(90, 90)
	op.ColorScale.SetR(1.0)
	op.ColorScale.SetG(0)
	op.ColorScale.SetB(0)

	screen.DrawImage(tri, op)
}
