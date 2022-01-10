package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/dusk125/pixelutils"
	"github.com/inkyblackness/imgui-go"
	_ "image/png"
	_ "runtime"
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

var win *pixelgl.Window
var ticker *pixelutils.Ticker
var framerateTarget = int64(120)

var white, _ = colorful.MakeColor(colornames.Antiquewhite)
var black, _ = colorful.MakeColor(colornames.Black)

var sprites [][]pixel.Sprite

var logBuf bytes.Buffer

var generationProgress = float32(0.0)
var threadProgress = float32(0.0)
var totalProgress = float32(0.0)
var gridSize = float32(0.0)

var exporting = false
var imported = -1 // if -1 then its not being imported and is not in memory, if 0 then its being imported, if its 1 then its imported

var simType = 0
var evolutions = int32(256)
var history = true
var path = ""
var leftPart = true

var randomInitialLength = int32(256)

func run() {
	initWindow()

	var (
		camPos       = pixel.ZV
		camSpeed     = 800.0
		camZoom      = 1.0
		camZoomSpeed = 1.2
	)

	ui := pixelui.NewUI(win, 0)
	defer ui.Destroy()

	gob.Register([][]string{})

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		cam := pixel.IM.Scaled(camPos, camZoom).Moved(win.Bounds().Center().Sub(camPos))
		win.SetMatrix(cam)

		// Quit
		if win.JustReleased(pixelgl.KeyEscape) {
			win.SetClosed(true)
		}

		// Controls
		if !imgui.CurrentIO().WantCaptureMouse() && win.Pressed(pixelgl.MouseButtonLeft) {
			camPos = camPos.Sub(win.MousePosition().Sub(win.MousePreviousPosition()).Scaled(1 / camZoom))
		}
		if win.Pressed(pixelgl.KeyLeft) {
			camPos.X -= camSpeed * dt * 1 / camZoom
		}
		if win.Pressed(pixelgl.KeyRight) {
			camPos.X += camSpeed * dt * 1 / camZoom
		}
		if win.Pressed(pixelgl.KeyDown) {
			camPos.Y -= camSpeed * dt * 1 / camZoom
		}
		if win.Pressed(pixelgl.KeyUp) {
			camPos.Y += camSpeed * dt * 1 / camZoom
		}
		if !imgui.CurrentIO().WantCaptureMouse() {
			camZoom *= math.Pow(camZoomSpeed, win.MouseScroll().Y)
		}

		ui.NewFrame()
		imgui.ShowDemoWindow(nil)

		//TODO imgui_demo.cpp
		// TODO import a batch of all the chunks to use for the translation matrix

		win.Clear(colornames.Antiquewhite)
		drawUI()
		drawSpriteMatrix()
		//sprite.Draw(win, pixel.IM.Moved(win.Bounds().Max.ScaledXY(pixel.V(0.5, 0.5)).Sub(sprite.Frame().Center())))

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
	imgui.Spacing()
	imgui.ProgressBarV(generationProgress/float32(len(sim)), imgui.Vec2{X: 375, Y: 22}, strconv.Itoa(int(generationProgress))+"/"+strconv.Itoa(len(sim)))
	imgui.SameLine()
	imgui.Text("Generation")
	imgui.SameLine()
	imgui.Text("       Rule 110 Calculator")
	imgui.ProgressBarV(threadProgress/float32(len(sim)), imgui.Vec2{X: 375, Y: 22}, strconv.Itoa(int(threadProgress))+"/"+strconv.Itoa(len(sim)))
	imgui.SameLine()
	imgui.Text("Threads")
	imgui.SameLine()
	imgui.Text("           By:Julian Carrier")
	imgui.ProgressBarV(totalProgress/gridSize, imgui.Vec2{X: 375, Y: 22}, strconv.Itoa(int(totalProgress))+"/"+strconv.Itoa(int(gridSize)))
	imgui.SameLine()
	imgui.Text("Total Progress")
	imgui.SameLine()
	imgui.Text("      FPS:")
	imgui.SameLine()
	imgui.Text(fmt.Sprintf("%.2f", framerate))
	imgui.Spacing()
	imgui.Separator()
	imgui.Spacing()

	if imgui.Button("Generate") {
		addLog("[Activity]:" + time.Now().Local().String() + ": Starting Generation... \n")
		generationProgress = 0.0
		switch simType {
		case 0:
			simulate(history, int(evolutions), randStart(int(randomInitialLength))) // Random
		case 1:
			simulate(history, int(evolutions), r110Default()) // Template Start
		case 2:
			simulate(history, int(evolutions), simpleCTS()) // CTS
			// TODO create function here to define cts configuration
		case 3:
			simulate(history, int(evolutions), r110Default())
			// TODO change to the input of the user
		}
		addLog("[Activity]:" + time.Now().Local().String() + ": Generation Complete... \n")
	}
	imgui.SameLine()
	if imgui.Checkbox("History", &history) {

	}
	imgui.SameLine()
	imgui.InputInt("Evolutions", &evolutions)
	if imgui.Button("Export") && exporting == false {
		addLog("[Activity]:" + time.Now().Local().String() + ": Starting Image Export... \n")
		exporting = true
		go createImages(true, 2)
		path = "/export/chunks"
	}
	imgui.SameLine()
	if imgui.Button("Import") && (imported == -1 || imported == 1) && exporting == false {
		addLog("[Activity]:" + time.Now().Local().String() + ": Importing Images... \n")
		imported = 0
		go importSpriteMatrix()
	}
	imgui.SameLine()
	if imgui.Checkbox("Left", &leftPart) { // TODO change to not include left part of simulation if unchecked

	}
	imgui.SameLine()
	imgui.InputText("Path", &path)

	imgui.Spacing()

	imgui.RadioButtonInt("Random", &simType, 0)
	imgui.SameLine()
	imgui.RadioButtonInt("Template", &simType, 1)
	imgui.SameLine()
	imgui.RadioButtonInt("CTS", &simType, 2)
	imgui.SameLine()
	imgui.RadioButtonInt("Custom", &simType, 3)
	switch simType {
	case 0:
		imgui.InputInt("Initial Length", &randomInitialLength)
	case 2:
		imgui.Selectable("test")
	case 3:
		imgui.Selectable("test1")
	}

	if imgui.CollapsingHeader("Logging") {
		if imgui.Button("Clear") {
			logBuf.Reset()
		}
		imgui.SameLine()
		if imgui.Button("Copy") {
			win.SetClipboard(logBuf.String())
		}
		imgui.Separator()
		imgui.BeginChildV("scrolling", imgui.Vec2{0, 200}, true, imgui.WindowFlagsHorizontalScrollbar)
		imgui.Text(logBuf.String())
		imgui.EndChild()
	}
	//imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{0,0})
	////logBuf.ReadString(b'\n')

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

func addLog(str string) {
	fmt.Print(str)
	logBuf.WriteString(str)
}

func createImages(history bool, pSize int) {
	threadProgress = 0.0
	totalProgress = 0.0
	gridSize = 0.0

	if sim == nil || len(sim) == 0 {
		addLog("[Error]:" + time.Now().Local().String() + ": No Simulation to Export (Nothing in Memory)... \n")
		exporting = false
		return
	}
	width := len(sim[0]) * pSize * 64
	height := len(sim) * pSize

	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	var wg sync.WaitGroup
	wg.Add(len(sim))
	if history {
		img = image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
		addLog("[Activity]:" + time.Now().Local().String() + ": Passing Layers to Threads... \n")
		for i := range sim {
			go renderRow(&wg, img, i, pSize)
			threadProgress++
		}
	} else {
		// TODO add non history version
	}
	wg.Wait()

	addLog("[Activity]:" + time.Now().Local().String() + ": Creating Files... \n")
	_ = os.RemoveAll("export/chunks/backup")
	_ = os.MkdirAll("export/chunks/backup", os.ModePerm)
	dir, _ := os.ReadDir("export/chunks")
	for entry := range dir {
		_ = os.Rename("export/chunks/"+dir[entry].Name(), "export/chunks/backup/"+dir[entry].Name())
	}
	f, _ := os.Create("export/image.png")
	addLog("[Activity]:" + time.Now().Local().String() + ": Files Created Successfully... \n")

	grid := gridSplit(img.Rect)
	addLog("[Activity]:" + time.Now().Local().String() + ": Starting Image Export... \n")
	_ = png.Encode(f, img)
	addLog("[Activity]:" + time.Now().Local().String() + ": Initial Export Done... \n")

	if grid == nil {
		addLog("[Activity]:" + time.Now().Local().String() + ": Image is already under maximum texture size... \n")
		gridSize = 1
		totalProgress = 1
	} else {
		gridSize = float32(len(grid)*len(grid[0])) + 1
		totalProgress++

		exportMatrix := make([][]string, len(grid))
		addLog("[Activity]:" + time.Now().Local().String() + ": Splitting Image... \n")
		for i := range exportMatrix {
			exportMatrix[i] = make([]string, len(grid[0]))
		}

		addLog("[Activity]:" + time.Now().Local().String() + ": Exporting Sub Images... \n")
		for i := range grid {
			for j := range grid[i] {
				name := strconv.Itoa(i) + "_" + strconv.Itoa(j)
				exportMatrix[i][j] = name + ".png"
				f, _ = os.Create("export/chunks/" + name + ".png")
				_ = png.Encode(f, img.SubImage(grid[i][j]))
				totalProgress++
			}
		}

		addLog("[Activity]:" + time.Now().Local().String() + ": Exporting Configuration... \n")
		writeToFile("export/chunks/matrix.bin", exportMatrix)
	}
	addLog("[Activity]:" + time.Now().Local().String() + ": Image Generation and Export Complete... \n")
	exporting = false
}

func drawSpriteMatrix() {
	if sprites != nil && imported == 1 {
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

func importSpriteMatrix() {
	sprites = nil
	var temp interface{}
	temp, err := readFromFile("export/chunks/matrix.bin", temp)
	index, _ := temp.([][]string)
	if errors.Is(err, os.ErrNotExist) {
		addLog("[Warn]:" + time.Now().Local().String() + ": No Matrix configuration file found check for matrix.bin file in /export/chunks... \n")
		img, err := loadPicture("export/image.png")
		if errors.Is(err, os.ErrNotExist) {
			addLog("[Error]:" + time.Now().Local().String() + ": No alternative image found import failed... \n")
			return
		} else {
			sprites = append(sprites, []pixel.Sprite{})
			sprites[0] = append(sprites[0], *pixel.NewSprite(img, img.Bounds()))
			imported = 1
			addLog("[Activity]:" + time.Now().Local().String() + ": Images Imported Successfully... \n")
		}
	} else {
		for i := range index {
			sprites = append(sprites, []pixel.Sprite{})
			for j := range index[i] {
				img, err := loadPicture("export/chunks/" + index[i][j])
				if err != nil {
					addLog("[Error]:" + time.Now().Local().String() + ": Could not open chunk file! (It may not exist, you may not have exported anything, one or more pieces could be missing, must be in /export/chunks)... \n")
					sprites = nil
				}
				sprites[i] = append(sprites[i], *pixel.NewSprite(img, img.Bounds()))
			}
		}
		imported = 1
		addLog("[Activity]:" + time.Now().Local().String() + ": Images Imported Successfully... \n")
	}

} // TODO add nessecary checks to make sure that the program just imports the image if no matrix was nessecary
//TODO NOT SPLITTING IMAGE CORRECTLY FOR CTS EVOLUTIONS 500 FIND OUT WHY

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
