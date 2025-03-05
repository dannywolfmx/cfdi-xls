package table

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/dannywolfmx/cfdi-xls/complemento"
)

var (
	filterPUE = cfdiFilterOption{ID: "PUE", Text: "Pago en una sola exhibición"}
	filterPPD = cfdiFilterOption{ID: "PPD", Text: "Pago en parcialidades o diferido"}

	filterFormaPagoEfectivo       = cfdiFilterOption{ID: "01", Text: "Efectivo"}
	filterFormaPagoCheque         = cfdiFilterOption{ID: "02", Text: "Cheque nominativo"}
	filterFormaPagoTransferencia  = cfdiFilterOption{ID: "03", Text: "Transferencia electrónica de fondos"}
	filterFormaPagoTarjetaCredito = cfdiFilterOption{ID: "04", Text: "Tarjeta de crédito"}
	filterFormaPagoMonederoElect  = cfdiFilterOption{ID: "05", Text: "Monedero electrónico"}
	filterFormaDineroElectronico  = cfdiFilterOption{ID: "06", Text: "Dinero electrónico"}
	filterFormaPagoCondonacion    = cfdiFilterOption{ID: "15", Text: "Condonación"}
	filterFormaPagoTarjetaDebito  = cfdiFilterOption{ID: "28", Text: "Tarjeta de débito"}
	filterFormaPagoAnticipo       = cfdiFilterOption{ID: "30", Text: "Aplicación de anticipos"}
	filterFormaPorDefinir         = cfdiFilterOption{ID: "99", Text: "Por definir"}
	filterUsoCFDIG01              = cfdiFilterOption{ID: "G01", Text: "Adquisición de mercancias"}
	filterUsoCFDIG02              = cfdiFilterOption{ID: "G02", Text: "Devoluciones, descuentos o bonificaciones"}
	filterUsoCFDIG03              = cfdiFilterOption{ID: "G03", Text: "Gastos en general"}
	filterIgnoreFilter            = cfdiFilterOption{ID: "IGNORE", Text: "Ignorar filtros"}

	//Tipo de comprobante
	filterTipoIngreso  = cfdiFilterOption{ID: "I", Text: "Ingreso"}
	filterTipoEgreso   = cfdiFilterOption{ID: "E", Text: "Egreso"}
	filterTipoTraslado = cfdiFilterOption{ID: "T", Text: "Traslado"}
	filterTipoPago     = cfdiFilterOption{ID: "P", Text: "Pago"}
)

var listFilters = map[string]cfdiFilterOption{
	filterPUE.ID: filterPUE,
	filterPPD.ID: filterPPD,

	filterFormaPagoEfectivo.ID:       filterFormaPagoEfectivo,
	filterFormaPagoCheque.ID:         filterFormaPagoCheque,
	filterFormaPagoTransferencia.ID:  filterFormaPagoTransferencia,
	filterFormaPagoTarjetaCredito.ID: filterFormaPagoTarjetaCredito,
	filterFormaPagoMonederoElect.ID:  filterFormaPagoMonederoElect,
	filterFormaDineroElectronico.ID:  filterFormaDineroElectronico,
	filterFormaPagoCondonacion.ID:    filterFormaPagoCondonacion,
	filterFormaPagoTarjetaDebito.ID:  filterFormaPagoTarjetaDebito,
	filterFormaPagoAnticipo.ID:       filterFormaPagoAnticipo,
	filterFormaPorDefinir.ID:         filterFormaPorDefinir,
	filterUsoCFDIG01.ID:              filterUsoCFDIG01,
	filterUsoCFDIG02.ID:              filterUsoCFDIG02,
	filterUsoCFDIG03.ID:              filterUsoCFDIG03,

	filterTipoIngreso.ID:  filterTipoIngreso,
	filterTipoEgreso.ID:   filterTipoEgreso,
	filterTipoTraslado.ID: filterTipoTraslado,
	filterTipoPago.ID:     filterTipoPago,

	filterIgnoreFilter.ID: filterIgnoreFilter,
}

var listFilterMetodoDePago = []string{
	filterPUE.ID,
	filterPPD.ID,
}

var listFilterFormaPago = []string{
	filterFormaPagoEfectivo.ID,
	filterFormaPagoCheque.ID,
	filterFormaPagoTransferencia.ID,
	filterFormaPagoTarjetaCredito.ID,
	filterFormaPagoMonederoElect.ID,
	filterFormaDineroElectronico.ID,
	filterFormaPagoCondonacion.ID,
	filterFormaPagoTarjetaDebito.ID,
	filterFormaPagoAnticipo.ID,
	filterFormaPorDefinir.ID,
}

var listFilterTipoComprobante = []string{
	filterTipoIngreso.ID,
	filterTipoEgreso.ID,
	filterTipoTraslado.ID,
	filterTipoPago.ID,
}

var listFilterUsoCFDI = []string{
	filterUsoCFDIG01.ID,
	filterUsoCFDIG02.ID,
	filterUsoCFDIG03.ID,
}

var filterTabsTitles = []string{
	"Metodo de pago",      //PUE, PPD
	"Forma de pago",       //Efectivo, Transferencia, Tarjeta de credito, etc
	"Uso CFDI",            //G01, G02, G03
	"Tipo de comprobante", //I, E, T, P
}

var filterTabsContent = [][]string{
	listFilterMetodoDePago,
	listFilterFormaPago,
	listFilterUsoCFDI,
	listFilterTipoComprobante,
}

var activeFilters = map[string]cfdiFilterOption{}

var originalCFDIS []complemento.CFDI = make([]complemento.CFDI, 0)

func calcularResumen(cfdis []complemento.CFDI) resumen {
	var r resumen

	r.CantidadFacturas = len(cfdis)

	for _, c := range cfdis {
		if c.TipoCambio == "" || c.TipoCambio == "1" {
			r.Descuento += c.Descuento
			r.SubTotal += c.SubTotal
			r.Total += c.Total
		} else {
			fTipoDeCambio, err := strconv.ParseFloat(c.TipoCambio, 64)
			if err != nil {
				fmt.Printf("%v\n", c)
				log.Fatal(err)
			}

			r.Descuento += c.Descuento * fTipoDeCambio
			r.SubTotal += c.SubTotal * fTipoDeCambio
			r.Total += c.Total * fTipoDeCambio
		}
	}

	return r
}

// filterCFDIS es ahora un alias para GenericFilterCFDIS que se encuentra en filters.go
// Mantenemos esta función para compatibilidad con el código existente
func filterCFDIS(cfdis []complemento.CFDI) []complemento.CFDI {
	return GenericFilterCFDIS(cfdis)
}

func transformCFDIToRow(cfdi []complemento.CFDI) []table.Row {
	rows := make([]table.Row, 0)

	for _, c := range cfdi {
		row := table.Row{
			c.Emisor.Nombre,
			c.Receptor.Nombre,
			c.Complemento.TimbreFiscalDigital.FechaTimbrado,
			ac.FormatMoney(c.Total),
		}
		rows = append(rows, row)
	}

	return rows
}
