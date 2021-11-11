package main

import (
	"github.com/dusk125/pixelutils"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	_ "image/png"
	"log"
	_ "math"
)

type cell struct {
	alive  bool
	blend  float64
	square pixel.Rect
}

var win *pixelgl.Window
var imd *imdraw.IMDraw
var ticker *pixelutils.Ticker
var framerate = int64(120)

var white, _ = colorful.MakeColor(colornames.Antiquewhite)
var black, _ = colorful.MakeColor(colornames.Black)
var blendMax = float64(8)

var cells [][]cell

func run() {
	initWindow()
	loadCells(true)

	drawInitial := false
	offset := int(-blendMax)
	for !win.Closed() {
		win.Clear(colornames.Antiquewhite)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			drawInitial = true
		}

		if drawInitial {
			if cells[len(cells)-1][0].blend < blendMax {
				for i := offset; i <= offset+int(blendMax) && i < len(cells); i++ {
					if i >= 0 {
						updateColors(uint64(i))
					}
				}
				offset++
			} else {
				drawInitial = false
			}
		}
		// might be easier to store this as a constant
		if drawInitial && int((win.Bounds().Size().Y-20.0)/6.0) < len(cells) && cells[int((win.Bounds().Size().Y-20.0)/6.0)][0].blend == blendMax {
			drawInitial = false
			for i := range cells {
				for j := range cells[i] {
					imd.Color = black
					imd.Push(cells[i][j].square.Min, cells[i][j].square.Max)
					imd.Rectangle(0)
				}
			}
		}

		// TODO figure out how many evolutions it takes to reach the end of the window
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

func loadCells(history bool) {
	data := decompress(history)

	cells = make([][]cell, len(data))
	for i := range cells {
		cells[i] = make([]cell, len(data[0]))
	}
	if history {
		for i := range data {
			for j := range data[i] {
				if data[i][j] == 1 {
					cells[i][j].square = pixel.R(win.Bounds().Size().X-float64((len(data[i])-j)*6), win.Bounds().Size().Y-15-(float64(i*6)),
						win.Bounds().Size().X-float64(((len(data[i])-j)*6)+6), win.Bounds().Size().Y-15-(float64(i*6))+6)
					imd.Color = white
					cells[i][j].alive = true
					cells[i][j].blend = 0
					imd.Push(cells[i][j].square.Min, cells[i][j].square.Max)
					imd.Rectangle(0)
				}
			}
		}
	} else {

	}
	for i := range cells {
		cells[i] = removeDead(cells[i])
	}
}

func updateColors(layer uint64) {
	for i := range cells[layer] {
		if cells[layer][i].blend < blendMax {
			imd.Color = white.BlendRgb(black, cells[layer][i].blend/blendMax)
			imd.Push(cells[layer][i].square.Min, cells[layer][i].square.Max)
			imd.Rectangle(0)
			cells[layer][i].blend++
		}
	}
}

func translateCells(vector pixel.Vec) {
	for i := range cells {
		for j := range cells[i] {
			cells[i][j].square = cells[i][j].square.Moved(vector)
			imd.Push(cells[i][j].square.Min, cells[i][j].square.Max)
			imd.Rectangle(0)
		}
	}
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

func removeDead(objs []cell) []cell {
	count := 0
	for i := range objs {
		if objs[i].alive {
			count++
		}
	}
	temp := make([]cell, count)
	count = 0
	for i := range objs {
		if objs[i].alive {
			temp[count] = objs[i]
			count++
		}
	}
	return temp
}
