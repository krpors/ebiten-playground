package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageFont struct {
	glyphs string
	height int
	imgmap map[rune](*ebiten.Image)
}

func NewImageFont(file, glyphs string) (*ImageFont, error) {
	font := &ImageFont{
		glyphs: glyphs,
		imgmap: make(map[rune](*ebiten.Image)),
	}

	derp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	decodedImage, err := png.Decode(derp)
	if err != nil {
		return nil, fmt.Errorf("could not decode file %s as PNG: %v", file, err)
	}

	font.height = decodedImage.Bounds().Dy()

	separatorColor := decodedImage.At(0, 0)
	glyphIndex := 0
	glyphWidth := 0
	lastSeparatorIndex := 0

	var images []*ebiten.Image

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
				eimage := ebiten.NewImageFromImage(subImage)

				images = append(images, eimage)

				glyphIndex++
			}

			lastSeparatorIndex = x
			glyphWidth = 0
			continue
		}

		// No separator color, so we can safely increment the current glyph's width.
		glyphWidth++
	}

	if utf8.RuneCountInString(glyphs) != len(images) {
		return nil, fmt.Errorf("%d glyphs (runes) wanted, but image contained %d glyphs", len(glyphs), len(images))
	}

	counter := 0
	for _, r := range glyphs {
		font.imgmap[r] = images[counter]
		counter++
	}

	return font, nil
}

type ImageText struct {
	font *ImageFont

	letterSpacing int
	lineSpacing   int
	textImage     *ebiten.Image
}

func NewImageText(font *ImageFont) *ImageText {
	return &ImageText{
		font:          font,
		letterSpacing: 1,
		lineSpacing:   -2,
	}
}

// calculateImageSize calculates the maximum size which is required to draw
// all the glyphs on a target image.
func (i *ImageText) calculateImageSize(text string) (width int, height int) {
	maxX, maxY := 0, i.font.height
	for _, r := range text {
		if r == '\n' {
			maxY += i.font.height + i.lineSpacing
			continue
		}

		if glyph, ok := i.font.imgmap[r]; ok {
			maxX += glyph.Bounds().Dx() + i.letterSpacing
		} else {
			fmt.Printf("Unable to calculate image size correctly, because the rune %s is not in the image map\n", string(r))
		}
	}
	return maxX, maxY
}

func (i *ImageText) SetText(text string) {
	width, height := i.calculateImageSize(text)
	i.textImage = ebiten.NewImage(width, height)

	x := 0
	y := 0
	for _, r := range text {
		if r == '\n' {
			y += i.font.height + i.lineSpacing
			x = 0
			continue
		}

		if glyph, ok := i.font.imgmap[r]; ok {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(x), float64(y))
			i.textImage.DrawImage(glyph, opts)

			x += glyph.Bounds().Dx() + i.letterSpacing
		} else {
			fmt.Printf("Rune %s not found?\n", string(r))
			// TODO glyph not found, now what? probably draw some placeholder thing or some stuff
		}

	}
}

func (i *ImageText) Update() {

}

func (i *ImageText) Draw(target *ebiten.Image, opts *ebiten.DrawImageOptions) {
	target.DrawImage(i.textImage, opts)
}
