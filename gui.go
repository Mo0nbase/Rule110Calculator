package main

import (
	"encoding/gob"
	"fmt"
	"github.com/dusk125/pixelutils"
	"github.com/inkyblackness/imgui-go"
	_ "image/png"
	"runtime"
	"strconv"
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

var onlyOnce sync.Once

var win *pixelgl.Window
var ticker *pixelutils.Ticker
var framerateTarget = int64(120)

var white, _ = colorful.MakeColor(colornames.Antiquewhite)
var black, _ = colorful.MakeColor(colornames.Black)

var generationProgress = float32(0.0)
var threadProgress = float32(0.0)
var routines = 0.0

func run() {
	initWindow()
	//pic, _ := loadPicture("export/image.png")

	var (
		camPos       = pixel.ZV
		camSpeed     = 800.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	ui := pixelui.NewUI(win, 0)
	defer ui.Destroy()

	ui.AddTTFFont("resources/03b04.ttf", 16)

	var sprites [][]pixel.Sprite

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		if win.JustReleased(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}
		if !imgui.CurrentIO().WantCaptureMouse() && win.Pressed(pixelgl.MouseButtonLeft) {
			camPos = camPos.Sub(win.MousePosition().Sub(win.MousePreviousPosition()).Scaled(1 / camZoom))
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
		if win.Pressed(pixelgl.KeyP) {
			//simulate(true, 57800, simpleCTS())
			simulate(true, 3500, randStart(3500))
		}
		if win.Pressed(pixelgl.KeyO) {
			onlyOnce.Do(func() { go createImages(true, 2) })
		}
		if win.Pressed(pixelgl.KeyI) {
			gob.Register([][]string{})
			sprites = importSpriteMatrix()
		}

		ui.NewFrame()
		imgui.ShowDemoWindow(nil)

		//TODO imgui_demo.cpp
		// TODO import a batch of all the chunks to use for the translation matrix

		drawUI()

		win.Clear(colornames.Antiquewhite)
		//sprite.Draw(win, pixel.IM.Moved(win.Bounds().Max.ScaledXY(pixel.V(0.5, 0.5)).Sub(sprite.Frame().Center())))

		drawSpriteMatrix(sprites)
		// TODO add checking to not call the splitter function if under maximum size

		ui.Draw(win)

		win.Update()
		ticker.Wait()
	}
}

func drawUI() {
	imgui.BeginV("Rule 110 Simulator", nil, 0b00000000) //0b00000110 //noResize, noMove

	imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorWindowBg, pixelui.ColorA(119, 134, 127, 245))
	imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorHeader, pixelui.ColorA(130, 79, 80, 220))
	imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorHeaderHovered, pixelui.ColorA(130, 94, 95, 215))
	imgui.Style.SetColor(imgui.CurrentStyle(), imgui.StyleColorPlotHistogram, pixelui.ColorA(130, 79, 80, 220))

	_, framerate := ticker.Tick()
	imgui.Text(fmt.Sprintf("%.2f", framerate))
	imgui.ProgressBarV(generationProgress/float32(len(sim)), imgui.Vec2{X: 375, Y: 22}, "Generation")
	imgui.ProgressBarV(threadProgress/float32(len(sim)), imgui.Vec2{X: 375, Y: 22}, "Threads")
	imgui.ProgressBarV(float32(routines/float64(runtime.NumGoroutine())), imgui.Vec2{X: 375, Y: 22}, "Progress")

	imgui.End()
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
	fmt.Println("This is bad if it happens more than once!")
	width := len(sim[0]) * pSize * 64
	height := len(sim) * pSize

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	threadProgress = 0
	var wg sync.WaitGroup
	routines = float64(runtime.NumGoroutine() + len(sim))
	wg.Add(len(sim))
	if history {
		img = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
		for i := range sim {
			go renderRow(&wg, img, i, pSize)
			threadProgress++
		}
	} else {
		// TODO add non history version
	}
	wg.Wait()

	fmt.Println("Finished passing to threads")

	//_, err := os.OpenFile("export/chunks/backup", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	//if err != nil {
	//	_ = os.RemoveAll("export/chunks/backup")
	//}
	//dir, err := os.ReadDir("export/chunks")
	//if err != nil {
	//	return
	//}
	//for entry := range dir {
	//	os.WriteFile(os.ReadFile(string(entry)))
	//}

	_ = os.MkdirAll("export/chunks", os.ModePerm)
	f, _ := os.Create("export/image.png")
	_ = png.Encode(f, img)

	grid := gridSplit(img.Rect)
	exportMatrix := make([][]string, len(grid))
	for i := range exportMatrix {
		exportMatrix[i] = make([]string, len(grid[0]))
	}
	if grid == nil {
		fmt.Println("Image is already under maximum texture size")
	} else {
		for i := range grid {
			for j := range grid[i] {
				name := strconv.Itoa(i) + "_" + strconv.Itoa(j)
				exportMatrix[i][j] = name + ".png"
				f, _ = os.Create("export/chunks/" + name + ".png")
				_ = png.Encode(f, img.SubImage(grid[i][j]))
			}
		}
	}
	gob.Register([][]string{})
	writeToFile("export/chunks/matrix.bin", exportMatrix)
	fmt.Println("Finished Generation")
}

func drawSpriteMatrix(sprites [][]pixel.Sprite) {
	if sprites != nil {
		delta := pixel.IM.Moved(win.Bounds().Max.ScaledXY(pixel.V(0.5, 0.5)).Sub(sprites[0][len(sprites[0])-1].Frame().Center()))
		sprites[0][len(sprites[0])-1].Draw(win, delta)
		DX, DY := -sprites[0][len(sprites[0])-1].Frame().Max.X, -sprites[0][len(sprites[0])-1].Frame().Max.Y

		for i := len(sprites[0]) - 2; i >= 0; i-- {
			sprites[0][i].Draw(win, delta.Moved(pixel.Vec{X: DX, Y: 0}))
			DX -= sprites[0][i].Frame().Max.X
		}

		DX = -sprites[0][len(sprites[0])-1].Frame().Max.X

		for i := 1; i < len(sprites); i++ {
			sprites[i][len(sprites[i])-1].Draw(win, delta.Moved(pixel.Vec{X: 0, Y: DY}))
			for j := len(sprites[i]) - 2; j >= 0; j-- {
				sprites[i][j].Draw(win, delta.Moved(pixel.Vec{X: DX, Y: DY}))
				DX -= sprites[i][j].Frame().Max.X
			}

			DY -= sprites[i][len(sprites)-1].Frame().Max.Y
			DX = -sprites[0][len(sprites[0])-1].Frame().Max.X
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
			for l := lpc * pSize; l < (lpc*pSize)+pSize; l++ {
				for m := i * pSize; m < (i*pSize)+pSize; m++ {
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
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func getTranslationMatrix(vector pixel.Vec) pixel.Matrix {
	return pixel.IM.Moved(vector)
}

func importSpriteMatrix() [][]pixel.Sprite {
	var sprites [][]pixel.Sprite
	var index [][]string
	index = readFromFile("export/chunks/matrix.bin", index).([][]string)
	for i := range index {
		sprites = append(sprites, []pixel.Sprite{})
		for j := range index[i] {
			img, err := loadPicture("export/chunks/" + index[i][j])
			if err != nil {
				fmt.Println("Could not open chunk file! (It may not exist, you may not have exported anything, one or more pieces could be missing, must be in /export/chunks)")
				return nil
			}
			sprites[i] = append(sprites[i], *pixel.NewSprite(img, img.Bounds()))
		}
	}
	return sprites
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
	fmt.Println("called") // TODO remove this when done testing
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
