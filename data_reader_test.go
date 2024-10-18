package main

import (
	"strings"
	//"fmt"
	"testing"

	"encoding/json"
)

func TestJsonImporter(t *testing.T) {
	var err error

	const jsonSrc = `  ["m", 1, {"d":4, "c":3, "b":2, "a":1, "e": true}]  {"i":"i", "j":"j"}`
	const expected = `[["m",1,{"a":1,"b":2,"c":3,"d":4,"e":true}],{"i":"i","j":"j"}]`

	u, err := importJSON(strings.NewReader(jsonSrc))
	if err != nil {
		t.Fatal(err)
	}

	restored, err := json.Marshal(u)
	if string(restored) != expected {
		t.Fatalf("restored data does not match")
	}
}

func TestCSVImporter(t *testing.T) {
	var err error
	const csvSrc = `col1,col2, col3
"row11",row 12, 行1列3
r2c1 , r2c2, r2c3,r2c4`
	const expected = `[["col1","col2"," col3"],["row11","row 12"," 行1列3"],["r2c1 "," r2c2"," r2c3","r2c4"]]`

	u, err := importCSV(strings.NewReader(csvSrc))
	if err != nil {
		t.Fatal(err)
	}

	restored, err := json.Marshal(u)
	if string(restored) != expected {
		t.Fatalf("restored data does not match")
	}
}

func TestTSVImporter(t *testing.T) {
	var err error
	const tsvSrc = `col1	col2	 col3
"row11"	row 12	 行1列3
r2c1 	 r2c2	 r2c3	r2c4`
	const expected = `[["col1","col2"," col3"],["row11","row 12"," 行1列3"],["r2c1 "," r2c2"," r2c3","r2c4"]]`

	u, err := importTSV(strings.NewReader(tsvSrc))
	if err != nil {
		t.Fatal(err)
	}

	restored, err := json.Marshal(u)
	if string(restored) != expected {
		t.Fatalf("restored data does not match")
	}
}
