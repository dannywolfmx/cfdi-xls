package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

const DIR_NAME = "./cfdis"

type CFDI struct {
	XMLName  xml.Name `xml:"Comprobante"`
	Version  string   `xml:"Version,attr"`
	Receptor Receptor `xml:"Emisor"`
	Folio    string   `xml:"Folio,attr"`
	Fecha    string   `xml:"Fecha,attr"`
	Total    string   `xml:"Total,attr"`
}

type Receptor struct {
	RFC    string `xml:"Rfc,attr"`
	Nombre string `xml:"Nombre,attr"`
}

func main() {
	if !directoryExist(DIR_NAME) {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal("Error al detectar el path del ejecutable")
		}
		path := filepath.Dir(ex)
		log.Fatalf("No existe el directorio %s en la ruta actual %s", DIR_NAME, path)
	}

	files, err := os.ReadDir(DIR_NAME)

	if err != nil {
		log.Fatal(err)
	}

	f := NewFile("Book1.xlsx")
	// Create a new sheet.
	f.SetCellRight("Fecha").SetCellRight("Nombre").SetCellRight("Total").SetCellRight("Folio").SetCellRight("Version")

	for _, file := range files {
		//Read the file
		fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			return
		}

		pathFile := path.Join(DIR_NAME, file.Name())

		content, err := os.ReadFile(pathFile)

		if err != nil {
			log.Fatal(err)
		}

		var cfdi CFDI

		xml.Unmarshal(content, &cfdi)

		//move to the next row down

		f.MoveRowDownAndResetColumn()

		f.SetCellRight(cfdi.Fecha).
			SetCellRight(cfdi.Receptor.Nombre).
			SetCellRight(cfdi.Total).
			SetCellRight(cfdi.Folio).
			SetCellRight(cfdi.Version)

	}
	// Save spreadsheet by the given path.
	if err := f.Save(); err != nil {
		fmt.Println(err)
	}
}

func directoryExist(name string) bool {
	_, err := os.Stat(name)

	return !os.IsNotExist(err)
}

type BookFile struct {
	actualSheet  string
	actualColumn byte
	actualRow    uint
	Err          error

	*excelize.File
}

func NewFile(path string) *BookFile {
	sheet := "Sheet1"
	file := excelize.NewFile()
	file.Path = path
	index := file.NewSheet(sheet)
	file.SetActiveSheet(index)

	return &BookFile{
		File:         file,
		actualSheet:  sheet,
		actualColumn: 'A',
		actualRow:    1,
	}
}

func (b *BookFile) SetCellRight(value string) *BookFile {
	//Check if previes cells has an error
	if b.Err != nil {
		return b
	}
	axis := fmt.Sprintf("%c%d", b.actualColumn, b.actualRow)
	b.Err = b.SetCellValue(b.actualSheet, axis, value)
	b.NextColumn()

	return b
}

func (b *BookFile) NextColumn() {
	b.actualColumn++
}

func (b *BookFile) MoveRowDown() {
	b.actualRow++
}

func (b *BookFile) MoveRowUpAndResetColumn() {
	b.actualColumn = 'A'
	b.MoveRowUp()
}

func (b *BookFile) MoveRowDownAndResetColumn() {
	b.actualColumn = 'A'
	b.MoveRowDown()
}

func (b *BookFile) MoveRowUp() {
	if b.actualRow == 1 {
		return
	}
	b.actualRow--
}

//ACell := fmt.Sprintf("A%d", index+2)
//BCell := fmt.Sprintf("B%d", index+2)
//CCell := fmt.Sprintf("C%d", index+2)
//DCell := fmt.Sprintf("D%d", index+2)
//ECell := fmt.Sprintf("E%d", index+2)
