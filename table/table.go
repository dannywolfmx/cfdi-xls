package table

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	filterPUE.ID:                     filterPUE,
	filterPPD.ID:                     filterPPD,
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

var orderedMapKeys = []string{
	filterPUE.ID,
	filterPPD.ID,

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

	filterUsoCFDIG01.ID,
	filterUsoCFDIG02.ID,
	filterUsoCFDIG03.ID,

	filterTipoIngreso.ID,
	filterTipoEgreso.ID,
	filterTipoTraslado.ID,
	filterTipoPago.ID,

	filterIgnoreFilter.ID,
}

var activeFilters = map[string]cfdiFilterOption{}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var originalCFDIS []complemento.CFDI = make([]complemento.CFDI, 0)

func PrintTable(cfdi []complemento.CFDI) {
	originalCFDIS = cfdi

	columns := []table.Column{
		{Title: "Emisor", Width: 40},
		{Title: "Receptor", Width: 40},
		{Title: "Fecha de timbrado", Width: 20},
		{Title: "Importe pagado", Width: 20},
	}

	cfdi = filterCFDIS(cfdi)

	rows := transformCFDIToRow(cfdi)

	t := generateCFDITable(columns, rows)

	textArea := ""
	if len(cfdi) != 0 {
		textArea = rowView(cfdi[0])
	}

	m := model{
		table:      t,
		textarea:   textArea,
		cfdis:      cfdi,
		focusState: focusTable,
		filter:     filterListView(),
		resumen:    calcularResumen(cfdi),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func generateCFDITable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)

	return t
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
}

// Init
func (m model) Init() tea.Cmd {
	return nil
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		//change focus
		case "tab":
			if m.focusState == focusTable {
				m.focusState = focusFilter
			} else {
				m.focusState = focusTable
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			//Add the selected row to the textarea
			if m.table.Focused() {
				selectedRow := m.table.Cursor()

				cfdi := m.cfdis[selectedRow]
				//Open the file with the default program
				cmd := exec.Command("cmd", "/c", "start", fmt.Sprintf("%s.xml", cfdi.Complemento.TimbreFiscalDigital.UUID))
				//Obtener el directorio actual
				ex, err := os.Executable()
				if err != nil {
					log.Fatal("Error al detectar el path del ejecutable")
				}

				path := filepath.Dir(ex)
				cmd.Dir = fmt.Sprintf("%s\\cfdis", path)
				err = cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				return m, nil
			}
		//Filter
		case " ":
			if m.focusState == focusFilter {
				selectedIndex := m.filter.Index()

				selectedFilter := listFilters[orderedMapKeys[selectedIndex]]
				if _, ok := activeFilters[selectedFilter.ID]; ok {
					delete(activeFilters, selectedFilter.ID)
				} else {
					activeFilters[selectedFilter.ID] = selectedFilter
				}

				m.cfdis = filterCFDIS(originalCFDIS)
				m.table = generateCFDITable(m.table.Columns(), transformCFDIToRow(m.cfdis))

				m.filter = filterListView()
				m.filter.Select(selectedIndex)
				//Update resumen
				m.resumen = calcularResumen(m.cfdis)
				m.cur = 0
			}
		case "up", "k":
			if m.focusState == focusTable {
				if m.cur > 0 && m.table.Focused() {
					m.cur--
					m.textarea = rowView(m.cfdis[m.cur])
				}
			}

		case "down", "j":
			if m.focusState == focusTable {
				if m.cur < len(m.cfdis)-1 && m.table.Focused() {
					m.cur++
					m.textarea = rowView(m.cfdis[m.cur])
				}

			}
		}
	}

	if m.focusState == focusTable {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.focusState == focusFilter {
		m.filter, cmd = m.filter.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

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

// View
func (m model) View() string {
	var s string
	if m.focusState == focusTable {
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Top,
				baseStyle.
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("69")).
					Render(m.table.View()),
				baseStyle.
					Width(128).
					Render(m.textarea),
				baseStyle.Width(128).Render(resumenView(m.resumen)),
			),
			lipgloss.JoinVertical(
				lipgloss.Top,
				baseStyle.Render(m.filter.View()),
			),
		)
	} else {
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.JoinVertical(
				lipgloss.Top,
				baseStyle.
					Render(m.table.View()),
				baseStyle.
					Width(128).
					Render(m.textarea),
				baseStyle.Width(128).Render(resumenView(m.resumen)),
			),
			lipgloss.JoinVertical(
				lipgloss.Top,
				baseStyle.
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("69")).
					Render(m.filter.View()),
			),
		)
	}
	return s
}

var (
	focusedModelStyle = lipgloss.NewStyle().
		Width(120).
		Height(10).
		Align(lipgloss.Center, lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))
)

