package sunspec

import "errors"

// Device describes a sunspec compliant device.
type Device interface {
	// Model returns the first immediate model identified by id.
	Model(id uint16) Model
	// Models returns all models from the device.
	// If ids are omitted all models are returned.
	Models(ids ...uint16) Models
}

// collect retrieves all the distinct points in a given address range.
// 	notice: invalid ranges may become valid when merged
func collect(d Device, idx ...Index) (Points, error) {
	var pts Points
	for _, idx := range merge(idx) {
		for _, m := range d.Models() {
			if !intersect(idx, m) {
				continue
			}
			if err := iterate(m, func(g Group) error {
				switch {
				case idx.Address() > ceil(g.Points().index()):
				case idx.Address() <= g.Address() && ceil(idx) >= ceil(g.Points().index()):
					pts = append(pts, g.Points()...)
				case g.Atomic():
					return errors.New("sunspec: the operation can not be done for an atomic group")
				default:
					for _, p := range g.Points() {
						if intersect(idx, p) {
							if idx.Address() > p.Address() || ceil(idx) < ceil(p) {
								return errors.New("sunspec: point not fully contained by index")
							}
							pts = append(pts, p)
						}
					}
				}
				return nil
			}); err != nil {
				return nil, err
			}
		}
	}
	if len(pts) == 0 {
		return nil, errors.New("sunspec: index does not reference any points in the device")
	}
	return pts, nil
}
