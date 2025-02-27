package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/dannywolfmx/cfdi-xls/complemento"
)

const DIR_NAME = "./cfdis"

func main() {
	//Check if the directory exists
	if !directoryExist(DIR_NAME) {
		ex, err := os.Executable()
		if err != nil {
			log.Fatal("Error al detectar el path del ejecutable")
		}
		path := filepath.Dir(ex)
		log.Fatalf("No existe el directorio %s en la ruta actual %s", DIR_NAME, path)
	}

	//Get the files in the directory
	files, err := os.ReadDir(DIR_NAME)

	if err != nil {
		log.Fatal(err)
	}

	//ComplementoDePagoPrint(files)
	CFDIPrint(files)

	//Prevent the console from closing
	fmt.Scanln()
}

func CFDIPrint(files []os.DirEntry) {
	cfdis := make([]complemento.CFDI, 0)

	for _, file := range files {
		//Read the file
		if file.IsDir() {
			continue
		}
		pathFile := path.Join(DIR_NAME, file.Name())

		//check if the extension is a XML
		if filepath.Ext(pathFile) != ".xml" {
			continue
		}

		content, err := os.ReadFile(pathFile)

		if err != nil {
			log.Fatal(err)
		}

		var cfdi complemento.CFDI

		err = xml.Unmarshal(content, &cfdi)

		if err != nil {
			log.Fatal(err)
		}

		cfdis = append(cfdis, cfdi)
	}

	if len(cfdis) == 0 {
		log.Fatal("No se encontraron facturas")
	}

	//Sort the cfdis by date
	sort.Slice(cfdis, func(i, j int) bool {
		fechaI, err := time.Parse("2006-01-02T15:04:05", cfdis[i].Fecha)
		if err != nil {
			log.Fatal(err)
		}
		fechaJ, err := time.Parse("2006-01-02T15:04:05", cfdis[j].Fecha)
		if err != nil {
			log.Fatal(err)
		}
		return fechaI.Before(fechaJ)
	})

	cfdisPUE := make([]complemento.CFDI, 0)
	//cfdisPPD := make([]complemento.CFDI, 0)

	for _, cfdi := range cfdis {
		//		if cfdi.MetodoPago == "PUE" && cfdi.FormaPago != "01" && cfdi.FormaPago != "15" && cfdi.FormaPago != "30" && (cfdi.Receptor.UsoCFDI == "G01" || cfdi.Receptor.UsoCFDI == "G03") {
		//			cfdisPUE = append(cfdisPUE, cfdi)
		//		} else if cfdi.MetodoPago == "PPD" && (cfdi.Receptor.UsoCFDI == "G01" || cfdi.Receptor.UsoCFDI == "G03") {
		//			cfdisPPD = append(cfdisPPD, cfdi)
		//		}

		//Solo forma de pago 01
		if cfdi.FormaPago == "01" && (cfdi.Receptor.UsoCFDI == "G01" || cfdi.Receptor.UsoCFDI == "G03") {
			cfdisPUE = append(cfdisPUE, cfdi)
		}

	}

	total := 0.0
	subTotal := 0.0
	formasDePago := make(map[string]int)
	for _, cfdi := range cfdisPUE {
		cfdi.Print()
		f, err := strconv.ParseFloat(cfdi.Total, 64)
		if err != nil {
			log.Fatal(err)
		}

		if f == 0 {
			continue
		}

		s, err := strconv.ParseFloat(cfdi.SubTotal, 64)

		if err != nil {
			log.Fatal(err)
		}

		formasDePago[cfdi.FormaPago]++

		if cfdi.TipoCambio != "" {
			tipoCambio, err := strconv.ParseFloat(cfdi.TipoCambio, 64)
			if err != nil {
				log.Fatal(err)
			}

			if tipoCambio != 1 {
				f = f * tipoCambio
				s = s * tipoCambio
			}
		}

		total += f
		subTotal += s
	}

	fmt.Println("Total Subtotal: ", subTotal)
	fmt.Println("Total: ", total)
	fmt.Println("Total de facturas PUE: ", len(cfdisPUE))

	//FormasDePago
	fmt.Println("Formas de pago")
	for k, v := range formasDePago {
		fmt.Println(k, ":", v)
	}

	//count
	fmt.Println("-------------------------------------------------")
	fmt.Println(len(cfdisPUE))

}

func ComplementoDePagoPrint(files []os.DirEntry) {
	pagos := make([]complemento.PrintablePagos, 0)

	for _, file := range files {
		//Read the file
		if file.IsDir() {
			continue
		}
		pathFile := path.Join(DIR_NAME, file.Name())

		//check if the extension is a XML
		if filepath.Ext(pathFile) != ".xml" {
			continue
		}

		content, err := os.ReadFile(pathFile)

		if err != nil {
			log.Fatal(err)
		}

		var complementoDePago complemento.ComplementoDePago

		err = xml.Unmarshal(content, &complementoDePago)

		if err != nil {
			log.Fatal(err)
		}

		//Transform the data to the struct PrintablePagos
		printablePagos := complemento.PrintablePagos{
			Emisor:        complementoDePago.Emisor.Nombre,
			Receptor:      complementoDePago.Receptor.Nombre,
			FechaTimbrado: transformFecha(complementoDePago.Fecha),
		}

		pago20 := complementoDePago.Complemento.Pagos20

		for _, documento := range pago20.Pago.DoctoRelacionado {
			printablePago := complemento.PrintablePago{
				FechaPago:     transformFecha(pago20.Pago.FechaPago),
				ImportePagado: documento.ImpPagado,
				Folio:         documento.Folio,
			}
			printablePagos.Pagos = append(printablePagos.Pagos, printablePago)
		}

		pagos = append(pagos, printablePagos)
	}
	if len(pagos) == 0 {
		log.Fatal("No se encontraron pagos")
	}

	//Sort the pagos by date
	sort.Slice(pagos, func(i, j int) bool {
		return pagos[i].Pagos[0].FechaPago.Before(pagos[j].Pagos[0].FechaPago)
	})

	complemento.PrintPagos(pagos)

}

func transformFecha(fecha string) time.Time {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, fecha)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func directoryExist(name string) bool {
	_, err := os.Stat(name)

	return !os.IsNotExist(err)
}