// View for the individual row
func rowView(cfdi complemento.CFDI) string {
	var s string

	s += "Folio: " + cfdi.Folio + "\t"
	s += "Fecha de timbrado: " + cfdi.Fecha + "\t"
	s += "UUID: " + cfdi.Complemento.TimbreFiscalDigital.UUID + "\n"
	s += "Emisor: " + cfdi.Emisor.Nombre + "\n"
	s += "Receptor: " + cfdi.Receptor.Nombre + "\n"
	s += "Forma de pago: " + FormaDePago(cfdi.FormaPago) + "\t"
	s += "Metodo de pago: " + cfdi.MetodoPago + "\t"
	s += "Uso CFDI: " + cfdi.Receptor.UsoCFDI + "\n"
	s += "Importe pagado: " + ac.FormatMoney(cfdi.Total) + "\n"

	return s
}

// View resumen
func resumenView(r resumen) string {
	var s string

	s += "SubTotal: " + ac.FormatMoney(r.SubTotal) + "\n"
	s += "Descuento: " + ac.FormatMoney(r.Descuento) + "\n"
	s += "Total: " + ac.FormatMoney(r.Total) + "\n"
	s += "Cantidad de facturas: " + strconv.Itoa(r.CantidadFacturas) + "\n"

	return s
}

// Currency style
var ac = accounting.Accounting{Symbol: "$", Precision: 2}

// Incluye solo los CFDIS que cumplan con al menos uno de los filtros y evitar agregar duplicados
func filterCFDIS(cfdis []complemento.CFDI) []complemento.CFDI {

	//check if ignore is active
	_, ok := activeFilters[filterIgnoreFilter.ID]

	if len(activeFilters) == 0 || ok {
		return cfdis
	}

	filteredCFDIS := make([]complemento.CFDI, 0)

	//check if the filters PUE or PPD are disabled
	_, pueActive := activeFilters[filterPUE.ID]
	_, ppdActive := activeFilters[filterPPD.ID]

	if !pueActive && !ppdActive {
		filteredCFDIS = make([]complemento.CFDI, len(cfdis))
		copy(filteredCFDIS, cfdis)
	} else {
		for _, c := range cfdis {
			if _, ok := activeFilters[c.MetodoPago]; ok {
				filteredCFDIS = append(filteredCFDIS, c)
				continue
			}
		}
	}

	//check if the filters Forma de pago are disabled
	var filtroActivo bool
	for _, c := range listFilterFormaPago {
		if _, filtroActivo = activeFilters[c]; filtroActivo {
			break
		}
	}

	if !filtroActivo {
		cfdis = make([]complemento.CFDI, len(filteredCFDIS))
		copy(cfdis, filteredCFDIS)
	} else {
		cfdis = make([]complemento.CFDI, 0)
		for _, c := range filteredCFDIS {
			if _, ok := activeFilters[c.FormaPago]; ok {
				cfdis = append(cfdis, c)
				continue
			}
		}
	}

	//check if the filters UsoCFDI are disabled
	filtroActivo = false
	for _, c := range listFilterUsoCFDI {
		if _, filtroActivo = activeFilters[c]; filtroActivo {
			break
		}
	}

	if !filtroActivo {
		filteredCFDIS = make([]complemento.CFDI, len(cfdis))
		copy(filteredCFDIS, cfdis)
	} else {
		filteredCFDIS = make([]complemento.CFDI, 0)
		for _, c := range cfdis {
			if _, ok := activeFilters[c.Receptor.UsoCFDI]; ok {
				filteredCFDIS = append(filteredCFDIS, c)
				continue
			}
		}
	}

	//check if the filters Tipo de comprobante are disabled
	filtroActivo = false
	for _, c := range listFilterTipoComprobante {
		if _, filtroActivo = activeFilters[c]; filtroActivo {
			break
		}
	}

	if !filtroActivo {
		cfdis = make([]complemento.CFDI, len(filteredCFDIS))
		copy(cfdis, filteredCFDIS)
	} else {
		cfdis = make([]complemento.CFDI, 0)
		for _, c := range filteredCFDIS {
			if _, ok := activeFilters[c.TipoDeComprobante]; ok {
				cfdis = append(cfdis, c)
				continue
			}
		}
	}

	return cfdis
}

type item struct {
	text string
}

func (i item) FilterValue() string {
	return i.text
}

func filterListView() list.Model {
	items := []list.Item{}

	for _, key := range orderedMapKeys {
		f := listFilters[key]
		if _, ok := activeFilters[key]; ok {
			items = append(items, item{text: fmt.Sprintf(" (•) %s - %s", f.ID, f.Text)})
		} else {
			items = append(items, item{text: fmt.Sprintf(" ( ) %s - %s", f.ID, f.Text)})
		}
	}

	const defaultListWidth = 40
	const defaultListHeight = 30

	l := list.New(items, itemDelegate{}, defaultListWidth, defaultListHeight)
	l.Title = "Filtros"
	l.InfiniteScrolling = true

	return l
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.text)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)
