package gv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/emicklei/dot"
)

func TestNewNodeOptions(t *testing.T) {
	tests := []struct {
		opts []nodeAttribute
		want string
	}{
		{
			[]nodeAttribute{},
			`digraph  {n1[fontname="Fira Code",fontsize="12",label="",margin="0.2,0.2",width="2"];}`,
		},

		{
			[]nodeAttribute{nodeFillColor("#ff0000")},
			`digraph  {n1[fillcolor="#ff0000",fontname="Fira Code",fontsize="12",label="",margin="0.2,0.2",style="filled",width="2"];}`,
		},

		{
			[]nodeAttribute{nodeShape("box")},
			`digraph  {n1[fontname="Fira Code",fontsize="12",label="",margin="0.2,0.2",shape="box",width="2"];}`,
		},

		{
			[]nodeAttribute{nodeShape("hexagon"), nodeFillColor("#00ff00")},
			`digraph  {n1[fillcolor="#00ff00",fontname="Fira Code",fontsize="12",label="",margin="0.2,0.2",shape="hexagon",style="filled",width="2"];}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			gv := dot.NewGraph()
			createNode(gv, "ID", tt.opts...)
			if got := flatten(gv.String()); got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

func TestNodeHTMLabel(t *testing.T) {
	tests := []struct {
		html string
		want string
	}{
		{
			`<b>Bold</b>`,
			`digraph  {n1[fontname="Fira Code",fontsize="12",label=<<b>Bold</b>>,margin="0.2,0.2",width="2"];}`,
		},
		{
			`<table><tr><td>col 1</td></tr></table>`,
			`digraph  {n1[fontname="Fira Code",fontsize="12",label=<<table><tr><td>col 1</td></tr></table>>,margin="0.2,0.2",width="2"];}`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			gv := dot.NewGraph()
			createNode(gv, "ID", nodeLabel(tt.html, true))
			if got := flatten(gv.String()); got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}
}

// remove tabs and newlines and spaces
func flatten(s string) string {
	return strings.Replace((strings.Replace(s, "\n", "", -1)), "\t", "", -1)
}
