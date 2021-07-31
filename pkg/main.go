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

var wodStatefulLexer = stateful.Must(stateful.Rules{
	"Root": {
		{Name: `Quantity`, Pattern: `\d+(?:\.\d+)?(?i)[a-z]*`, Action: nil},
		{Name: `WodTitle`, Pattern: `#[^#\n]+`, Action: nil},              // Workout Title
		{Name: `tStart`, Pattern: `##+`, Action: stateful.Push("CTitle")}, // Circuit Title Start
		{Name: "comment", Pattern: `//[^\n]*`, Action: nil},               // comments (elided)
		{Name: `GString`, Pattern: `[^\n]+`, Action: nil},                 // Greedy String (sink hole)
		{Name: "whitespace", Pattern: `\s+`, Action: nil},                 // orphan whitespaces (elided)
	},
	"CTitle": {
		{Name: "newline", Pattern: `\r?\n`, Action: stateful.Pop()},
		{Name: `WodType`, Pattern: `(?i)(AMRAP|EMOM|tabata|for-time)`, Action: nil},
		{Name: `Quantity`, Pattern: `\d+(?:\.\d+)?(?i)[a-z]*`, Action: nil},
		{Name: `MetaDiv`, Pattern: `(?:\s-+\s|;)`, Action: nil},
		{Name: `Colon`, Pattern: `:`, Action: nil},
		{Name: `Ident`, Pattern: `(?i)[a-z][\w-]+`, Action: nil},
		{Name: `Punct`, Pattern: `[!/@[` + "`" + `{~]`, Action: nil},
		{Name: "space", Pattern: `[\t ]+`, Action: nil},
		stateful.Return(),
	},
})

type Workout struct {
	Identifier string     `parser:"@WodTitle?"`
	Sets       []*Set     `parser:"( @@"`
	Circuits   []*Circuit `parser:"| @@ )*"`
}

type Circuit struct {
	Title Title  `parser:"@@"`
	Sets  []*Set `parser:"@@*"`
}

type Title struct {
	TitleFragments []*Fragment `parser:"@@*"`
	Metadata       []*Metadata `parser:"(MetaDiv @@+)?"`
}

type Fragment struct {
	Ident string `parser:"( @Ident"`
	Punct string `parser:"| @Punct"`
	Colon string `parser:"| @Colon)"`
}

type Metadata struct {
	WodType  string    `parser:"( @WodType"`
	Quantity *Quantity `parser:"| @Quantity"`
	Tags     []*Tag    `parser:"| @@)"`
}

type Tag struct {
	Key      string    `parser:"@Ident Colon"`
	Quantity *Quantity `parser:"(@Quantity"`
	Value    string    `parser:"| @Ident)"`
}

type Set struct {
	Pos lexer.Position

	Quantity *Quantity `parser:"@Quantity"`
	Exercise *Exercise `parser:"@@"`
}

type Exercise struct {
	Pos lexer.Position

	GString string `parser:"@GString"`
}

type Quantity struct {
	Value float64
	Unit  string
	// TODO add IsTime() bool method
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

var parser = participle.MustBuild(&Workout{},
	participle.Lexer(wodStatefulLexer),
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
