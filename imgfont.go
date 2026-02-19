package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageFont struct {
	glyphs string
	Height int
	// todo rune
	imgmap map[byte](*ebiten.Image)
}

func NewImageFont(file, glyphs string) (*ImageFont, error) {
	font := &ImageFont{
		glyphs: glyphs,
		imgmap: make(map[byte](*ebiten.Image)),
	}

	derp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	decodedImage, err := png.Decode(derp)
	if err != nil {
		return nil, fmt.Errorf("could not decode file %s as PNG: %v", file, err)
	}

	font.Height = decodedImage.Bounds().Dy()

	separatorColor := decodedImage.At(0, 0)
	glyphIndex := 0
	glyphWidth := 0
	lastSeparatorIndex := 0

	for x := 0; x < decodedImage.Bounds().Dx(); x++ {
		// Check if the pixel at the current x position is our designated separator color.
		// If so, check if it's time to assign a glyph.
		if decodedImage.At(x, 0) == separatorColor {
			if glyphWidth > 0 {
				// fmt.Printf("Glyph at %d, size: %d\n", lastSeparatorIndex+1, glyphWidth)
				rect := image.Rectangle{
					Min: image.Pt(lastSeparatorIndex+1, 0),
					Max: image.Pt(lastSeparatorIndex+1+glyphWidth, decodedImage.Bounds().Dy()),
				}

				subImage := decodedImage.(*image.NRGBA).SubImage(rect)
				bleh := ebiten.NewImageFromImage(subImage)
				// todo: check if glyphindex is larger than string length, or else panic
				font.imgmap[glyphs[glyphIndex]] = bleh

				glyphIndex++
			}

			lastSeparatorIndex = x
			glyphWidth = 0
			continue
		}

		// No separator color, so we can safely increment the current glyph's width.
		glyphWidth++
	}

	if len(glyphs) != len(font.imgmap) {
		return nil, fmt.Errorf("%d glyphs wanted, but image contained %d glyphs", len(glyphs), len(font.imgmap))
	}

	return font, nil
}

type ImageText struct {
	font *ImageFont

	TextImage *ebiten.Image
}

func NewImageText(font *ImageFont) *ImageText {
	return &ImageText{
		font: font,
	}
}

func (i *ImageText) SetText(text string) {
	// TODO: determine max image size before?
	i.TextImage = ebiten.NewImage(320, 240)

	x := 0.0
	y := 0.0
	for _, r := range text {
		if r == '\n' {
			y += float64(i.font.Height)
			x = 0
			continue
		}

		glyph := i.font.imgmap[byte(r)]
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, y)
		i.TextImage.DrawImage(glyph, opts)

		x += float64(glyph.Bounds().Dx() + 1)

	}

	// calculate max width of image based on string
	// create image
}

func (i *ImageText) Update() {

}

func (i *ImageText) Draw(target *ebiten.Image, opts *ebiten.DrawImageOptions) {
	target.DrawImage(i.TextImage, opts)
}
