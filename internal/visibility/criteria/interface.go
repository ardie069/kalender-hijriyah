package criteria

type Category string

const (
	CategoryA Category = "A" // Easily visible
	CategoryB Category = "B" // Visible under perfect conditions
	CategoryC Category = "C" // May need optical aid
	CategoryD Category = "D" // Visible with optical aid only
	CategoryE Category = "E" // Not visible even with optical aid
	CategoryF Category = "F" // Below Danjon limit / Impossible
)

type Context struct {
	Altitude   float64 // h
	Elongation float64 // L
	ArcOfLight float64 // arcV
	Width      float64 // w
}

type Evaluator interface {
	Evaluate(ctx Context) Category
	Name() string
}