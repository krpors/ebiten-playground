package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

// State defines the interface for a game state.
type State interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type Game struct {
	currentState State
}

func NewGame() *Game {
	game := &Game{}
	game.currentState = &StateTest{}
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
