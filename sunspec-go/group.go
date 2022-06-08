package sunspec

// Group defines a sunspec container for points.
type Group interface {
	// Index defines the locality of the entire group in a modbus address space.
	Index
	// Name returns the group`s identifier.
	Name() string
	// Atomic specifies if all immediate points in the group must be read atomically.
	Atomic() bool
	// Origin returns the group´s parent container.
	Origin() Group
	// Point returns the first immediate point identified by name.
	Point(name string) Point
	// Points returns all immediate points identified by names.
	// If names are omitted all points are returned.
	Points(names ...string) Points
	// Group returns the first immediate group identified by name.
	Group(name string) Group
	// Groups returns all immediate groups identified by names.
	// If names are omitted all groups are returned.
	Groups(names ...string) Groups
}

// GroupDef is the definition of a sunspec Group element.
type GroupDef struct {
	Name        string      `json:"name"`
	Atomic      atomic      `json:"type"`
	Count       interface{} `json:"count"`
	Points      []PointDef  `json:"points"`
	Groups      []GroupDef  `json:"groups"`
	Label       string      `json:"label,omitempty"`
	Description string      `json:"desc,omitempty"`
	Detail      string      `json:"detail,omitempty"`
	Notes       string      `json:"notes,omitempty"`
	Comments    []string    `json:"comments,omitempty"`
}

// iterate executes callback recursively for group g and all its sub-groups.
// The function immediately stops if the callback returns an error.
func iterate(g Group, callback func(g Group) error) error {
	if err := callback(g); err != nil {
		return err
	}
	for _, g := range g.Groups() {
		iterate(g, callback)
	}
	return nil
}

// group is internally used to build out a model.
type group struct {
	name   string
	atomic bool
	origin *group
	points Points
	groups Groups
}

// Address returns the modbus starting address of the given Group.
func (g *group) Address() uint16 { return g.points.address() }

// Quantity returns the number of registers (2-byte/words) of all points in the group including sub-groups.
func (g *group) Quantity() uint16 { return g.points.Quantity() + g.groups.Quantity() }

// Name returns the groups identifier.
func (g *group) Name() string { return g.name }

// Atomic specifies whether the immediate points of the group have to be manipulated atomically.
func (g *group) Atomic() bool { return g.atomic }

// Origin returns the group´s parent container.
func (g *group) Origin() Group { return g.origin }

// Point returns the first immediate point identified by name.
func (g *group) Point(name string) Point { return g.points.Point(name) }

// Points returns all immediate points identified by names.
func (g *group) Points(names ...string) Points { return g.points.Points(names...) }

// Group returns the first immediate group identified by name.
func (g *group) Group(name string) Group { return g.groups.Group(name) }

// Groups returns all immediate groups identified by names.
func (g *group) Groups(names ...string) Groups { return g.groups.Groups(names...) }

// Groups wraps a collection of multiple groups.
type Groups []Group

// First returns the first group from the collection.
func (gps Groups) First() Group { return gps[0] }

// Last returns the last group from the collection.
func (gps Groups) Last() Group { return gps[len(gps)-1] }

// Quantity determines the total word size of the points in the collection including all sub-groups.
func (gps Groups) Quantity() uint16 {
	var l uint16
	for _, g := range gps {
		l += g.Quantity()
	}
	return l
}

// Group returns the first immediate group identified by name.
func (gps Groups) Group(name string) Group {
	for _, g := range gps {
		if g.Name() == name {
			return g
		}
	}
	return nil
}

// Groups returns all immediate groups identified by names.
// If names are omitted all groups are returned.
func (gps Groups) Groups(names ...string) Groups {
	if len(names) == 0 {
		return append(Groups(nil), gps...)
	}
	col := make(Groups, 0, len(names))
	for _, g := range gps {
		for _, id := range names {
			if g.Name() == id {
				col = append(col, g)
				break
			}
		}
	}
	return col
}

// Index returns the merged indexes of all groups in the collection.
func (gps Groups) Index() []Index {
	idx := make([]Index, 0, len(gps))
	for _, g := range gps {
		idx = append(idx, g)
	}
	return merge(idx)
}
