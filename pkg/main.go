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
		stateful.Include("Common"),
		{Name: `Quantity`, Pattern: `\d+(?:\.\d+)?(?i)[a-z]*`, Action: nil},
		{Name: `WodTitle`, Pattern: `#[^#\n]+`, Action: nil},                  // Workout Title
		{Name: `titleStart`, Pattern: `##+`, Action: stateful.Push("CTitle")}, // Circuit Title Start
		{Name: `block`, Pattern: "```\\w*", Action: nil},                      // MD Block (elided)
		{Name: "whitespace", Pattern: `\s+`, Action: nil},                     // orphan whitespaces (elided)
		{Name: `GString`, Pattern: `\S[^\n]+`, Action: nil},                   // Greedy String
	},
	"CTitle": {
		stateful.Include("Common"),
		{Name: "TitleEnd", Pattern: `\r?\n`, Action: stateful.Pop()},                //
		{Name: `MetaDiv`, Pattern: `(?:-+\s|;)`, Action: stateful.Push("Metadata")}, //
		{Name: `Puntc`, Pattern: `(-|<)`, Action: nil},                              //
		{Name: `TGString`, Pattern: `[^-;<\n]+`, Action: nil},                       // Title Greedy String
		stateful.Return(),
	},
	"Metadata": {
		stateful.Include("Common"),
		{Name: `WodType`, Pattern: `(?i)(AMRAP|EMOM|tabata|superset|dropset)`, Action: nil},
		{Name: `Quantity`, Pattern: `\d+(?:\.\d+)?(?i)[a-z]*`, Action: nil},
		{Name: `Colon`, Pattern: `:`, Action: nil},
		{Name: `Ident`, Pattern: `(?i)[a-z][\w-]+`, Action: nil},
		{Name: "titleEnd", Pattern: `\r?\n`, Action: stateful.Pop()},
		{Name: "space", Pattern: `[\t ]+`, Action: nil},
		stateful.Return(),
	},
	"Common": {
		{Name: "commentStart", Pattern: `<!--+`, Action: stateful.Push("Comment")},
	},
	"Comment": {
		{Name: "commentFinish", Pattern: `--+>`, Action: stateful.Pop()},
		{Name: "dash", Pattern: `-`, Action: nil},
		{Name: "comment", Pattern: `[^-]+`, Action: nil},
	},
})

type Workout struct {
	Identifier string     `parser:"@WodTitle?"`
	Sets       []*Set     `parser:"( @@"`
	Circuits   []*Circuit `parser:"| @@ )*"`
}

type Circuit struct {
	Title []*Title `parser:"@@+"`
	Sets  []*Set   `parser:"@@*"`
}

type Title struct {
	TitleFragments []*Fragment `parser:"@@*"`
	Metadata       []*Metadata `parser:"(MetaDiv @@+)? TitleEnd"`
}

type Fragment struct {
	String string `parser:"( @TGString"`
	Puntc  string `parser:"| @Puntc)"`
}

type Metadata struct {
	WodType  string    `parser:"( @WodType"`
	Quantity *Quantity `parser:"| @Quantity"`
	Words    []*Word   `parser:"| @@ )"`
}

type Word struct {
	String string `parser:"@Ident"`
	Tag    *Tag   `parser:"(Colon @@)?"`
}

type Tag struct {
	Quantity *Quantity `parser:"( @Quantity"`
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
