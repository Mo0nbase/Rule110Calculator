package main

import (
	"fmt"
	"github.com/dusk125/pixelutils"
	"github.com/inkyblackness/imgui-go"
	_ "image/png"
	"strconv"
	"sync"

	"encoding/csv"
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
	createImages(true, 2)

	pic, _ := loadPicture("export/image.png")
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
			camPos = camPos.Sub(win.MousePosition().Sub(win.MousePreviousPosition())).ScaledXY(pixel.V(1, 1)) // /camZoom
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
		// TODO import a batch of all the chunks to use for the translation matrix

		imgui.BeginV("Rule 110 Simulator", nil, 0) //0b00000110

		imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorWindowBg, pixelui.ColorA(119, 134, 127, 245))
		imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorHeader, pixelui.ColorA(130, 79, 80, 220))
		imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorHeaderHovered, pixelui.ColorA(130, 94, 95, 215))
		imgui.ProgressBarV(0.0, imgui.Vec2{X: 200, Y: 20}, "Total")

		imgui.End()

		win.Clear(colornames.Antiquewhite)
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Max.ScaledXY(pixel.V(0.5, 0.5)).Sub(sprite.Frame().Center())))
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

func createImages(history bool, pSize int) {
	width := len(sim[0]) * pSize * 64
	height := len(sim) * pSize

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	var wg sync.WaitGroup
	wg.Add(len(sim))
	if history {
		csvFile, _ := os.Create("multithreading.csv")
		csvFile1, _ := os.Create("single-thread.csv")

		writer := csv.NewWriter(csvFile)
		writer1 := csv.NewWriter(csvFile1)

		for i := range sim {
			renderRowNormal(img, i, pSize) // TODO make multithreaded again
		}
		for i := 0; i < len(img.Pix)/4; i++ {
			_ = writer1.Write([]string{strconv.Itoa(int(img.Pix[i*4])), strconv.Itoa(int(img.Pix[(i*4)+1])), strconv.Itoa(int(img.Pix[(i*4)+2])), strconv.Itoa(int(img.Pix[(i*4)+3]))})
		}
		img = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
		for i := range sim {
			go renderRow(&wg, img, i, pSize) // TODO make multithreaded again
		}
		wg.Wait()
		for i := 0; i < len(img.Pix)/4; i++ {
			_ = writer.Write([]string{strconv.Itoa(int(img.Pix[i*4])), strconv.Itoa(int(img.Pix[(i*4)+1])), strconv.Itoa(int(img.Pix[(i*4)+2])), strconv.Itoa(int(img.Pix[(i*4)+3]))})
		}
		writer.Flush()
		writer1.Flush()
		csvFile.Close()
		csvFile1.Close()
	} else {
		// TODO add non history version
	}
	//wg.Wait()

	fmt.Println("Finished passing to threads")
	_ = os.MkdirAll("export/chunks", os.ModePerm)
	f, _ := os.Create("export/image.png")
	_ = png.Encode(f, img)

	grid := gridSplit(img.Rect)
	if grid == nil {
		fmt.Println("Image is already under maximum texture size")
	} else {
		for i := range grid {
			for j := range grid[i] {
				name := strconv.Itoa(i) + "_" + strconv.Itoa(j)
				f, _ = os.Create("export/chunks/" + name + ".png")
				_ = png.Encode(f, img.SubImage(grid[i][j]))
			}
		}
	}
	fmt.Println("Finished Generation")
}

func renderRowNormal(img *image.RGBA, i int, pSize int) {
	// i=the current row of the simulation
	// j=loop through each row for the 64 bit integers
	// k=the current bit on the current row looping through horizontally
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

func importSpriteMatrix() {
	entries, _ := os.ReadDir("export/chunks")
	for _, entry := range entries {
		fmt.Println(entry.Type())
	}
	//fmt.Println("Could not open chunk file! (It may not exist or you may not have exported anything)")
}

func gridSplit(r image.Rectangle) [][]image.Rectangle {
	MAX := 8192
	if r.Dx()*r.Dy() < MAX*MAX {
		return nil
	} // check to make sure original is not less than max already
	var temp [][]image.Rectangle
	temp = append(temp, nil)
	temp[0] = append(temp[0], r)
	return getSubRectangles(temp, true)
}

func getSubRectangles(rectangles [][]image.Rectangle, horizontal bool) [][]image.Rectangle {
	fmt.Println("called")
	MAX := 8192
	var temp [][]image.Rectangle
	if horizontal {
		for i := range rectangles {
			temp = append(temp, nil) // add a row 1x previous
			for j := range rectangles[i] {
				temp[i] = append(temp[i], splitRect(rectangles[i][j], true)[0]...) // populate columns 2x previous
			}
		}
	} else {
		for i := range rectangles {
			temp = append(temp, nil) // add rows 2x previous
			temp = append(temp, nil)
			for j := range rectangles[i] {
				temp[i*2] = append(temp[i*2], splitRect(rectangles[i][j], false)[0]...) // populate columns 1x previous
				temp[(i*2)+1] = append(temp[i*2+1], splitRect(rectangles[i][j], false)[1]...)
			}
		}
	}
	// check if its below the max
	// otherwise call the same function again with the inverse of the horizontal
	if temp[0][0].Dx()*temp[0][0].Dy() < MAX*MAX {
		return temp
	} else {
		return getSubRectangles(temp, !horizontal)
	}
}

func splitRect(r image.Rectangle, horizontal bool) [][]image.Rectangle {
	// if horizontal is true then split horizontally otherwise split vertically
	var temp [][]image.Rectangle
	if horizontal {
		temp = append(temp, nil)
		diff := int(math.Floor(float64(r.Max.X-r.Min.X) / 2))
		temp[0] = append(temp[0], image.Rect(r.Min.X, r.Min.Y, r.Min.X+diff, r.Max.Y))
		temp[0] = append(temp[0], image.Rect(r.Min.X+diff, r.Min.Y, r.Max.X, r.Max.Y))
		return temp
	} else {
		temp = append(temp, nil)
		temp = append(temp, nil)
		diff := int(math.Floor(float64(r.Max.Y-r.Min.Y) / 2))
		temp[0] = append(temp[0], image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+diff))
		temp[1] = append(temp[1], image.Rect(r.Min.X, r.Min.Y+diff, r.Max.X, r.Max.Y))
		return temp
	}
}
