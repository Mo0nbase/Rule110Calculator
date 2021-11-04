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
)

var win *pixelgl.Window

type cell struct {
	alive  bool
	square pixel.Rect
	obj    *imdraw.IMDraw
}

var white, _ = colorful.MakeColor(colornames.Antiquewhite)
var black, _ = colorful.MakeColor(colornames.Black)

var cells [][]cell

func run() {
	initWindow()
	loadCells(true)
	ticker := pixelutils.NewTicker(120)

	simmer := false
	cl, clpc := 0, 0 //color layer, color loop count
	for !win.Closed() {
		win.Clear(colornames.Antiquewhite)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			simmer = true
		}

		if simmer {
			if cl < len(cells) {
				updateColors(uint64(cl), uint64(clpc))
				//TODO this here is whats most likely broken check that you area actually updating the colors correctly
				clpc++
				if clpc == 120 {
					clpc = 0
					cl++
				}
			}

			for i := range cells {
				for j := range cells[i] {
					cells[i][j].obj.Draw(win)
				}
			}
		}

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
				cells[i][j].square = pixel.R(win.Bounds().Size().X-float64((len(data[i])-j)*10), win.Bounds().Size().Y-20-(float64(i*10)),
					win.Bounds().Size().X-float64(((len(data[i])-j)*10)+10), win.Bounds().Size().Y-20-(float64(i*10))+10)
				cells[i][j].obj = imdraw.New(nil)
				if data[i][j] == 1 {
					cells[i][j].obj.Color = black
					cells[i][j].alive = true
					cells[i][j].obj.Push(cells[i][j].square.Min, cells[i][j].square.Max)
					cells[i][j].obj.Rectangle(0)
				} else {
					cells[i][j].obj.Color = white
					cells[i][j].obj.Push(cells[i][j].square.Min, cells[i][j].square.Max)
					cells[i][j].obj.Rectangle(0)
				}
			}
		}
	} else {

	}
}

func updateColors(layer uint64, iteration uint64) {
	//NOTE iterations start at 0
	for i := range cells[layer] {
		if cells[layer][i].alive {
			cells[layer][i].obj.Color = white.BlendRgb(black, float64(iteration/199))
		}
	}
	// (250,235,215)
	// for all cells that are alive
	// the degree should reduce the colors to zero in a quarter of the framerate iterations
	// it should not begin the next row until the previous has been completed

}

func resizeCells(degree float64) {

}
