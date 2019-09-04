package sprite

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

var (
	zero  = uint16(0)
	max   = uint16(65535)
	red   = color.RGBA64{R: max, G: zero, B: zero, A: max}
	white = color.RGBA64{R: max, G: max, B: max, A: max}
)

func TestSprite(t *testing.T) {
	t.Run("copying source image data", func(t *testing.T) {
		src := image.NewRGBA64(image.Rect(0, 0, 2, 2))
		src.Set(0, 0, white)
		src.Set(0, 1, red)
		src.Set(1, 0, red)
		src.Set(1, 1, white)

		var srcAsImg image.Image = src
		s := NewSprite(&srcAsImg)

		assertEquals(white, s.img.At(0, 0), t)
		assertEquals(red, s.img.At(0, 1), t)
		assertEquals(red, s.img.At(1, 0), t)
		assertEquals(white, s.img.At(1, 1), t)
	})

	t.Run("copying source image with offset", func(t *testing.T) {
		src := image.NewRGBA64(image.Rect(0, 0, 3, 2))
		src.Set(0, 0, white)
		src.Set(0, 1, red)
		src.Set(1, 0, red)
		src.Set(1, 1, white)
		src.Set(2, 0, white)
		src.Set(2, 1, red)

		var srcAsImg = src.SubImage(image.Rect(1, 0, 3, 2))
		s := NewSprite(&srcAsImg)

		assertEquals(red, s.img.At(0, 0), t)
		assertEquals(white, s.img.At(0, 1), t)
		assertEquals(white, s.img.At(1, 0), t)
		assertEquals(red, s.img.At(1, 1), t)
	})

	t.Run("instantiating sprite from png", func(t *testing.T) {
		file, err := os.Open("./alien.png")
		if err != nil {
			fmt.Println("img.jpg file not found!")
			os.Exit(1)
		}
		defer file.Close()

		img, err := png.Decode(file)
		if err != nil {
			fmt.Println("could not decode image!")
			os.Exit(1)
		}

		NewSprite(&img)
	})
}

func assertEquals(c1, c2 color.Color, t *testing.T) {
	if !sameColor(c1, c2) {
		fmt.Printf("%v is not the same as %v\n", c1, c2)
		t.Fail()
	}
}

func sameColor(a, b color.Color) bool {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}
