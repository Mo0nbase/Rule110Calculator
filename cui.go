package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
)

var cont = true

func maind5gd() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	data := []float64{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}

	sl3 := widgets.NewSparkline()
	sl3.Title = "Enlarged Sparkline"
	sl3.Data = data
	sl3.LineColor = ui.ColorYellow

	slg2 := widgets.NewSparklineGroup(sl3)
	slg2.Title = "Tweeked Sparkline"
	slg2.SetRect(20, 0, 100, 15)
	slg2.BorderStyle.Fg = ui.ColorCyan

	ui.Render(slg2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		//if cont {
		//	rand.Seed(time.Now().UnixNano())
		//	rand.Shuffle(len(data), func(i, j int) { data[i], data[j] = data[j], data[i] })
		//	sl3.Data = data
		//	ui.Render(slg2)
		//}
		switch e.ID {
		case "q", "<C-c>":
			return
			//case "s":
			//	cont = !cont
		}
	}
}
