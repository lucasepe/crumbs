package gv

import (
	"fmt"
	"io"
	"strings"

	"github.com/emicklei/dot"
	"github.com/lucasepe/crumbs"
	"github.com/lucasepe/crumbs/text"
)

// Render translates the mind note tree to a
// graphviz dot language definition.
func Render(wr io.Writer, note *crumbs.Entry, opts ...GraphOption) error {
	gr := newGraph(opts...)
	renderTree(gr, note.Root(), sanitizer(30))
	_, err := io.WriteString(wr, gr.String())
	return err
}

// render a tree node (the node, and its children)
func renderTree(gr *dot.Graph, el *crumbs.Entry, sanify func(string) string) {
	tintFor := colorSupplier()
	asHTML := htmlLabelMaker(30)

	if el.Level() > 0 {
		createNode(gr, el.ID(), nodeLabel(asHTML(el), true))
	}

	if el.Parent() != nil {
		createEdge(gr, el.Parent().ID(), el.ID(), tintFor(el.Level()))
	}

	for _, child := range el.Childrens() {
		renderTree(gr, child, sanify)
	}
}

// sanitizer wraps the given string within lim
// width in characters and does some sanification
func sanitizer(lim uint) func(string) string {
	escaper := strings.NewReplacer(
		`&`, "&amp;",
		`'`, "&#39;",
		`"`, "&#34;",
	)

	return func(txt string) string {
		res := text.WrapString(txt, lim)
		res = escaper.Replace(res)
		res = strings.ReplaceAll(res, "\n", "<br/>")
		return res
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
		label := text.WrapString(note.Text(), lim)
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
