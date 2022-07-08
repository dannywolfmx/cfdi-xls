package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"

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
	files, err := os.ReadDir(DIR_NAME)

	if err != nil {
		log.Fatal(err)
	}

	f := excelize.NewFile()
	// Create a new sheet.
	index := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "Fecha")
	f.SetCellValue("Sheet1", "B1", "Nombre")
	f.SetCellValue("Sheet1", "C1", "Total")
	f.SetCellValue("Sheet1", "D1", "Folio")
	f.SetCellValue("Sheet1", "D1", "Version")

	for index, file := range files {
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

		fmt.Println(cfdi)
		ACell := fmt.Sprintf("A%d", index+2)
		BCell := fmt.Sprintf("B%d", index+2)
		CCell := fmt.Sprintf("C%d", index+2)
		DCell := fmt.Sprintf("D%d", index+2)
		ECell := fmt.Sprintf("E%d", index+2)

		fmt.Println(cfdi.Receptor.Nombre)

		f.SetCellValue("Sheet1", ACell, cfdi.Fecha)
		f.SetCellValue("Sheet1", BCell, cfdi.Receptor.Nombre)
		f.SetCellValue("Sheet1", CCell, cfdi.Total)
		f.SetCellValue("Sheet1", DCell, cfdi.Folio)
		f.SetCellValue("Sheet1", ECell, cfdi.Version)

	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}
}
