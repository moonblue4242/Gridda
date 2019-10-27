package actions

// Expandable describes which window gets expanded and how
type Expandable struct {
	Where Where
	How   How
}

// Which describes window positioned inside the grid
type Where struct {
	Column int
	Row    int
	Span   *Span
}

// How defines how a window gets expanded
type How struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}
