package bytesize

// UnitString returns a bytesize unit string (e.g. "B", "KB").
func (b ByteSize) UnitString() string {
	return b.unit.String()
}

// CloneByValue returns a new cloned bytesize.
func (b ByteSize) CloneByValue() ByteSize {
	return *b.Clone()
}

// ParseString parses a bytesize string. It is recommended to use the
// [ParseString] function in the same package instead of this method.
func (ByteSize) ParseString(s string) (ByteSize, error) {
	return ParseString(s)
}
