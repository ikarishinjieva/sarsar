package ui

import (
	"github.com/jroimartin/gocui"
	"fmt"
	"bytes"
	"sync"
	"strings"
	"regexp"
)

const PREFIX_LEVEL_INDENT = "  "
const PREFIX_COLLAPSE = "+ "
const PREFIX_LEAF = ". "
const PREFIX_EXPAND = "- "

type TreeNode struct {
	Name          string
	Nodes         []*TreeNode
	bindKeyOnce   sync.Once
	enterCallback TreeNodeEnterCallbackFn
	isExpand      bool
	HideName      bool
}

func (n *TreeNode) AddSubNode(name string, nodes []*TreeNode) {
	n.Nodes = append(n.Nodes, &TreeNode{Name: name, Nodes: nodes})
}

func (n *TreeNode) Render(g *gocui.Gui, view *gocui.View) error {
	view.Clear()
	n.bindKey(g, view)

	output := n.innerRender()
	fmt.Fprintf(view, "%s", output)

	return nil
}

func (n *TreeNode) bindKey(g *gocui.Gui, v *gocui.View) error {
	n.bindKeyOnce.Do(func() {
		g.SetKeybinding(v.Name(), gocui.KeyArrowDown, gocui.ModNone, n.onCursorDown)
		g.SetKeybinding(v.Name(), gocui.KeyArrowUp, gocui.ModNone, n.onCursorUp)
		g.SetKeybinding(v.Name(), gocui.KeyEnter, gocui.ModNone, n.onEnter)
	})
	return nil
}

func (n *TreeNode) onCursorDown(g *gocui.Gui, v *gocui.View) error {
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

func (n *TreeNode) onCursorUp(g *gocui.Gui, v *gocui.View) error {
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

func (n *TreeNode) findNode(name string) *TreeNode {
	for _, node := range n.Nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func (n *TreeNode) Expand() {
	n.isExpand = true
}

func (n *TreeNode) Switch() {
	n.isExpand = !n.isExpand
}

func (n *TreeNode) onEnter(g *gocui.Gui, v *gocui.View) error {
	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		line = ""
	}

	lineTrimSpace := strings.TrimLeft(line, PREFIX_LEVEL_INDENT)
	level := (len(line) - len(lineTrimSpace)) / 2

	var segs []string
	for i := cy; i >= cy-level; i -- {
		seg, _ := v.Line(i)
		seg = n.getRawLabel(seg)
		segs = append(segs, seg)
	}

	if strings.HasPrefix(lineTrimSpace, PREFIX_COLLAPSE) || strings.HasPrefix(lineTrimSpace, PREFIX_EXPAND) {
		var curr *TreeNode
		for i := len(segs) - 1; i >= 0; i-- {
			if nil == curr {
				curr = n
			} else {
				curr = curr.findNode(segs[i])
			}

			if 0 == i {
				curr.Switch()
			}
		}
		return n.Render(g, v)
	}

	if nil != n.enterCallback {
		n.enterCallback(g, v, segs)
	}

	return nil
}

func (n *TreeNode) getRawLabel(l string) string {
	for {
		trimmed := l
		trimmed = strings.TrimPrefix(trimmed, PREFIX_LEVEL_INDENT)
		trimmed = strings.TrimPrefix(trimmed, PREFIX_COLLAPSE)
		trimmed = strings.TrimPrefix(trimmed, PREFIX_LEAF)
		trimmed = strings.TrimPrefix(trimmed, PREFIX_EXPAND)
		if l == trimmed {
			return l
		}
		l = trimmed
	}
}

type TreeNodeEnterCallbackFn func(g *gocui.Gui, v *gocui.View, keys []string) error

func (n *TreeNode) SetEnterCallback(callback TreeNodeEnterCallbackFn) {
	n.enterCallback = callback
}

var regexpLineHeader = regexp.MustCompile("(?m)^([^$])")

func (n *TreeNode) innerRender() string {
	buf := bytes.NewBufferString("")
	if len(n.Nodes) > 0 {
		if n.isExpand {
			if !n.HideName {
				fmt.Fprintln(buf, PREFIX_EXPAND+n.Name)
			}
			for _, subNode := range n.Nodes {
				subNodeStr := subNode.innerRender()
				subNodeStr = regexpLineHeader.ReplaceAllString(subNodeStr, PREFIX_LEVEL_INDENT+"$1")
				fmt.Fprint(buf, subNodeStr)
			}
		} else {
			if !n.HideName {
				fmt.Fprintln(buf, PREFIX_COLLAPSE+n.Name)
			}
		}
	} else {
		if !n.HideName {
			fmt.Fprintln(buf, PREFIX_LEAF+n.Name)
		}
	}
	return buf.String()
}
