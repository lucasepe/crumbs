package gv

import (
	"fmt"
	"io"
	"strings"

	"github.com/emicklei/dot"
	"github.com/lucasepe/crumbs"
	"github.com/lucasepe/crumbs/text"
)

// RenderConfig defines some render parameters.
type RenderConfig struct {
	VerticalLayout bool
	WrapTextLimit  uint
}

// Render translates the mind note tree to a
// graphviz dot language definition.
func Render(wr io.Writer, note *crumbs.Entry, cfg RenderConfig) error {
	htmlize := htmlLabelMaker(cfg.WrapTextLimit)

	gr := newGraph(Vertical(cfg.VerticalLayout))

	renderTree(gr, note.Root(), htmlize)

	_, err := io.WriteString(wr, gr.String())
	return err
}

// render a tree node (the node, and its children)
func renderTree(gr *dot.Graph, el *crumbs.Entry, htmlize func(*crumbs.Entry) string) {
	tintFor := colorSupplier()

	if el.Level() > 0 {
		createNode(gr, el.ID(), nodeLabel(htmlize(el), true))
	}

	if el.Parent() != nil {
		createEdge(gr, el.Parent().ID(), el.ID(), tintFor(el.Level()))
	}

	for _, child := range el.Childrens() {
		renderTree(gr, child, htmlize)
	}
}

func colorSupplier() func(lvl int) string {
	palette := map[int]string{
		0: "#264653",
		1: "#2A9D8F",
		2: "#E9C46A",
		3: "#E76F51",
		4: "#FFCDB2",
		5: "#B5838D",
		6: "#6D6875",
	}

	return func(lvl int) string {
		if val, ok := palette[lvl]; ok {
			return val
		}
		return "#000000"
	}
}

func htmlLabelMaker(lim uint) func(*crumbs.Entry) string {
	escaper := strings.NewReplacer(
		`&`, "&amp;",
		`'`, "&#39;",
		`"`, "&#34;",
	)

	return func(note *crumbs.Entry) string {
		label := strings.TrimSpace(note.Text())
		if lim > 0 {
			label = text.WrapString(label, lim)
		}
		label = escaper.Replace(label)
		label = strings.ReplaceAll(label, "\n", "<br/>")

		var sb strings.Builder
		sb.WriteString(`<table border="0" cellborder="0">`)

		if len(note.Icon()) > 0 {
			sb.WriteString("<tr>")
			fmt.Fprintf(&sb, `<td fixedsize="true" width="48" height="48"><img src="%s" /></td>`, note.Icon())
			sb.WriteString("</tr>")
		}

		switch {
		case note.Level() == 1:
			fmt.Fprintf(&sb, `<tr><td><font point-size="14"><b>%s</b></font></td></tr>`, label)
		case note.Level() > 1:
			fmt.Fprintf(&sb, `<tr><td><font point-size="12">%s</font></td></tr>`, label)
		}

		sb.WriteString("</table>")

		return sb.String()
	}
}
