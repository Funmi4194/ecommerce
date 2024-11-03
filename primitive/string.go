package primitive

type String string

// String returns the string value
func (s String) String() string {
	return string(s)
}
