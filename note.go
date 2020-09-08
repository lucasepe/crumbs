package crumbs

// Entry is a thought, an idea.
type Entry struct {
	id        string
	level     int
	text      string
	icon      string
	parent    *Entry
	childrens []*Entry
}

// ID returns the node identifier.
func (ti *Entry) ID() string {
	return ti.id
}

// Childrens returns the node childs.
func (ti *Entry) Childrens() []*Entry {
	return ti.childrens
}

// Text returns the node data.
func (ti *Entry) Text() string {
	return ti.text
}

// Icon returns the icon path.
func (ti *Entry) Icon() string {
	return ti.icon
}

// Level returns the node depth.
func (ti *Entry) Level() int {
	return ti.level
}

// Parent returns the node parent
func (ti *Entry) Parent() *Entry {
	return ti.parent
}

// Root returns the root note
func (ti *Entry) Root() *Entry {
	if ti.parent == nil {
		return ti
	}
	return ti.parent.Root()
}
