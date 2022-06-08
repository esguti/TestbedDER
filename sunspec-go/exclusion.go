package sunspec

import (
	"encoding/json"
	"errors"
)

type exclusion interface {
	false() string
	true() string
}

func marshalExclusion(sel bool, e exclusion) ([]byte, error) {
	if sel {
		return []byte(`"` + e.true() + `"`), nil
	}
	return []byte(`"` + e.false() + `"`), nil
}

func unmarshalExclusion(b []byte, e exclusion) (bool, error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return false, err
	}
	switch s {
	case e.false():
		return false, nil
	case e.true():
		return true, nil
	}
	return false, errors.New("sunspec: illegal enumerated value")
}

type atomic bool

func (t atomic) false() string { return "group" }

func (t atomic) true() string { return "sync" }

// MarshalJSON puts the enumerated string-value in a json document.
func (t atomic) MarshalJSON() ([]byte, error) {
	return marshalExclusion(bool(t), t)
}

// UnmarshalJSON sets atomic from its enumerated json string-value.
func (t *atomic) UnmarshalJSON(b []byte) error {
	v, err := unmarshalExclusion(b, t)
	*t = atomic(v)
	return err
}

type writable bool

func (t writable) false() string { return "R" }

func (t writable) true() string { return "RW" }

// MarshalJSON puts the enumerated string-value in a json document.
func (t writable) MarshalJSON() ([]byte, error) {
	return marshalExclusion(bool(t), t)
}

// UnmarshalJSON sets writable from its enumerated json string-value.
func (t *writable) UnmarshalJSON(b []byte) (err error) {
	v, err := unmarshalExclusion(b, t)
	*t = writable(v)
	return err
}

type mandatory bool

func (t mandatory) false() string { return "O" }

func (t mandatory) true() string { return "M" }

// MarshalJSON puts the enumerated string-value in a json document.
func (t mandatory) MarshalJSON() ([]byte, error) {
	return marshalExclusion(bool(t), t)
}

// UnmarshalJSON sets mandatory from its enumerated json string-value.
func (t *mandatory) UnmarshalJSON(b []byte) (err error) {
	v, err := unmarshalExclusion(b, t)
	*t = mandatory(v)
	return err
}

type static bool

func (t static) false() string { return "D" }

func (t static) true() string { return "S" }

// MarshalJSON puts the enumerated string-value in a json document.
func (t static) MarshalJSON() ([]byte, error) {
	return marshalExclusion(bool(t), t)
}

// UnmarshalJSON sets static from its enumerated json string-value.
func (t *static) UnmarshalJSON(b []byte) (err error) {
	v, err := unmarshalExclusion(b, t)
	*t = static(v)
	return err
}
