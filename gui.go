package main

import (
	"github.com/dusk125/pixelutils"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
	"image/color"
	_ "image/png"
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

//var batches map[*pixel.Batch]pixel.Batch

//func main() {
//	pixelgl.Run(run)
//}

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

	square := imdraw.New(nil)
	square.Color = colornames.Black
	square.Push(pixel.V(960, 540), pixel.V(965, 545))
	square.Rectangle(0)

	win.SetComposeMethod(pixel.ComposeIn)

	//batch := pixel.NewBatch(&pixel.TrianglesData{}, square)

	lpc := 0
	for !win.Closed() {
		//win.Clear(colornames.Antiquewhite)
		win.Clear(colornames.White)

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			square.Clear()
			square.Color = color.RGBA{
				R: 255,
				G: 0,
				B: 135,
				A: 255,
			}
			square.Push(pixel.V(960, 540), pixel.V(1200, 800))
			square.Rectangle(0)
		}

		square.Draw(win)

		lpc++
		win.Update()
		ticker.Wait()
	}

	// Method 1:
	// Method 2:
	// Method 3:
	// Method 4:
}
