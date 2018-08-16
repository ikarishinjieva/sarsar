package sarsar

import (
	"github.com/jroimartin/gocui"
	"fmt"
	"github.com/gizak/termui"
	"github.com/miguelmota/cointop/pkg/color"
)

const (
	CHART_HEIGHT = 10
)

func renderChartView(g *gocui.Gui, labels []string, values []float64) error {
	maxX, _ := g.Size()

	g.DeleteView("chart")
	if v, err := g.SetView("chart", MENU_WIDTH, 0, maxX-1, CHART_HEIGHT); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		makeChartView(g, v, maxX-2, CHART_HEIGHT, labels, values)
	}
	return nil
}

func makeChartView(g *gocui.Gui, view *gocui.View, maxX int, height int, labels []string, values []float64) {
	chartPoints := makeChartPoints(maxX, height, labels, values)

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

func makeChartPoints(maxX int, height int, labels []string, values []float64) [][]termui.Cell {
	chart := termui.NewLineChart()
	chart.Height = height
	chart.Width = maxX
	chart.AxesColor = termui.ColorWhite
	chart.LineColor = termui.ColorGreen
	chart.Border = false
	chart.Mode = "braille"
	chart.Data = values
	chart.DataLabels = labels

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
