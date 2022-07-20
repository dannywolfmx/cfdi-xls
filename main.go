package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/dannywolfmx/cfdi-xls/sheet"
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

	f := sheet.NewFile("Book1.xlsx")
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
