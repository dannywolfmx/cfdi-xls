package complemento

import (
	"encoding/xml"
	"fmt"
)

type CFDI struct {
	XMLName     xml.Name        `xml:"Comprobante"`
	Version     string          `xml:"Version,attr"`
	Receptor    Receptor        `xml:"Receptor"`
	Emisor      Emisor          `xml:"Emisor"`
	Folio       string          `xml:"Folio,attr"`
	Fecha       string          `xml:"Fecha,attr"`
	SubTotal    string          `xml:"SubTotal,attr"`
	Total       string          `xml:"Total,attr"`
	MetodoPago  string          `xml:"MetodoPago,attr"`
	Complemento ComplementoCFDI `xml:"Complemento"`
	FormaPago   string          `xml:"FormaPago,attr"`
	TipoCambio  string          `xml:"TipoCambio,attr"`
}

type ComplementoCFDI struct {
	TimbreFiscalDigital TimbreFiscalDigital `xml:"TimbreFiscalDigital"`
}

type TimbreFiscalDigital struct {
	UUID          string `xml:"UUID,attr"`
	FechaTimbrado string `xml:"FechaTimbrado,attr"`
}

type Receptor struct {
	RFC     string `xml:"Rfc,attr"`
	Nombre  string `xml:"Nombre,attr"`
	UsoCFDI string `xml:"UsoCFDI,attr"`
}

type Emisor struct {
	RFC    string `xml:"Rfc,attr"`
	Nombre string `xml:"Nombre,attr"`
}

func (c CFDI) Print() {
	fmt.Println("UUID: ", c.Complemento.TimbreFiscalDigital.UUID)
	fmt.Println("Emisor: ", c.Emisor.Nombre)
	fmt.Println("Emisor RFC: ", c.Emisor.RFC)
	fmt.Println("Receptor: ", c.Receptor.Nombre)
	fmt.Println("Fecha de timbrado: ", c.Fecha)
	fmt.Println("Metodo de pago", c.MetodoPago)
	fmt.Println("Uso CFDI: ", c.Receptor.UsoCFDI)
	fmt.Println("Forma de pago: ", c.FormaPago)
	fmt.Println("Total: ", c.Total)
	fmt.Println("Tipo de cambio: ", c.TipoCambio)
	fmt.Println("-------------------------------------------------")
}
