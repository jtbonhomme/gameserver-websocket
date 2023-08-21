package skyjo

const (
	Zero int = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Eleven
	Twelve
	MinusTwo int = -2
	MinusOne int = -1
)

type Card struct {
	Value   int
	Visible bool
}
