package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	"github.com/alecthomas/repr"
)

// var timeUnits = []string{"s", "m", "h", "sec(s)?", "min(s)", "hr(s)?"}
// var distanceUnits = []string{"mt(s)?", "km(s)?", "mi", "mile(s)?", "ft(s)?"}
// var repetitionUnits = []string{"x"}

var workoutLexer = stateful.MustSimple([]stateful.Rule{
	{Name: `Quantity`, Pattern: `\d+(?:\.\d+)?(?i)[a-z]*`, Action: nil},
	{Name: `CTitle`, Pattern: `##+[^#\n]*`, Action: nil},  // Circuit Title
	{Name: `WTitle`, Pattern: `#[^#\n]*`, Action: nil},    // Workout Title
	{Name: "comment", Pattern: `//[^\n]*`, Action: nil},   //
	{Name: "whitespace", Pattern: `\s+`, Action: nil},     //
	{Name: `GString`, Pattern: `\S+[^/\n]+`, Action: nil}, // Greedy String
})

type Workout struct {
	Identifier string     `parser:"@WTitle?"`
	Sets       []*Set     `parser:"( @@"`
	Circuits   []*Circuit `parser:"| @@ )*"`
}

type Tag struct {
	Pos        lexer.Position
	Identifier string `parser:"@Label"`
	GString    string `parser:"@GString"`
}

type Circuit struct {
	Identifier string `parser:"@CTitle"`
	Sets       []*Set `parser:"@@*"`
}

type Set struct {
	Pos lexer.Position

	Quantity *Quantity `parser:"@Quantity"`
	Exercise *Exercise `parser:"@@"`
}

type Quantity struct {
	Value float64
	Unit  string
}

var notNumber = regexp.MustCompile(`\D+`)

func (q *Quantity) Capture(values []string) (err error) {
	val := values[0]
	var unit string

	match := notNumber.FindStringIndex(val)
	if match != nil {
		unit = val[match[0]:]
		val = val[0:match[0]] //always re-assign last
	}

	q.Unit = unit
	q.Value, err = strconv.ParseFloat(val, 64)

	return err
}

type Exercise struct {
	Pos lexer.Position

	GString string `parser:"@GString"`
}

var parser = participle.MustBuild(&Workout{},
	participle.Lexer(workoutLexer),
)

func Parse(wod *Workout, r io.Reader) error {
	if err := parser.Parse("", r, wod); err != nil {
		return err
	}

	return nil
}

func main() {
	wod := &Workout{}

	if err := Parse(wod, os.Stdin); err != nil {
		fmt.Printf("+%v\n", err)
	}

	repr.Println(wod, repr.Indent("  "), repr.OmitEmpty(true))

	fmt.Println(parser.String())
}
