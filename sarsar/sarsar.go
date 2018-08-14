package sarsar

import (
	"github.com/jroimartin/gocui"
	"github.com/gizak/termui"
	"github.com/miguelmota/cointop/pkg/color"
	"github.com/miguelmota/cointop/pkg/table"
	"log"
	"fmt"
)

func HelloWorld() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

const (
	CHART_HEIGHT = 10
	MENU_WIDTH = 30
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("menu", -1, -1, MENU_WIDTH, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		makeMenuView(g, v)
	}

	if v, err := g.SetView("chart", MENU_WIDTH+1, 0, maxX-1, CHART_HEIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		makeChartView(g, v, maxX-2, 10, []float64{1001.0, 2000.0, 3000.0, 4000.0, 5000.0, 6000.0, 5000.0, 4000.0, 3000.0, 2.0, 1000.0, 2.5})
	}

	if v, err := g.SetView("table", MENU_WIDTH+1, 11, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		makeTableView(g, v, maxX)
	}

	g.SetCurrentView("menu")
	return nil
}


func makeMenuView(g *gocui.Gui, v *gocui.View) error {
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack
	fmt.Fprintln(v, "Item 1")
	fmt.Fprintln(v, "Item 2")
	fmt.Fprintln(v, "Item 3")

	if err := g.SetKeybinding(v.Name(), gocui.KeyArrowDown, gocui.ModNone, menuCursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding(v.Name(), gocui.KeyArrowUp, gocui.ModNone, menuCursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(v.Name(), gocui.KeyEnter, gocui.ModNone, menuEnter); err != nil {
		return err
	}
	return nil
}

func menuCursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func menuCursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func menuEnter(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, _ := g.Size()

	g.DeleteView("chart")

	if v, err := g.SetView("chart", MENU_WIDTH, 0, maxX-1, CHART_HEIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		makeChartView(g, v, maxX-2, 10, []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 5.0, 4.0, 3.0, 2.0, 1000.0, 2.5})
	}

	l = l

	return nil
}

func makeTableView(g *gocui.Gui, view *gocui.View, maxX int) {
	t := table.New().SetWidth(maxX)
	t.AddCol("first")
	t.AddCol("second")
	t.AddRow("1", "2")
	t.AddRow("3", "4")

	g.Update(func(gui *gocui.Gui) error {
		t.Format().Fprint(view)
		return nil
	})
}

func makeChartView(g *gocui.Gui, view *gocui.View, maxX int, height int, data []float64) {
	chartPoints := makeChartPoints(maxX, height, data)

	var body string
	for i := range chartPoints {
		var s string
		for j := range chartPoints[i] {
			p := chartPoints[i][j]
			s = fmt.Sprintf("%s%c", s, p.Ch)
		}
		body = fmt.Sprintf("%s%s\n", body, s)
	}

	g.Update(func(gui *gocui.Gui) error {
		fmt.Fprint(view, color.White(body))
		return nil
	})

}

func makeChartPoints(maxX int, height int, data []float64) [][]termui.Cell {
	chart := termui.NewLineChart()
	chart.Height = height
	chart.Width = maxX
	chart.AxesColor = termui.ColorWhite
	chart.LineColor = termui.ColorCyan
	chart.Border = false
	chart.Mode = "dot"
	chart.DotStyle = rune('.')
	chart.Data = data

	termui.Body = termui.NewGrid()
	termui.Body.Width = maxX
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(12, 0, chart),
		),
	)

	var points [][]termui.Cell

	{
		// calculate layout
		termui.Body.Align()
		w := termui.Body.Width
		h := height
		row := termui.Body.Rows[0]
		b := row.Buffer()
		for i := 0; i < h; i = i + 1 {
			var rowPoints []termui.Cell
			for j := 0; j < w; j = j + 1 {
				p := b.At(j, i)
				rowPoints = append(rowPoints, p)
			}
			points = append(points, rowPoints)
		}
	}

	return points
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
