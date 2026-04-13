package throughput

// UnitString returns a throughput unit string (e.g. "Bps", "KBps").
func (t Throughput) UnitString() string {
	return t.unit.String()
}

// CloneByValue returns a new cloned throughput.
func (t Throughput) CloneByValue() Throughput {
	return *t.Clone()
}

// ParseString parses a throughput string. It is recommended to use the
// [ParseString] function in the same package instead of this method.
func (Throughput) ParseString(s string) (Throughput, error) {
	return ParseString(s)
}
