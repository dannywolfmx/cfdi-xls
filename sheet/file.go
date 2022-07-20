package sheet

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type SheetFile struct {
	actualSheet  string
	actualColumn byte
	actualRow    uint
	Err          error

	*excelize.File
}

func NewFile(path string) *SheetFile {
	sheet := "Sheet1"
	file := excelize.NewFile()
	file.Path = path
	index := file.NewSheet(sheet)
	file.SetActiveSheet(index)

	return &SheetFile{
		File:         file,
		actualSheet:  sheet,
		actualColumn: 'A',
		actualRow:    1,
	}
}

func (b *SheetFile) SetCellRight(value string) *SheetFile {
	//Check if previes cells has an error
	if b.Err != nil {
		return b
	}
	axis := fmt.Sprintf("%c%d", b.actualColumn, b.actualRow)
	b.Err = b.SetCellValue(b.actualSheet, axis, value)
	b.NextColumn()

	return b
}

func (b *SheetFile) NextColumn() {
	b.actualColumn++
}

func (b *SheetFile) MoveRowDown() {
	b.actualRow++
}

func (b *SheetFile) MoveRowUpAndResetColumn() {
	b.actualColumn = 'A'
	b.MoveRowUp()
}

func (b *SheetFile) MoveRowDownAndResetColumn() {
	b.actualColumn = 'A'
	b.MoveRowDown()
}

func (b *SheetFile) MoveRowUp() {
	if b.actualRow == 1 {
		return
	}
	b.actualRow--
}
