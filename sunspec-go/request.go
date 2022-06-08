package sunspec

// Request describes a received sunspec server request.
type Request interface {
	// Writing specifies whether the request is attempting to set point values.
	Writing() bool
	// Ingest updates the affected point values in accordance to the request.
	// For read only requests no change is applied to the points.
	Ingest() error
	// Points returns all points that are affected by the request.
	Points() Points
	// Flush ends the request.
	// It is mandatory to do so after finishing the processing.
	Flush() error
}

type request struct {
	points  Points
	writing bool
	buffer  []byte
}

// Writing specifies whether the request is attempting to set point values.
func (r *request) Writing() bool { return r.writing }

// Ingest updates the affected point values in accordance to the request.
// For read only requests no change is applied to the points.
func (r *request) Ingest() error {
	if !r.Writing() {
		return nil
	}
	return r.points.decode(r.buffer)
}

// Points returns all points that are affected by the request.
func (r *request) Points() Points { return r.points.Points() }

// Close ends the request.
// It is mandatory to do so after finishing the processing.
func (r *request) Flush() error {
	return r.points.encode(r.buffer)
}
