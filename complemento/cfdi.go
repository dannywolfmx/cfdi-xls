package complemento

import (
	"encoding/xml"
)

type CFDI struct {
	XMLName           xml.Name        `xml:"Comprobante"`
	Version           string          `xml:"Version,attr"`
	Receptor          Receptor        `xml:"Receptor"`
	Emisor            Emisor          `xml:"Emisor"`
	Folio             string          `xml:"Folio,attr"`
	Fecha             string          `xml:"Fecha,attr"`
	Descuento         float64         `xml:"Descuento,attr"`
	SubTotal          float64         `xml:"SubTotal,attr"`
	Total             float64         `xml:"Total,attr"`
	MetodoPago        string          `xml:"MetodoPago,attr"`
	Complemento       ComplementoCFDI `xml:"Complemento"`
	FormaPago         string          `xml:"FormaPago,attr"`
	TipoCambio        string          `xml:"TipoCambio,attr"`
	TipoDeComprobante string          `xml:"TipoDeComprobante,attr"`
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
