package largenumber

// UnitString returns a largenumber unit string (e.g. "10", "10 K").
func (m LargeNumber) UnitString() string {
	return m.unit.String()
}

// CloneByValue returns a new cloned large number.
func (m LargeNumber) CloneByValue() LargeNumber {
	return *m.Clone()
}

// ParseString parses a large number string. It is recommended to use the
// [ParseString] function in the same package instead of this method.
func (LargeNumber) ParseString(s string) (LargeNumber, error) {
	return ParseString(s)
}
