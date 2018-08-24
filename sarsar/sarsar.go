package sarsar

import (
	"github.com/jroimartin/gocui"
	"github.com/miguelmota/cointop/pkg/table"
	"fmt"
	"github.com/ikarishinjieva/sarsar/sarsar/ui"
)

var file *sarFile

func SarSar(inputFile string) error {
	var err error
	file, err = parseSarFile(inputFile)
	if nil != err {
		return err
	}

	return startUi()
}

func startUi() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if nil != err {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); nil != err {
		return err
	}

	if err := g.MainLoop(); nil != err && err != gocui.ErrQuit {
		return err
	}

	return nil
}

const (
	MENU_WIDTH = 30
)

func layout(g *gocui.Gui) error {
	_, maxY := g.Size()

	if v, err := g.SetView("menu", -1, -1, MENU_WIDTH, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		makeMenuView(g, v)
	}

	g.SetCurrentView("menu")
	return nil
}

func makeMenuView(g *gocui.Gui, v *gocui.View) error {
	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack

	treeRoot := &ui.TreeNode{
		Name:     "root",
		Nodes:    []*ui.TreeNode{},
		HideName: true,
	}
	treeRoot.Expand()

	for sectionId := range file.sections {
		name := section2Name[sectionId]
		section := file.sections[sectionId]

		var nodes []*ui.TreeNode
		if len(section.records) > 0 {
			for col := range section.records[0].data {
				nodes = append(nodes, &ui.TreeNode{
					Name: col,
				})
			}
		}
		treeRoot.AddSubNode(name, nodes)
	}

	treeRoot.SetEnterCallback(menuEnter)

	if err := treeRoot.Render(g, v); nil != err {
		return err
	}
	return nil
}

func menuEnter(g *gocui.Gui, v *gocui.View, keys []string) error {
	if len(keys) != 3 {
		return fmt.Errorf("unexpected menu key depth: %+v", keys)
	}

	sectionId, err := file.getSectionId(keys[1])
	if nil != err {
		return err
	}

	labels, values, err := file.getDataSeriesByName(keys[1], keys[0])
	if nil != err {
		return err
	}

	if err := renderChartView(g, labels, values); nil != err {
		return err
	}

	return renderTableView(g, sectionId)
}

func renderTableView(g *gocui.Gui, sectionId int) error {
	maxX, maxY := g.Size()

	section := file.sections[sectionId]
	tbl := table.New().SetWidth(maxX)

	if len(section.records) == 0 {
		return nil
	}

	for col := range section.records[0].data {
		tbl.AddCol(fmt.Sprintf("%8s", col))
	}

	for _, rec := range section.records {
		var vals []interface{}
		for _, val := range rec.data {
			vals = append(vals, fmt.Sprintf("%8s", val))
		}
		tbl.AddRow(vals...)
	}

	g.DeleteView("table")
	if v, err := g.SetView("table", MENU_WIDTH+1, 11, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false

		g.Update(func(gui *gocui.Gui) error {
			tbl.Format().Fprint(v)
			return nil
		})
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
