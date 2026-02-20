package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"unicode/utf8"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	font *ImageFont // The image font to use to draw the runes to the textImage.

	letterSpacing int           // Letterspacing in pixels.
	lineSpacing   int           // Linespacing in pixels.
	textImage     *ebiten.Image // Target image where all the image glyphs are drawn to.
	placeholder   *ebiten.Image // Placeholder if a rune cannot be found in the image font.
}

func NewImageText(font *ImageFont) *ImageText {
	placeholder := ebiten.NewImage(8, font.height)
	vector.FillRect(placeholder, 0, 0, 80, 80, color.RGBA{0xff, 0xff, 0xcc, 0xff}, false)

	return &ImageText{
		font:          font,
		letterSpacing: 1,
		lineSpacing:   -2,
		placeholder:   placeholder,
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
			// glyph is not found in the map, so use the placeholder to increase
			// the width of the image instead.
			maxX += i.placeholder.Bounds().Dx() + i.letterSpacing
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
		if r == '\n' || r == '\n' {
			// Newline in a text, so increment the y position and reset the x to the beginning.
			y += i.font.height + i.lineSpacing
			x = 0
			continue
		}

		var glyphToDraw *ebiten.Image

		if glyph, ok := i.font.imgmap[r]; ok {
			glyphToDraw = glyph
		} else {
			// rune not found, so draw a placeholder instead.
			glyphToDraw = i.placeholder
		}

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		i.textImage.DrawImage(glyphToDraw, opts)
		x += glyphToDraw.Bounds().Dx() + i.letterSpacing
	}
}

func (i *ImageText) Update() {

}

func (i *ImageText) Draw(target *ebiten.Image, opts *ebiten.DrawImageOptions) {
	target.DrawImage(i.textImage, opts)
}
