package complemento

import "encoding/xml"

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
