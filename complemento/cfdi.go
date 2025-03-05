package complemento

import (
	"encoding/xml"
)

type CFDI struct {
	Complemento       ComplementoCFDI `xml:"Complemento"`
	Descuento         float64         `xml:"Descuento,attr"`
	Emisor            Emisor          `xml:"Emisor"`
	Fecha             string          `xml:"Fecha,attr"`
	Folio             string          `xml:"Folio,attr"`
	FormaPago         string          `xml:"FormaPago,attr"`
	MetodoPago        string          `xml:"MetodoPago,attr"`
	Moneda            string          `xml:"Moneda,attr"`
	Receptor          Receptor        `xml:"Receptor"`
	Serie             string          `xml:"Serie,attr"`
	SubTotal          float64         `xml:"SubTotal,attr"`
	TipoCambio        string          `xml:"TipoCambio,attr"`
	TipoDeComprobante string          `xml:"TipoDeComprobante,attr"`
	Total             float64         `xml:"Total,attr"`
	Version           string          `xml:"Version,attr"`
	XMLName           xml.Name        `xml:"Comprobante"`
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
