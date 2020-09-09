package crumbs

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/teris-io/shortid"
)

// ParseLines parses a slice of text lines and builds the tree.
func ParseLines(lines []string, iconspath string) (*Entry, error) {
	mkID := idGenerator()
	checkIcon := lookForIcon(iconspath)

	// generate a short id for the root node
	rootID, err := mkID()
	if err != nil {
		return nil, err
	}

	// create the root node
	root := newEmptyNote(rootID)

	node := root
	nodeDepth := 0
	for _, el := range lines {
		// skip empty lines
		if strings.TrimSpace(el) == "" {
			continue
		}

		// count depth
		childDepth := depth(el)

		// trim leading 'stars', then the spaces
		text := el[childDepth:]
		text = strings.TrimSpace(text)

		// create the child
		childID, err := mkID()
		if err != nil {
			return nil, err
		}
		child := newNote(childID, childDepth, text)
		// check if has an icon
		checkIcon(child)

		// case: the current 'node' is the parent
		if childDepth > nodeDepth {
			// update tree
			child.parent = node
			node.childrens = append(node.childrens, child)

			// update loop state
			node = child
			nodeDepth++

			// case: the current 'node' is not the parent of our child
		} else if childDepth <= nodeDepth {
			// adjust 'node' until it's correct
			for childDepth <= nodeDepth {
				node = node.parent
				nodeDepth--
			}

			// update tree
			child.parent = node
			node.childrens = append(node.childrens, child)

			// update loop state
			node = child
			nodeDepth++
		}
	}

	return root, nil
}

// depth space-counting helper (probably done in a dumb way, dunno)
func depth(line string) int {
	i := 0
	for line[i] == '*' {
		i++
	}

	return i
}

// newNote creates a new note element
func newNote(id string, lvl int, txt string) *Entry {
	f := new(Entry)
	f.id = id
	f.text = txt
	f.level = lvl
	return f
}

// newNote creates a new note element
func newEmptyNote(id string) *Entry {
	f := new(Entry)
	f.id = id
	f.level = -1
	return f
}

// idGenerator generates a new short id at each invocation.
func idGenerator() func() (string, error) {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return func() (string, error) {
			return "", err
		}
	}

	return func() (string, error) {
		return sid.Generate()
	}
}

func lookForIcon(iconspath string) func(note *Entry) {
	re := regexp.MustCompile(`^\[{2}(.*?)\]{2}`)

	return func(note *Entry) {
		str := note.text
		res := re.FindStringSubmatch(str)
		if len(res) > 0 {
			note.icon = filepath.Join(iconspath, strings.TrimSpace(res[1]))
			note.text = re.ReplaceAllString(str, "")
		}
	}
}
