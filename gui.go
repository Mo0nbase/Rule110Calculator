package main

import (
	"fmt"
	"github.com/dusk125/pixelutils"
	"github.com/inkyblackness/imgui-go"
	_ "image/png"
	"sync"

	"github.com/dusk125/pixelui"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"image"
	"image/png"
	"log"
	"math"
	_ "math"
	"os"
	"time"
)

var win *pixelgl.Window
var ticker *pixelutils.Ticker
var framerateTarget = int64(120)

var white, _ = colorful.MakeColor(colornames.Antiquewhite)
var black, _ = colorful.MakeColor(colornames.Black)

func run() {
	initWindow()
	createImage(true, 4)

	pic, _ := loadPicture("image.png")
	sprite := pixel.NewSprite(pic, pic.Bounds())

	var (
		camPos       = pixel.ZV
		camSpeed     = 800.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	ui := pixelui.NewUI(win, 0)
	defer ui.Destroy()

	ui.AddTTFFont("resources/03b04.ttf", 16)
	ui.AddTTFFont("resources/Roboto-Medium.ttf", 16)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)
		//_, framerate := ticker.Tick()

		if win.JustReleased(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		if !imgui.CurrentIO().WantCaptureMouse() && win.Pressed(pixelgl.MouseButtonLeft) {
			camPos = camPos.Sub(win.MousePosition().Sub(win.MousePreviousPosition()))
		}
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt
		}
		if !imgui.CurrentIO().WantCaptureMouse() {
			camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		}

		ui.NewFrame()
		imgui.ShowDemoWindow(nil)

		//TODO imgui_demo.cpp

		imgui.BeginV("Rule 110 Simulator", nil, 0) //0b00000110

		imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorWindowBg, pixelui.ColorA(119, 134, 127, 245))
		imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorHeader, pixelui.ColorA(130, 79, 80, 220))
		imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorHeaderHovered, pixelui.ColorA(130, 94, 95, 215))
		imgui.ProgressBarV(0.0, imgui.Vec2{X: 200, Y: 20}, "Total")

		imgui.End()

		win.Clear(colornames.Antiquewhite)

		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Max.Sub(pic.Bounds().Center())))
		ui.Draw(win)

		win.Update()
		ticker.Wait()
	}
}

func initWindow() {
	if win != nil {
		win = nil
	}

	cfg := pixelgl.WindowConfig{
		Title:     "Cellular Automata Calculator",
		Bounds:    pixel.R(0, 0, 1920, 1080),
		VSync:     true,
		Maximized: true,
	}

	err := error(nil)
	win, err = pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	win.SetComposeMethod(pixel.ComposeIn)
	ticker = pixelutils.NewTicker(framerateTarget)
}

func createImage(history bool, pSize int) {
	width := len(sim[0]) * pSize * 64
	height := len(sim) * pSize

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	var wg sync.WaitGroup
	wg.Add(len(sim))
	if history {
		for i := range sim {
			go renderRow(&wg, img, i, pSize)
		}
	} else {

	}
	wg.Wait()

	fmt.Println("Finished passing to threads")
	f, _ := os.Create("image.png")
	_ = png.Encode(f, img)
	fmt.Println("Finished Generation")
}

func renderRow(wg *sync.WaitGroup, img *image.RGBA, i int, pSize int) {
	// i=the current row of the simulation
	// j=loop through each row for the 64 bit integers
	// k=the current bit on the current row looping through horizontally
	defer wg.Done()
	var lpc = 0
	for j := 0; j < 64; j++ {
		for k := range sim[i] {
			for l := lpc * pSize; l <= (lpc*pSize)+pSize; l++ {
				for m := i * pSize; m <= (i*pSize)+pSize; m++ {
					if getBit(sim[i][k], j) == 1 {
						img.Set(l, m, black)
					} else {
						img.Set(l, m, white)
					}
				}
			}
			lpc++
		}
	}
	// function loops the proper amount
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func getTranslationMatrix(vector pixel.Vec) pixel.Matrix {
	return pixel.IM.Moved(vector)
}

func maxRectDim(largeX, int, largeY int, smallX int, smallY int) {

}

// TODO fix this function
//func cellCount() int {
//	counter := 0
//	for i := range cells {
//		counter += len(cells[i])
//	}
//	return counter
//}
