package main

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// State defines the interface for a game state.
type State interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	currentState State
	text         *ImageText
}

func NewGame() *Game {
	game := &Game{}
	game.currentState = &StateTest{}

	font, err := NewImageFont("font.png", " abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-=_+|/\\:;'\"<>,.?")
	if err != nil {
		panic(err)
	}

	game.text = NewImageText(font)
	game.text.SetText("The quick brown foxeh...\n... jumps over the lazy doggeh")

	return game
}

func (g *Game) Update() error {
	if g.currentState != nil {
		return g.currentState.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentState != nil {
		g.currentState.Draw(screen)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.SetR(rand.Float32())
	opts.ColorScale.SetG(rand.Float32())
	opts.ColorScale.SetB(rand.Float32())
	opts.ColorScale.SetA(1)
	opts.GeoM.Translate(4*rand.Float64(), 4*rand.Float64())
	// opts.GeoM.Scale(rand.Float64()*2+1, rand.Float64()*2+1)
	opts.GeoM.Translate(320/2-50, 240/2-50)

	g.text.Draw(screen, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
