package main

import (
	"github.com/dusk125/pixelutils"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"image"
	"image/png"
	_ "image/png"
	"log"
	_ "math"
	"os"
)

var win *pixelgl.Window
var imd *imdraw.IMDraw
var ticker *pixelutils.Ticker
var framerate = int64(120)

var white, _ = colorful.MakeColor(colornames.Antiquewhite)
var black, _ = colorful.MakeColor(colornames.Black)

var cells [][]uint64

func run() {
	initWindow()
	createImage(true, 10)

	for !win.Closed() {
		win.Clear(colornames.Antiquewhite)

		// TODO add zoom and pan ability for the user
		imd.Draw(win)
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
	imd = imdraw.New(nil)
	ticker = pixelutils.NewTicker(framerate)
}

func createImage(history bool, pSize int) {
	data := decompress(history)
	cells = data

	width := len(data[0]) * pSize
	height := len(data) * pSize

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	if history {
		for i := range data {
			for j := range data[i] {
				for k := j * pSize; k <= (j*pSize)+pSize; k++ {
					for l := i * pSize; l <= (i*pSize)+pSize; l++ {
						if data[i][j] == 1 {
							img.Set(k, l, black)
						} else {
							img.Set(k, l, white)
						}
					}
				}
			}
		}
	} else {

	}
	f, _ := os.Create("image.png")
	_ = png.Encode(f, img)
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

func translateCells(vector pixel.Vec) {

}

func resizeCells(degree float64) {
	// cells by constant amount

	// you cant scale with floats due to errors
	// 1 = scale 10%
	// 11 = scale 110%
}

func cellCount() int {
	counter := 0
	for i := range cells {
		counter += len(cells[i])
	}
	return counter
}
