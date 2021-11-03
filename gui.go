package main

import (
	"github.com/dusk125/pixelutils"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image/color"
	_ "image/png"
	"log"
)

var cells [][]pixel.Rect

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Cellular Automata Calculator",
		Bounds:    pixel.R(0, 0, 1920, 1080),
		VSync:     true,
		Maximized: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ticker := pixelutils.NewTicker(120)

	cell := pixel.R(960, 540, 990, 570)
	square := imdraw.New(nil)
	square.Color = colornames.Black
	square.Push(cell.Min, cell.Max)
	square.Rectangle(0)

	win.SetComposeMethod(pixel.ComposeIn)

	lpc := 0
	grow := false
	shrink := false
	for !win.Closed() {
		win.Clear(colornames.Antiquewhite)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			square.Clear()
			square.Color = color.RGBA{
				R: 255,
				G: 0,
				B: 135,
				A: 255,
			}
			grow = true
		}

		// if pressed again its growing and shrinking at the same time in the same iteration
		if grow == true {
			square.Clear()
			//cell = cell.Resized(cell.Center(), cell.Size().Add(pixel.V(cell.Size().X+0.005, cell.Size().Y+0.005)))
			cell = cell.Resized(cell.Center(), cell.Size().Scaled(1.01))
			square.Push(cell.Min, cell.Max)
			square.Rectangle(0)
			if cell.Size().X >= 500 && cell.Size().Y >= 500 {
				shrink = true
				grow = false
			}
		}
		if shrink == true {
			square.Clear()
			cell = cell.Resized(cell.Center(), cell.Size().Scaled(0.99))
			square.Push(cell.Min, cell.Max)
			square.Rectangle(0)
			if cell.Size().X <= 30 && cell.Size().Y <= 30 {
				grow = true
				shrink = false
			}
		}

		square.Draw(win)

		lpc++
		win.Update()
		ticker.Wait()
	}
}

func loadCells() {

}

func resizeCells(degree float64) {

}
