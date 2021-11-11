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

var cells [][]cell

func run() {
	initWindow()
	loadCells(true)

	simmer := false
	//colorLayer := 0
	//multiplier := float64(framerate) * 1.0 // chagne to 0.8
	for !win.Closed() {
		win.Clear(colornames.Antiquewhite)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			simmer = !simmer
		}

		if simmer {
			if cells[len(cells)-1][0].blend < 10 {
				for i := 0; i < len(cells); i++ {
					updateColors(uint64(i))
					if cells[i][0].blend-1 == 0 {
						break
					}
					// begin looping through the rows for cells
					// if you run into a blend == 0 then update the row and break into the next loop
					// loop through until the blend is 0 then adjust the color then break out to start the next iteration
				}
			}
			// TODO fix so that the multiplier if over the frame rate the simulation starts updating multiple rows at once

			// TODO add inital bounds checks lock out input untill the simulation has exponentially rendered all of the availiable cells within a starting area
			// TODO add zoom and pan ability for the user
		}

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
					cells[i][j].square = pixel.R(win.Bounds().Size().X-float64((len(data[i])-j)*3), win.Bounds().Size().Y-20-(float64(i*3)),
						win.Bounds().Size().X-float64(((len(data[i])-j)*3)+3), win.Bounds().Size().Y-20-(float64(i*3))+3)
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
		if cells[layer][i].blend < 10 {
			imd.Color = white.BlendRgb(black, cells[layer][i].blend/10)
			imd.Push(cells[layer][i].square.Min, cells[layer][i].square.Max)
			imd.Rectangle(0)
			cells[layer][i].blend++
		}
	}
}

func resizeCells(degree float64) {
	// cells by constant amount

	// you cant scale with floats due to errors
	// 1 = scale 10%
	// 11 = scale 110%
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
