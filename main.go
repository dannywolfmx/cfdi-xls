package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

const DIR_NAME = "./cfdis2"

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
	f.SetCell("Fecha").SetCell("Nombre").SetCell("Total").SetCell("Folio").SetCell("Version")

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if file.IsDir() {
			return
		}

		path := fmt.Sprintf("%s/%s", DIR_NAME, file.Name())
		content, err := os.ReadFile(path)

		if err != nil {
			log.Fatal(err)
		}

		var cfdi CFDI

		xml.Unmarshal(content, &cfdi)

		f.SetCell(cfdi.Fecha).
			SetCell(cfdi.Receptor.Nombre).
			SetCell(cfdi.Total).
			SetCell(cfdi.Folio).
			SetCell(cfdi.Version)

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
	actualSheet, actualAxis string
	Err                     error

	*excelize.File
}

func NewFile(path string) *BookFile {
	sheet := "Sheet1"
	file := excelize.NewFile()
	file.Path = path
	index := file.NewSheet(sheet)
	file.SetActiveSheet(index)

	return &BookFile{
		File:        file,
		actualSheet: sheet,
		actualAxis:  "A",
	}
}

func (b *BookFile) SetCell(value string) *BookFile {
	b.Err = b.SetCellValue(b.actualSheet, b.actualAxis, value)
	b.nextAxis()

	return b
}

func (b *BookFile) nextAxis() {
	b.actualAxis = "B"
}

//ACell := fmt.Sprintf("A%d", index+2)
//BCell := fmt.Sprintf("B%d", index+2)
//CCell := fmt.Sprintf("C%d", index+2)
//DCell := fmt.Sprintf("D%d", index+2)
//ECell := fmt.Sprintf("E%d", index+2)
