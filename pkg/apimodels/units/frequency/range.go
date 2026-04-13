package frequency

// UnitString returns a frequency unit string (e.g. "Hz", "GHz").
func (f Frequency) UnitString() string {
	return f.unit.String()
}

// CloneByValue returns a new cloned frequency.
func (f Frequency) CloneByValue() Frequency {
	return *f.Clone()
}

// ParseString parses a frequency string. It is recommended to use the
// [ParseString] function in the same package instead of this method.
func (Frequency) ParseString(s string) (Frequency, error) {
	return ParseString(s)
}
