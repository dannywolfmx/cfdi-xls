package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
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
			FechaTimbrado: complementoDePago.Fecha,
		}

		pago20 := complementoDePago.Complemento.Pagos20

		for _, documento := range pago20.Pago.DoctoRelacionado {
			printablePago := PrintablePago{
				FechaPago:     pago20.Pago.FechaPago,
				ImportePagado: documento.ImpPagado,
			}
			printablePagos.Pagos = append(printablePagos.Pagos, printablePago)
		}

		//Print the data
		ImprimirPantallaPagos(printablePagos)
	}
}

func directoryExist(name string) bool {
	_, err := os.Stat(name)

	return !os.IsNotExist(err)
}

type PrintablePagos struct {
	Emisor        string
	Receptor      string
	FechaTimbrado string
	Pagos         []PrintablePago
}

type PrintablePago struct {
	FechaPago     string
	ImportePagado string
}

func ImprimirPantallaPagos(pago PrintablePagos) {
	fmt.Println("Emisor: ", pago.Emisor)
	fmt.Println("Receptor: ", pago.Receptor)
	fmt.Println("Fecha de timbrado: ", pago.FechaTimbrado)
	for _, p := range pago.Pagos {
		fmt.Println("	Fecha de pago: ", p.FechaPago)
		fmt.Println("	Importe pagado: ", p.ImportePagado)
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
	ImpPagado        string      `xml:"ImpPagado,attr"`
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
