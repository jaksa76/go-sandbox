package main

import (
	"fmt"
	"image"
	"image/color"
	"syscall/js"
	"time"
)

var (
	window                     js.Value
	canvas, jsCtx, jsImageData js.Value
	beep                       js.Value
	windowSize                 struct{ w, h float64 }
	buffer                     *image.RGBA
	lastFPS                    time.Time
	jsArray                    js.TypedArray
	frameNo                    int
	periodStart                int
)

func main() {
	runGameForever := make(chan bool)

	setup()
	fmt.Println("setup completed")

	var renderer js.Func

	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updateGame()
		draw()
		window.Call("requestAnimationFrame", renderer)
		return nil
	})
	window.Call("requestAnimationFrame", renderer)

	var mouseEventHandler js.Func = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		updateMouse(args[0])
		return nil
	})
	window.Call("addEventListener", "pointerdown", mouseEventHandler)

	<-runGameForever
}

func setup() {
	window = js.Global()
	doc := window.Get("document")
	body := doc.Get("body")

	// windowSize.h = window.Get("innerHeight").Float()
	// windowSize.w = window.Get("innerWidth").Float()
	windowSize.w = 1920
	windowSize.h = 1080

	canvas = doc.Call("createElement", "canvas")
	canvas.Set("height", windowSize.h)
	canvas.Set("width", windowSize.w)
	body.Call("appendChild", canvas)

	buffer = image.NewRGBA(image.Rect(0, 0, int(windowSize.w), int(windowSize.h)))

	jsCtx = canvas.Call("getContext", "2d")
	jsImageData = jsCtx.Call("createImageData", windowSize.w, windowSize.h)
	jsArray = js.TypedArrayOf(buffer.Pix)

	// http://www.iandevlin.com/blog/2012/09/html5/html5-media-and-data-uri/
	beep = window.Get("Audio").New("data:audio/mp3;base64,SUQzBAAAAAAAI1RTU0UAAAAPAAADTGF2ZjU2LjI1LjEwMQAAAAAAAAAAAAAA/+NAwAAAAAAAAAAAAFhpbmcAAAAPAAAAAwAAA3YAlpaWlpaWlpaWlpaWlpaWlpaWlpaWlpaWlpaWlpaWlpaW8PDw8PDw8PDw8PDw8PDw8PDw8PDw8PDw8PDw8PDw8PDw////////////////////////////////////////////AAAAAExhdmYAAAAAAAAAAAAAAAAAAAAAACQAAAAAAAAAAAN2UrY2LgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP/jYMQAEvgiwl9DAAAAO1ALSi19XgYG7wIAAAJOD5R0HygIAmD5+sEHLB94gBAEP8vKAgGP/BwMf+D4Pgh/DAPg+D5//y4f///8QBhMQBgEAfB8HwfAgIAgAHAGCFAj1fYUCZyIbThYFExkefOCo8Y7JxiQ0mGVaHKwwGCtGCUkY9OCugoFQwDKqmHQiUCxRAKOh4MjJFAnTkq6QqFGavRpYUCmMxpZnGXJa0xiJcTGZb1gJjwOJDJgoUJG5QQuDAsypiumkp5TUjrOobR2liwoGBf/X1nChmipnKVtSmMNQDGitG1fT/JhR+gYdCvy36lTrxCVV8Paaz1otLndT2fZuOMp3VpatmVR3LePP/8bSQpmhQZECqWsFeJxoepX9dbfHS13/////aysppUblm//8t7p2Ez7xKD/42DE4E5z9pr/nNkRw6bhdiCAZVVSktxunhxhH//4xF+bn4//6//3jEvylMM2K9XmWSn3ah1L2MqVIjmNlJtpQux1n3ajA0ZnFSu5EpX////uGatn///////1r/pYabq0mKT//TRyTEFNRTMuOTkuNaqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq/+MQxNIAAANIAcAAAKqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqg==")

	lastFPS = time.Now()
	frameNo = 0
	periodStart = 0
}

func updateMouse(event js.Value) {
	mouseX := event.Get("clientX").Float()
	mouseY := event.Get("clientY").Float()
	fmt.Println("click at ", mouseX, mouseY)
}

func playSound() {
	beep.Call("play")
	window.Get("navigator").Call("vibrate", 300)
}

func updateGame() {
	now := time.Now()
	t := now.Nanosecond()/1000000 + now.Second()*1000
	t = t / 20
	width := int(windowSize.w)
	height := int(windowSize.h)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			base := (y*width + x) * 4
			buffer.Pix[base] = uint8(y) + uint8(t)
			buffer.Pix[base+1] = uint8(x)
			buffer.Pix[base+2] = uint8(t)
			buffer.Pix[base+3] = 255
		}
	}
}

func col(r, g, b uint8) color.RGBA {
	return color.RGBA{R: clamp(r), G: clamp(g), B: clamp(b), A: 255}
}

func clamp(n uint8) uint8 {
	return n
}

func draw() {
	printStats()
	// jsArray = js.TypedArrayOf(buffer.Pix)
	jsImageData.Get("data").Call("set", jsArray)
	// jsArray.Release()
	jsCtx.Call("putImageData", jsImageData, 0, 0)
}

func printStats() {
	frameNo++
	now := int(time.Now().UnixNano() / int64(time.Millisecond))
	delta := now - periodStart
	if delta > 1000 {
		fps := 1000.0 * float64(frameNo) / float64(delta)
		log(fmt.Sprintf("fps: %f", fps))
		periodStart = now
		frameNo = 0
	}
}

func log(s string) {
	window.Get("console").Call("log", s)
}
