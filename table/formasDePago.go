package table

import "fmt"

var formaDePago = map[string]string{
	"01": "Efectivo",
	"02": "Cheque nominativo",
	"03": "Transferencia electrónica de fondos",
	"04": "Tarjeta de crédito",
	"05": "Monedero electrónico",
	"06": "Dinero electrónico",
	"08": "Vales de despensa",
	"12": "Dación en pago",
	"13": "Pago por subrogación",
	"14": "Pago por consignación",
	"15": "Condonación",
	"17": "Compensación",
	"23": "Novación",
	"24": "Confusión",
	"25": "Remisión de deuda",
	"26": "Prescripción o caducidad",
	"27": "A satisfacción del acreedor",
	"28": "Tarjeta de débito",
	"29": "Tarjeta de servicios",
	"30": "Aplicación de anticipos",
	"31": "Intermediario pagos",
	"99": "Por definir",
}

// FormaDePago returns the forma de pago
func FormaDePago(key string) string {
	//check if the key exists
	_, ok := formaDePago[key]
	if !ok {
		return fmt.Sprintf(" (%s) No existe la forma de pago", key)
	}

	return fmt.Sprintf(" (%s) %s", key, formaDePago[key])
}
