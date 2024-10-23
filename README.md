
# data2xls: Create a Excel XLSX file using raw table in CSV, TSV or JSON format

This utility creates a Excel workbook file from raw text tables, like CSV, TSV or JSON.


## Install

```sh
go install github.com/mixcode/data2xlsx@latest
```

## Examples

### Quick example of creating an Excel sheet
Assume we have 3 tables, each in page1.csv, page2.csv, pag3.json. Make an Excel workbook with 3 sheets, each sheet for each table.
```sh
data2xlsx -o out.xlsx page1.csv page2.tsv page3.json
```

### Appending a sheet to an existing Excel workbook

Add a table 'newpage.csv' as a new sheet to an existing Excel workbook. The new sheet can be named with "-sheet-name" optional argument.
```sh
data2xlsx -a out.xlsx -sheet-name="my new page" newpage.csv
```

Or, you may use pipe to directly send data to a new sheet.
```sh
cat newpage.csv | data2xlsx -a out.xlsx -data-type=csv
```


---

