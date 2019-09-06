package main

import (
	"fmt"
	"image"
)

type SimpleSprite struct {
	x, y           float32
	rotation       float32
	scaleX, scaleY float32
	img            *image.RGBA
}

// TODO come up with a design that supports procedural sprites, tiled backgrounds etc.

/*
  BASIC ALGORITHM

  for each pixel of the screen determine the world coordinate
  find sprites that overlap with world coordinate
  ask each overlapping sprite the RGBA value of that world coordinate
  combine the colors in reverse z order of sprites

  O(resolution*(n*b + n_vis*A))

  CAMERA BASED APPROACH

  determine the sprites visible in the view
  ask the sprites to draw themselves on the screen in reverse z order
  sprite.Draw() must take the destination image and the camera transform

  O(n*b + n_vis*res*a)

  rough estimates: 100s of sprites on screen, 100s off-screen, 10% overlapping

  TILE BASED APPROACH

  like the camera based approach, but with many smaller cameras
  determine sprites visible on each tile
  for each tile calculate the tile transform and ask the corresponding sprites to render the tile

  O(n*b + n_vis_t*res*a)
*/

/*
  types of sprites (assuming camera/tile based approach)

  normal sprite - combines the camera transform with own transform and copies pixels from own image
  interpolated sprite - uses some interpolation method to determine inter-pixel color when scaled
  instanced sprite - like normal, but reuses image through pointer (can only be created with RGBA images)
  screen sprite - determines inital point on screen/tile and copies image data (no transform)
  composite sprite - uses a bunch of other sprites to draw various parts
  tiled sprite - like a tiled background
  procedural sprite - uses some algorithm like perlin noise to determine color
*/

type rect struct{ x1, x2, y1, y2 float32 }

func (r1 *rect) overlap(r2 rect) bool {
	return !(r1.x2 < r2.x1 || r1.y2 < r2.y1 || r2.x2 < r1.x1 || r2.y2 < r1.y1)
}

type Transform struct {
	offsetX, offsetY float32
	rotation         float32
	scaleX, scaleY   float32
}

type Sprite interface {
	Bounds() rect
	Draw(t Transform, canvas *image.RGBA)
}

type Camera interface {
	Transform() Transform
}

type World interface {
	Draw(camera Camera, rgba *image.RGBA)
}

type SpriteWorld struct {
	sprites []Sprite
}

func (world *SpriteWorld) Draw(camera Camera, rgba *image.RGBA) {
	width := float32(rgba.Bounds().Max.X - rgba.Bounds().Min.X)
	height := float32(rgba.Bounds().Max.Y - rgba.Bounds().Min.Y)
	screenBounds := rect{
		camera.Transform().offsetX - width/2,
		camera.Transform().offsetY - height/2,
		camera.Transform().offsetX + width/2,
		camera.Transform().offsetY + height/2}

	for _, sprite := range world.sprites {
		if screenBounds.overlap(sprite.Bounds()) {
			sprite.Draw(camera.Transform(), rgba)
		}
	}
}

func (s SimpleSprite) Bounds() rect {
	return rect{
		s.x,
		s.y,
		float32(s.img.Bounds().Max.X - s.img.Bounds().Min.X),
		float32(s.img.Bounds().Max.Y - s.img.Bounds().Max.Y)}
}

func (s SimpleSprite) Draw(t Transform, canvas *image.RGBA) {
	spriteBounds := s.img.Bounds()
	fmt.Printf("Sprite bounds: %v\n", spriteBounds)
	fmt.Printf("s.x: %v\n", s.x)

	canvasBounds := canvas.Bounds()

	for x := canvasBounds.Min.X; x < canvasBounds.Max.X; x++ {
		for y := canvasBounds.Min.Y; y < canvasBounds.Max.Y; y++ {
			base := y*canvas.Stride + x*4

			// canvas.Pix[base] = 100
			sx := int((float32(x) + t.offsetX) / t.scaleX)
			sy := int((float32(y) + t.offsetY) / t.scaleY)
			if sx > spriteBounds.Min.X && sx < spriteBounds.Max.X && sy > spriteBounds.Min.Y && sy < spriteBounds.Max.Y {
				srcBase := sy*s.img.Stride + sx*4
				canvas.Pix[base] = s.img.Pix[srcBase]
				canvas.Pix[base+1] = s.img.Pix[srcBase+1]
				canvas.Pix[base+2] = s.img.Pix[srcBase+2]
			}
		}
	}
}

func NewSprite(srcImg *image.Image) Sprite {
	srcBounds := (*srcImg).Bounds()
	rgbaImg := image.NewRGBA(image.Rect(0, 0, srcBounds.Size().X, srcBounds.Size().Y))

	rgbaBounds := rgbaImg.Bounds()
	for x := rgbaBounds.Min.X; x < rgbaBounds.Max.X; x++ {
		srcX := x + srcBounds.Min.X
		for y := rgbaBounds.Min.Y; y < rgbaBounds.Max.Y; y++ {
			srcY := y + srcBounds.Min.Y
			srcColor := (*srcImg).At(srcX, srcY)
			rgbaImg.Set(x, y, srcColor)
		}
	}

	return SimpleSprite{x: 0, y: 0, rotation: 0, scaleX: 1, scaleY: 1, img: rgbaImg}
}
