package main

import (
	"fmt"
	"io"
	"encoding/json"

	"github.com/tealeg/xlsx/v3" // tealeg's xlsx toolkit, v3
)

// TODO
// add column info

func addSheet(wb *xlsx.File, sheetId, sheetNo int, fType fileType, reader io.Reader) (sheetCount int, err error) {

	if fType == FILETYPE_CSV || fType == FILETYPE_TSV { // CSV/TSV raw values

		// load the table
		var table [][]string
		var e_ error
		if fType == FILETYPE_CSV {
			table, e_ = importCSV(reader)
			if e_ != nil {
				err = fmt.Errorf("CSV read fail: %w", e_)
				return
			}
		} else {
			table, e_ = importTSV(reader)
			if e_ != nil {
				err = fmt.Errorf("TSV read fail: %w", e_)
				return
			}
		}

		// add a sheet
		sh, e_ := wb.AddSheet(getSheetName(sheetId, sheetNo))
		if e_ != nil {
			err = fmt.Errorf("cannot add new sheet to xlsx: %w", e_)
			return
		}

		// add title row

		// add data row
		for i, row := range table {
			wRow := sh.AddRow()

			// TODO: row formatting
			for j, s := range row {

				_ = j

				wCel := wRow.AddCell()
				wCel.SetString(s)

				// TODO: column formatting
			}

			if i == 0 && flagFirstRowAsHeader {
				//	TODO: Format the first row as header
			}
		}

		// sheet added
		sheetCount = 1
		return

	} else if fType == FILETYPE_JSON {

		// load the JSON table
		jT, e_ := importJSON(reader)
		if e_ != nil {
			err = fmt.Errorf("JSON read fail: %w", e_)
			return
		}

		for _, table := range jT {

			var rows []any

			switch tbl := table.(type) {
			case map[string]any: // js object
				// column = tbl['column']	// column info
				rows = tbl["data"].([]any) // data rows

			case []any: // js array
				rows = tbl

				// TODO
				//if (flagFirstRowAsHeader) {
				//	header = rows[:1]
				//	rows = rows[1:]
				//}

			default:
				err = fmt.Errorf("table must be an json array")
				return
			}

			// add a sheet
			sh, e_ := wb.AddSheet(getSheetName(sheetId, sheetNo))
			if e_ != nil {
				err = fmt.Errorf("cannot add new sheet to xlsx: %w", e_)
				return
			}

			// TODO: format header

			for _, arow := range rows {

				row, ok := arow.([]any)
				if !ok {
					err = fmt.Errorf("row must be an json array")
					return
				}

				wRow := sh.AddRow()

				// TODO: row formatting
				for j, s := range row {

					wCel := wRow.AddCell()

					//
					// TODO: set data type by col type
					//
					switch v := s.(type) {
					case string:
						wCel.SetString(v)

					case json.Number:
						iVal, e := v.Int64()
						if e==nil {
							wCel.SetInt64(iVal)
							break
						}

						fVal, e := v.Float64()
						if e==nil {
							wCel.SetFloat(fVal)
							break
						}

					case int:
						wCel.SetInt(v)

					case int64:
						wCel.SetInt64(v)
					case float64:
						wCel.SetFloat(v)
					case bool:
						wCel.SetBool(v)

					default:
						//wCel.SetString(fmt.Sprintf("%v", v))
						wCel.SetValue(v)
					}

					// TODO: column formatting
					_ = j
				}

			}

			sheetCount++
			sheetId++
			sheetNo++
		}

		return

	}

	err = fmt.Errorf("unknown source data type")
	return

}
