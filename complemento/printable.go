package complemento

import (
	"fmt"
	"log"
	"time"

	"github.com/leekchan/accounting"
)

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

func PrintPagos(pagos []PrintablePagos) {

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
}
