package bitrate

// UnitString returns a bitrate unit string (e.g. "bit/s", "Gbit/s").
func (b Bitrate) UnitString() string {
	return b.unit.String()
}

// CloneByValue returns a new cloned bitrate.
func (b Bitrate) CloneByValue() Bitrate {
	return *b.Clone()
}

// ParseString parses a bitrate string. It is recommended to use the
// [ParseString] function in the same package instead of this method.
func (Bitrate) ParseString(s string) (Bitrate, error) {
	return ParseString(s)
}
