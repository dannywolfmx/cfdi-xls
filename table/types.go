package table

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/dannywolfmx/cfdi-xls/complemento"
	"github.com/leekchan/accounting"
)

type cfdiFilterOption struct {
	ID   string
	Text string
}

type focusState uint

const (
	focusTable focusState = iota
	focusFilter
)

type resumen struct {
	SubTotal         float64
	Descuento        float64
	Total            float64
	CantidadFacturas int
}

type model struct {
	focusState focusState
	table      table.Model
	filter     list.Model
	textarea   string
	cfdis      []complemento.CFDI
	cur        int
	resumen    resumen
	Tabs       []string
	activeTab  int
}

type item struct {
	text string
}

func (i item) FilterValue() string {
	return i.text
}

// Currency style
var ac = accounting.Accounting{Symbol: "$", Precision: 2}
