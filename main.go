package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tealeg/xlsx/v3"
)

// command line options
var (
	flagDataType  string // csv, tsv, json
	flagSheetName string // Sheet names

	flagFirstRowAsHeader bool //

	flagOutputXlsx string // filename for newly created xlsx
	flagOverwrite  bool   // overwrite if file exists

	flagAppendXlsx string
)

// data converted from command line options
var (
	dataTypes  []fileType // source file types converted from flagDataTypes
	sheetNames []string   // sheet names converted from flagSheetName
)

func getSheetName(sheetId, sheetNo int) string { // return a table sheet name
	if sheetId < len(sheetNames) {
		return sheetNames[sheetId]
	}
	return fmt.Sprintf("Table %d", sheetNo)
}

// actual main routine
func run() (err error) {

	// XLSX book
	var wb *xlsx.File
	var outFile string // xlsx filename to save

	if flagAppendXlsx != "" { // "append sheets to existing xlsx"
		// open an existing xlsx file
		wb, err = xlsx.OpenFile(flagAppendXlsx)
		if err != nil {
			if !flagOverwrite {
				return fmt.Errorf("cannot open XLSX file: %w", err)
			}
			wb = xlsx.NewFile()
		}
		outFile = flagAppendXlsx
		flagOverwrite = true
	} else if flagOutputXlsx != "" { // "create a new xlsx"
		// create a new book
		wb = xlsx.NewFile()
		outFile = flagOutputXlsx
	} else {
		return fmt.Errorf("no output file name")
	}

	// convert data type string like "csv,tsv,json" into an array of type codes
	if flagDataType != "" {
		tblType, _ := importCSV(strings.NewReader(flagDataType))
		if len(tblType) > 0 {
			dataTypes = make([]fileType, len(tblType[0]))
			for i, s := range tblType[0] {
				dataTypes[i] = getFileType("TMP." + s)
			}
		}
	}

	// convert sheet name strings in CSV to array of strings
	if flagSheetName != "" {
		tN, _ := importCSV(strings.NewReader(flagSheetName))
		if len(tN) > 0 {
			sheetNames = tN[0]
		}
	}

	sheetNo := 1 + len(wb.Sheets) // start page number of the newly added sheets
	inFile := flag.Args()

	if len(inFile) == 0 { // no file name given: use STDIN for data stream

		if len(dataTypes) == 0 {
			// neither file name and file type for STDIN are given
			return fmt.Errorf("data filename, or file type for stdin must be given")
		}

		fType := dataTypes[0]
		if fType == FILETYPE_UNKNOWN {
			return fmt.Errorf("unknown file type %s", flagDataType)
		}

		// read a table from STDIN, add it to xlsx
		count, e_ := addSheet(wb, 0, sheetNo, fType, os.Stdin)
		if e_ != nil {
			return fmt.Errorf("sheet creation failed: %w", e_)
		}
		sheetNo += count

	} else {

		// load sheets from specified filenames

		sheetId := 0
		for i, srcName := range inFile {

			e_ := func() (err error) {
				// open the source file
				fi, e := os.Open(srcName)
				if e != nil {
					return fmt.Errorf("cannot open data file %s", srcName)
				}
				defer func() {
					e := fi.Close()
					if e != nil && err == nil {
						err = fmt.Errorf("file close failed %s", srcName)
					}
				}()

				// determine the file type
				var fT fileType = FILETYPE_UNKNOWN
				if len(dataTypes) > 0 {
					// use the specified type
					if i < len(dataTypes) {
						fT = dataTypes[i]
					} else {
						// if the more filenames are given than the type list
						// then use the last type for leftovers
						fT = dataTypes[len(dataTypes)-1]
					}
				} else {
					// detect file type by its file extension
					fT = getFileType(srcName)
				}
				if fT == FILETYPE_UNKNOWN {
					return fmt.Errorf("cannot detect file type for %s", srcName)
				}

				// process sheets
				count, e := addSheet(wb, sheetId, sheetNo, fT, fi)
				if e != nil {
					return fmt.Errorf("sheet creation failed: %w", e)
				}
				sheetId += count
				sheetNo += count

				return nil
			}()

			if e_ != nil {
				return e_
			}
		}
	}

	// write XLS
	var fo *os.File
	if _, e_ := os.Stat(outFile); !os.IsNotExist(e_) {
		// file exists
		if !flagOverwrite {
			return fmt.Errorf("output file exists: %s", outFile)
		}
	}
	fo, e_ := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e_ != nil {
		return fmt.Errorf("cannot create output file: %s", outFile)
	}
	defer func() {
		e := fo.Close()
		if e != nil && err == nil {
			err = fmt.Errorf("output file not closed properly: %s", outFile)
		}
	}()
	e_ = wb.Write(fo)
	if e_ != nil {
		return fmt.Errorf("cannot write output file: %s", outFile)
	}

	return nil
}

// program entry point
func main() {
	var err error

	//flag.StringVar(&flagOutputXlsx, "output", "", "New Excel file name to create")
	flag.StringVar(&flagOutputXlsx, "o", "", "New Excel file name to create")

	flag.BoolVar(&flagOverwrite, "overwrite", false, "Overwrite existing Excel file if exists.\nIf '-a' option is used, create a new file if not exists.")
	flag.BoolVar(&flagOverwrite, "y", false, "Same with -overwrite")

	flag.BoolVar(&flagFirstRowAsHeader, "header-row", false, "Treat the first row as the column name header")
	flag.BoolVar(&flagFirstRowAsHeader, "h", false, "Same with -header-row")

	//flag.StringVar(&flagAppendXlsx, "append", "", "Existing Excel file name to append new sheets")
	flag.StringVar(&flagAppendXlsx, "a", "", "Existing Excel file to add sheets")

	flag.StringVar(&flagDataType, "data-type", "", "Input data type. One of csv, tsv, json\nIf omitted, type is guessed from the file extension")
	flag.StringVar(&flagSheetName, "sheet-name", "", "Name of the new sheets added\nUse comma to specify multiple sheet names")

	flag.Usage = func() {
		o := flag.CommandLine.Output()
		cmd := os.Args[0]
		fmt.Fprintf(o, "%s: Create MS Excel books from text data tables in CSV, TSV, JSON\n\n", cmd)
		fmt.Fprintf(o, "Usage:\n%s [options] -o NewExcel.xlsx [infile...]\n", cmd)
		fmt.Fprintf(o, "%s [options] -a Existing.xlsx [infile...]\n", cmd)
		fmt.Fprintln(o)
		flag.PrintDefaults()
		fmt.Fprintln(o)
		fmt.Fprintf(o, `Exmples:
   # Create 'new.xlsx' with 3 sheets.
   # The '-y' option makes the xlsx file to be overwritten if exists.
   $ %[1]s -o new.xlsx -y page1.csv page2.tsv page3.json

   # Add two sheets to 'existing.xlsx' file.
   # The newly add sheet's name will be '1st_page' and '2nd_page'.
   $ %[1]s -a existing.xlsx -sheet-name="1st_page,2nd_page" p1.csv p2.tsv

   # Add a new sheet from STDIN to an exisiting xlsx file.
   $ cat page1.csv | %[1]s -data-type=csv -a existing.xlsx

`, cmd)
	}

	flag.Parse()

	err = run()

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
