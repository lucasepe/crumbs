package gv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/emicklei/dot"
	"github.com/lucasepe/crumbs/text"
)

// GraphOption is a graph attribute.
type GraphOption func(*dot.Graph)

// ImagesPath specifies a list of directories in which
// to look for image files as specified by the
// image attribute or using the IMG element in HTML-like labels.
// The string should be a list of (absolute or relative)
// pathnames, each separated by a semicolon (for Windows)
// or a colon (all other OS).
func ImagesPath(paths string) GraphOption {
	return func(gr *dot.Graph) {
		gr.Attr("imagepath", paths)
	}
}

// Vertical enables/diables the top to bottom layout direction.
func Vertical(set bool) GraphOption {
	return func(gr *dot.Graph) {
		if set {
			gr.Attr("rankdir", "TB")
		} else {
			gr.Attr("rankdir", "LR")
		}
	}
}

// newGraph returns a new GraphViz DOT language graph
func newGraph(opts ...GraphOption) *dot.Graph {
	res := dot.NewGraph(dot.Undirected)

	res.Attr("rankdir", "LR")
	res.Attr("pad", "1")
	res.Attr("ranksep", "2.3")
	res.Attr("nodesep", "0.8")
	res.Attr("fontname", "Fira Code")
	res.Attr("fontsize", "14")
	res.Attr("splines", "curved")
	res.Attr("concentrate", "true")
	res.Attrs("orientation", "portrait")

	for _, opt := range opts {
		opt(res)
	}

	return res
}

// createEdge creates a new connection line between two nodes
func createEdge(gr *dot.Graph, fid, tid string, color string) error {
	a, ok := gr.FindNodeById(fid)
	if !ok {
		return fmt.Errorf("node with id=%s not found", fid)
	}

	b, ok := gr.FindNodeById(tid)
	if !ok {
		return fmt.Errorf("node with id=%s not found", tid)
	}

	res := gr.Edge(a, b)

	res.Attr("fontname", "Fira Code")
	res.Attr("fontsize", "10")
	res.Attr("penwidth", "2.5")
	//res.Attr("xlabels", strconv.Itoa(lvl))

	if strings.TrimSpace(color) != "" {
		res.Attr("color", color)
	} else {
		res.Attr("color", "#ced4da")
	}

	return nil
}

// createNode create and adds a new node to the graph.
// You can customize some attributes using
// the variadic node attributes.
func createNode(g *dot.Graph, id string, opts ...nodeAttribute) string {
	n := g.Node(id)
	n.Label("")
	n.Attr("fontname", "Fira Code")
	n.Attr("fontsize", "12")
	n.Attr("width", "2")
	n.Attr("margin", "0.2,0.2")
	n.Attr("shape", "plain")

	for _, opt := range opts {
		opt(&n)
	}

	return id
}

// nodeAttribute defines a function that
// apply a property to a node
type nodeAttribute func(*dot.Node)

// nodeFillColor sets the node fill color
func nodeFillColor(hex string) nodeAttribute {
	return func(el *dot.Node) {
		if strings.TrimSpace(hex) == "" {
			return
		}

		el.Attr("fillcolor", hex)

		const attr = "filled"

		val := el.AttributesMap.Value("style")
		if val == nil {
			el.Attr("style", attr)
			return
		}

		style, ok := val.(string)
		if !ok {
			el.Attr("style", attr)
			return
		}

		if _, found := text.Find(strings.Split(style, ","), attr); !found {
			el.Attr("style", strings.Join([]string{style, attr}, ","))
		}
	}
}

// nodeFontSize specify the font size, in points, used for text
func nodeFontSize(s int) nodeAttribute {
	return func(el *dot.Node) {
		el.Attr("fontsize", strconv.Itoa(s))
	}
}

// nodeLabel is the node caption
// if 'htm' is true the caption is treated as HTML code
func nodeLabel(label string, htm bool) nodeAttribute {
	return func(el *dot.Node) {
		if htm {
			el.Attr("label", dot.HTML(label))
		} else {
			el.Attr("label", label)
		}
	}
}

// nodeShape sets the shape of a node
func nodeShape(shape string) nodeAttribute {
	return func(el *dot.Node) {
		el.Attr("shape", shape)
	}
}
