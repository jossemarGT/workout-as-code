package main

import (
	"os"
	"testing"

	"github.com/alecthomas/repr"
)

func TestParse(t *testing.T) {
	samples := []string{
		"testdata/minimal.txt",
		"testdata/circuits.txt",
		"testdata/wod-annie.txt",
		"testdata/wod-mary.txt",
		"testdata/wod-murph.txt",
		"testdata/v0-parser-full.txt",
	}

	for _, path := range samples {
		t.Run(path, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			wod := &Workout{}
			if err := Parse(wod, f); err != nil {
				t.Log(repr.String(wod, repr.Indent(" "), repr.OmitEmpty(true)))
				t.Log(parser.String())
				t.Fatal(err)
			}
		})
	}
}
