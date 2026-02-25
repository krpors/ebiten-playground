package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type StateArt struct {
	Image  *ebiten.Image
	Ticks  int
	Sample int
}

func NewStateArt() *StateArt {
	s := &StateArt{}

	s.Image = ebiten.NewImage(1024, 768)

	// ebitenutil.DebugPrintAt(s.Image, "HERRO!", 12, 12)

	var inc int = 20
	for x := 0; x < s.Image.Bounds().Dx(); x += 10 {
		for y := 0; y < s.Image.Bounds().Dy(); y += inc + 2 {
			inc = rand.Intn(20) + 5
			clr := color.RGBA{
				R: uint8(rand.Intn(0xff)),
				G: uint8(rand.Intn(0xff)),
				B: uint8(rand.Intn(0xff)),
			}
			vector.FillRect(s.Image, float32(x), float32(y), 8, float32(inc), clr, false)
		}
	}

	s.Sample = 0

	return s
}

func (s *StateArt) sample() int {
	s.Sample++
	argh := math.Mod(float64(s.Sample)*math.Phi, 1.0)
	// fmt.Println(argh * 255)
	return int(argh * 255)
}

func (s *StateArt) Update() error {
	s.Ticks++

	if s.Ticks%60 == 0 {
		s.Image.Clear()
		var inc int = 20
		for x := 0; x < s.Image.Bounds().Dx(); x += 10 {
			for y := 0; y < s.Image.Bounds().Dy(); y += inc + 2 {
				inc = rand.Intn(5) + 5
				clr := color.RGBA{
					R: uint8(s.sample()),
					G: uint8(rand.Intn(255)),
					B: uint8(rand.Intn(255)),
				}
				vector.FillRect(s.Image, float32(x), float32(y), 8, float32(inc), clr, false)
			}
		}
	}

	return nil
}

func (s *StateArt) Draw(screen *ebiten.Image) {
	screen.DrawImage(s.Image, &ebiten.DrawImageOptions{})
}
