// +build windows

package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/lxn/walk"
	decl "github.com/lxn/walk/declarative"
)

//go:generate rsrc -manifest manifest/goiv.exe.manifest -o goiv_windows.syso

func display(image image.Image, width, height int) {
	mw := new(Window)
	mw.goImg = image

	keyEvent := func(key walk.Key) {
		switch key {
		case walk.KeyQ, walk.KeyEscape:
			mw.Close()
		case walk.KeyF11, walk.KeyF:
			mw.SetFullscreen(!mw.Fullscreen())
		}
	}

	if err := (decl.MainWindow{
		AssignTo:  &mw.MainWindow,
		OnKeyDown: keyEvent,
		MinSize:   decl.Size{320, 240},
		Size:      decl.Size{width, height},
		Layout:    decl.VBox{MarginsZero: true, SpacingZero: true},
		Children: []decl.Widget{
			decl.ImageView{
				AssignTo:   &mw.imageView,
				Background: decl.SolidColorBrush{Color: walk.RGB(0, 0, 0)},
				OnKeyDown:  keyEvent,
			},
		},
	}.Create()); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	mw.drawImageError()

	mw.Run()
}

type Window struct {
	*walk.MainWindow

	image     walk.Image
	imageView *walk.ImageView
	goImg     image.Image
}

func (mw *Window) drawImageError() {
	if err := mw.drawImage(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

func (mw *Window) drawImage() error {
	var err error

	if mw.image != nil {
		mw.image.Dispose()
		mw.image = nil
	}

	mw.image, err = walk.NewBitmapFromImage(mw.goImg)
	// mw.image, err = walk.NewImageFromFile(mw.images[mw.idx])
	if err != nil {
		return err
	}

	var succeeded bool
	defer func() {
		if !succeeded {
			if mw.image != nil {
				mw.image.Dispose()
			}
		}
	}()

	if mw.imageView == nil {
		mw.imageView, err = walk.NewImageView(mw)
		if err != nil {
			return err
		}
	}

	mw.imageView.SetMode(walk.ImageViewModeShrink)
	if err = mw.imageView.SetImage(mw.image); err != nil {
		return err
	}

	mw.SetTitle("sprites")

	succeeded = true

	return nil
}

func main() {
	img := loadImage("./alien.png")

	s := NewSprite(img)
	ss := s.(SimpleSprite)
	ss.x = 100

	canvas := image.NewRGBA(image.Rect(0, 0, 640, 480))
	camera := Transform{0, 0, 0, 0.5, 0.5}
	s.Draw(camera, canvas)

	display(canvas, 1024, 768)
}

func loadImage(path string) *image.Image {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(path + " not found!")
		os.Exit(1)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("could not decode image!")
		os.Exit(1)
	}

	return &img
}
