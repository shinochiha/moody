package helpers

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tealeg/xlsx"
)

func CreateXlsx(c Context, header []string, data []map[string]interface{}, filename string) error {
	w := c.Response().Writer
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	column := map[string]int{}

	for i, h := range header {
		cell := AddCell(sheet, 0, i)
		cell.Value = h
		column[h] = i
	}

	for i, x := range data {
		for k, v := range x {
			cell := AddCell(sheet, i+1, column[k])
			t := fmt.Sprintf("%T", v)
			if t == "string" {
				cell.Value = fmt.Sprintf("%v", v)
			} else if t == "float64" {
				cell.Value = fmt.Sprintf("%10.2f", v)
			} else {
				b, ok := v.([]byte)
				if ok {
					uid, err := uuid.ParseBytes(b)
					if err != nil {
						cell.Value = fmt.Sprintf("%v", v)
					} else {
						cell.Value = uid.String()
					}
				} else {
					cell.Value = fmt.Sprintf("%v", v)
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("File-Name", filename)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	return file.Write(w)
}

func AddCell(sheet *xlsx.Sheet, row, col int) *xlsx.Cell {
	for row >= len(sheet.Rows) {
		sheet.AddRow()
	}
	for col >= len(sheet.Rows[row].Cells) {
		sheet.Rows[row].AddCell()
	}
	return sheet.Cell(row, col)
}
