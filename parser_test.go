package crumbs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLines(t *testing.T) {
	test := `
* main idea
** topic 1
*** sub topic 1 1
*** sub topic 1 2
**** sub sub topic
** topic 2
*** sub topic 2 1
`
	got, err := ParseLines(strings.SplitAfter(test, "\n"), "")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(got.childrens))
	assert.Equal(t, got.childrens[0].text, "main idea")

	assert.Equal(t, 2, len(got.childrens[0].childrens))
	assert.Equal(t, got.childrens[0].childrens[0].text, "topic 1")
	assert.Equal(t, got.childrens[0].childrens[1].text, "topic 2")

	assert.Equal(t, 2, len(got.childrens[0].childrens[0].childrens))
	assert.Equal(t, got.childrens[0].childrens[0].childrens[1].text, "sub topic 1 2")
}

func TestLookForIcon(t *testing.T) {

	tests := []struct {
		imagespath string
		entry      Entry
		want       string
	}{
		{
			"./images/png",
			Entry{text: "[[blob.png]] La vispa Teresa avea tra l'erbetta, a volo sorpresa"},
			"images/png/blob.png",
		},
		{
			"/home/lus/Pictures/fontawesome/PNG",
			Entry{text: "[[blob.png]] La vispa Teresa avea tra l'erbetta, a volo sorpresa"},
			"/home/lus/Pictures/fontawesome/PNG/blob.png",
		},
	}

	for _, tt := range tests {
		fn := lookForIcon(tt.imagespath)
		fn(&tt.entry)

		t.Run(tt.imagespath, func(t *testing.T) {
			if got := tt.entry.icon; got != tt.want {
				t.Errorf("got [%v] want [%v]", got, tt.want)
			}
		})
	}

}

func TestDepth(t *testing.T) {
	tests := []struct {
		line string
		want int
	}{
		{"* main idea LV.1", 1},
		{"** topic 1 LV.2", 2},
		{"*** sub topic LV.3", 3},
		{"******* LV.7", 7},
	}

	for _, tt := range tests {

		t.Run(tt.line, func(t *testing.T) {
			assert.Equal(t, depth(tt.line), tt.want)
		})
	}
}
