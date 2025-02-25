package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/leekchan/accounting"
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

type Emisor struct {
	RFC    string `xml:"Rfc,attr"`
	Nombre string `xml:"Nombre,attr"`
}

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

	pagos := make([]PrintablePagos, 0)

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

		var complementoDePago ComplementoDePago

		err = xml.Unmarshal(content, &complementoDePago)

		if err != nil {
			log.Fatal(err)
		}

		//Transform the data to the struct PrintablePagos
		printablePagos := PrintablePagos{
			Emisor:        complementoDePago.Emisor.Nombre,
			Receptor:      complementoDePago.Receptor.Nombre,
			FechaTimbrado: transformFecha(complementoDePago.Fecha),
		}

		pago20 := complementoDePago.Complemento.Pagos20

		for _, documento := range pago20.Pago.DoctoRelacionado {
			printablePago := PrintablePago{
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

	lastMonth := time.Month(0)
	total := 0.0
	contadorFacturas := 0

	for _, pago := range pagos {
		if len(pagos[0].Pagos) == 0 {
			log.Fatal("No se encontraron pagos")
		}

		if lastMonth != pago.Pagos[0].FechaPago.Month() {
			fmt.Println("Total pagos del mes: ", ac.FormatMoney(total))
			fmt.Println("Total de facturas: ", contadorFacturas)
			fmt.Println("-------------------------------------------------")
			fmt.Println()
			fmt.Println()

			contadorFacturas = 0
			total = 0.0

			lastMonth = pago.Pagos[0].FechaPago.Month()
			fmt.Println("-------------------------------------------------")
			fmt.Println("-------------------------------------------------")
			fmt.Println("Pagos del mes de ", pago.Pagos[0].FechaPago.Month())
			fmt.Println("-------------------------------------------------")
			fmt.Println("-------------------------------------------------")
		}

		ImprimirPantallaPagos(pago)

		//Total pagos del mes
		for _, p := range pago.Pagos {
			total += p.ImportePagado
			contadorFacturas++
		}
	}

	fmt.Println("Total pagos del mes: ", ac.FormatMoney(total))
	fmt.Println("Total de facturas: ", contadorFacturas)
	fmt.Println("-------------------------------------------------")

	//Prevent the console from closing
	fmt.Scanln()
}

func printDate(date time.Time) {
	fmt.Println(date.Format("2006-01-02"))
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

type PrintablePagos struct {
	Emisor        string
	Receptor      string
	FechaTimbrado time.Time
	Pagos         []PrintablePago
}

type PrintablePago struct {
	FechaPago     time.Time
	ImportePagado float64
	Folio         string
}

var ac = accounting.Accounting{Symbol: "$", Precision: 2}

func ImprimirPantallaPagos(pago PrintablePagos) {
	fmt.Println("Emisor: ", pago.Emisor)
	fmt.Println("Receptor: ", pago.Receptor)
	fmt.Println("Fecha de timbrado: ", pago.FechaTimbrado.Format("2006-01-02"))
	for _, p := range pago.Pagos {
		fmt.Println("	Folio: ", p.Folio)
		fmt.Println("	Fecha de pago: ", p.FechaPago.Format("2006-01-02"))

		fmt.Println("	Importe pagado: ", ac.FormatMoney(p.ImportePagado))
		fmt.Println("")
	}
	fmt.Println("-------------------------------------------------")
}

type ComplementoDePago struct {
	XMLName     xml.Name    `xml:"Comprobante"`
	Version     string      `xml:"Version,attr"`
	Emisor      Emisor      `xml:"Emisor"`
	Receptor    Receptor    `xml:"Receptor"`
	Folio       string      `xml:"Folio,attr"`
	Fecha       string      `xml:"Fecha,attr"`
	Total       string      `xml:"Total,attr"`
	Complemento Complemento `xml:"Complemento"`
}

type Complemento struct {
	Pagos20 Pagos20 `xml:"Pagos"`
}

type Pagos20 struct {
	Totales Totales `xml:"Totales"`
	Pago    Pago    `xml:"Pago"`
}

type Totales struct {
	TotalTrasladosBaseIVA16     string `xml:"TotalTrasladosBaseIVA16,attr"`
	TotalTrasladosImpuestoIVA16 string `xml:"TotalTrasladosImpuestoIVA16,attr"`
	MontoTotalPagos             string `xml:"MontoTotalPagos,attr"`
}

type Pago struct {
	XMLName          xml.Name           `xml:"Pago"`
	FechaPago        string             `xml:"FechaPago,attr"`
	Monto            string             `xml:"Monto,attr"`
	FormaDePagoP     string             `xml:"FormaDePagoP,attr"`
	MonedaP          string             `xml:"MonedaP,attr"`
	TipoCambioP      string             `xml:"TipoCambioP,attr"`
	NumOperacion     string             `xml:"NumOperacion,attr"`
	DoctoRelacionado []DoctoRelacionado `xml:"DoctoRelacionado"`
}

type DoctoRelacionado struct {
	IDDocumento      string      `xml:"IdDocumento,attr"`
	Serie            string      `xml:"Serie,attr"`
	Folio            string      `xml:"Folio,attr"`
	ImpSaldoAnt      string      `xml:"ImpSaldoAnt,attr"`
	ImpPagado        float64     `xml:"ImpPagado,attr"`
	ImpSaldoInsoluto string      `xml:"ImpSaldoInsoluto,attr"`
	ImpuestosDR      ImpuestosDR `xml:"ImpuestosDR"`
}

type ImpuestosDR struct {
	TrasladosDR TrasladosDR `xml:"TrasladosDR"`
}

type TrasladosDR struct {
	TrasladoDR TrasladoDR `xml:"TrasladoDR"`
}

type TrasladoDR struct {
	BaseDR       string `xml:"BaseDR,attr"`
	ImpuestoDR   string `xml:"ImpuestoDR,attr"`
	TipoFactorDR string `xml:"TipoFactorDR,attr"`
	TasaOCuotaDR string `xml:"TasaOCuotaDR,attr"`
	ImporteDR    string `xml:"ImporteDR,attr"`
}
