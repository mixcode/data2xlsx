package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
)

type fileType int

const (
	FILETYPE_UNKNOWN fileType = 0
	FILETYPE_CSV     fileType = 1
	FILETYPE_TSV     fileType = 2
	FILETYPE_JSON    fileType = 3
)

func getFileType(fileName string) fileType {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".csv":
		return FILETYPE_CSV
	case ".tsv":
		return FILETYPE_TSV
	case ".json":
		return FILETYPE_JSON
	}
	return FILETYPE_UNKNOWN
}


// Read json stream.
// This function reads read multiple json objects in the stream
// then returns all of them in an array.
func importJSON(stream io.Reader) ([]any, error) {
	dec := json.NewDecoder(stream)
	dec.UseNumber()	// decode JSON nuber to json.Number object
	res := make([]any, 0)
	for dec.More() {
		var v any
		err := dec.Decode(&v)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}

	return res, nil
}

// Read CSV stream as multiple rows of string columns.
func importCSV(stream io.Reader) ([][]string, error) {
	dec := csv.NewReader(stream)
	dec.FieldsPerRecord = -1 // the number of entries of each line may vary
	return dec.ReadAll()
}

// Read TSV stream as multiple rows of string columns.
func importTSV(stream io.Reader) ([][]string, error) {
	dec := csv.NewReader(stream)
	dec.FieldsPerRecord = -1 // the number of entries of each line may vary
	dec.Comma = '\t'         // use tab as the separator
	return dec.ReadAll()
}
