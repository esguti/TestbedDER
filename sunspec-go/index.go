package sunspec

import "sort"

// Index describes the locality in a modbus address space model.
type Index interface {
	Address() uint16
	Quantity() uint16
}

// index is a helper for instantiating an Index interface.
type index struct {
	address  uint16
	quantity uint16
}

// Address returns the inclusive starting address
func (i index) Address() uint16 {
	return i.address
}

// Quantity returns the number of registers
func (i index) Quantity() uint16 {
	return i.quantity
}

// ceil calculates the ceiling of an index.
// This is the highest inclusive address in a modbus address space.
func ceil(idx Index) uint16 {
	return idx.Address() + idx.Quantity()
}

// intersect determines whether two indexes are overlapping.
func intersect(a, b Index) bool {
	return !(ceil(a) <= b.Address() || ceil(b) <= a.Address())
}

// merge combines all overlapping indexes into their least common sum.
// For instance given the following indexes:
//  A{address: 0; quantity: 4}
//  B{address: 3; quantity: 3}
//  C{address: 8: quantity: 2}
// Then A and B are intersecting each other while C is not contained.
//  So the resulting merged index will be
//  I{address: 0; quantity: 7}
//  J{address: 8; quantity: 2}
func merge(idx []Index) []Index {
	sort.Slice(idx, func(i, j int) bool { return idx[i].Address() < idx[j].Address() })
	var merged []Index
	curr := index{address: idx[0].Address(), quantity: idx[0].Quantity()}
	for _, idx := range idx[1:] {
		switch {
		case intersect(curr, idx):
			curr.quantity = ceil(idx) - curr.address
		default:
			merged = append(merged, curr)
			curr = index{address: idx.Address(), quantity: idx.Quantity()}
		}
	}
	merged = append(merged, curr)
	return merged
}
