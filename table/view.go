package table

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dannywolfmx/cfdi-xls/complemento"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// Estilos para la visualización de detalles en rowView y resumenView
var (
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			PaddingLeft(1)

	// Estilo compacto para secciones
	compactSectionStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(0, 1).
				MarginBottom(0) // Sin margen inferior para ahorrar espacio

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Background(lipgloss.Color("236")).
			Bold(true).
			Width(120) // No usar MarginBottom para ahorrar espacio

	moneyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("202")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true)

	// Estilos para layout compacto
	inlineValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				PaddingLeft(1).
				PaddingRight(3)
)

// Estilos adicionales para los filtros
var (
	filterTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Background(lipgloss.Color("236")).
				Bold(true).
				Align(lipgloss.Center).
				PaddingLeft(1).
				PaddingRight(1)

	filterActiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("42")).
				Bold(true)

	filterInactiveStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("250"))
)

// View resumen con layout compacto
func resumenView(r resumen) string {
	var doc strings.Builder

	// Header/Título
	doc.WriteString(headerStyle.Render("Resumen de Facturas"))
	doc.WriteString("\n")

	// Sección de información financiera compacta
	finanzasSection := strings.Builder{}

	// Primera línea: SubTotal y Descuento en la misma línea
	finanzasSection.WriteString(labelStyle.Render("SubTotal:") + moneyStyle.Render(ac.FormatMoney(r.SubTotal)))
	finanzasSection.WriteString("  ")

	if r.Descuento > 0 {
		finanzasSection.WriteString(labelStyle.Render("Descuento:") + warningStyle.Render(ac.FormatMoney(r.Descuento)))
	} else {
		finanzasSection.WriteString(labelStyle.Render("Descuento:") + valueStyle.Render(ac.FormatMoney(r.Descuento)))
	}
	finanzasSection.WriteString("\n")

	// Segunda línea: Total y cantidad de facturas
	finanzasSection.WriteString(labelStyle.Render("Total:") + moneyStyle.Render(ac.FormatMoney(r.Total)))
	finanzasSection.WriteString("  ")
	finanzasSection.WriteString(labelStyle.Render("Facturas:") + infoStyle.Render(strconv.Itoa(r.CantidadFacturas)))

	// Solo mostrar promedio si tenemos facturas
	if r.CantidadFacturas > 0 {
		promedio := r.Total / float64(r.CantidadFacturas)
		finanzasSection.WriteString("  ")
		finanzasSection.WriteString(labelStyle.Render("Promedio:") + valueStyle.Render(ac.FormatMoney(promedio)))
	}

	doc.WriteString(compactSectionStyle.Render(finanzasSection.String()))

	return doc.String()
}

// View for the individual row con layout compacto
func rowView(cfdi complemento.CFDI) string {
	var doc strings.Builder

	// Header/Título
	doc.WriteString(headerStyle.Render("Detalles de la Factura"))
	doc.WriteString("\n")

	// Sección de información general compacta
	generalSection := strings.Builder{}

	// Primera línea: Folio, Serie, UUID
	generalSection.WriteString(labelStyle.Render("Folio:") + inlineValueStyle.Render(cfdi.Folio))
	generalSection.WriteString(labelStyle.Render("Serie:") + inlineValueStyle.Render(cfdi.Serie))
	generalSection.WriteString(labelStyle.Render("UUID:") + valueStyle.Render(cfdi.Complemento.TimbreFiscalDigital.UUID))
	generalSection.WriteString("\n")

	// Segunda línea: Fechas
	generalSection.WriteString(labelStyle.Render("Emisión:") + inlineValueStyle.Render(cfdi.Fecha))
	generalSection.WriteString(labelStyle.Render("Timbrado:") + valueStyle.Render(cfdi.Complemento.TimbreFiscalDigital.FechaTimbrado))

	doc.WriteString(compactSectionStyle.Render(generalSection.String()))
	doc.WriteString("\n")

	// Sección de emisor y receptor compacta
	partesSection := strings.Builder{}
	partesSection.WriteString(labelStyle.Render("Emisor:") + inlineValueStyle.Render(cfdi.Emisor.Nombre) +
		labelStyle.Render("RFC:") + valueStyle.Render(cfdi.Emisor.RFC))
	partesSection.WriteString("\n")
	partesSection.WriteString(labelStyle.Render("Receptor:") + inlineValueStyle.Render(cfdi.Receptor.Nombre) +
		labelStyle.Render("RFC:") + valueStyle.Render(cfdi.Receptor.RFC))

	doc.WriteString(compactSectionStyle.Render(partesSection.String()))
	doc.WriteString("\n")

	// Sección de pago compacta
	pagoSection := strings.Builder{}

	// Primera línea: Método de pago, Forma de pago
	pagoSection.WriteString(labelStyle.Render("Método:") + inlineValueStyle.Render(cfdi.MetodoPago))
	pagoSection.WriteString(labelStyle.Render("Forma:") + valueStyle.Render(FormaDePago(cfdi.FormaPago)))
	pagoSection.WriteString("  ")

	// Segunda línea: Uso CFDI y Moneda (si aplica)
	pagoSection.WriteString(labelStyle.Render("Uso CFDI:") + valueStyle.Render(cfdi.Receptor.UsoCFDI))

	if cfdi.TipoCambio != "" && cfdi.TipoCambio != "1" {
		pagoSection.WriteString("  ")
		pagoSection.WriteString(labelStyle.Render("Moneda:") + valueStyle.Render(cfdi.Moneda+" (TC: "+cfdi.TipoCambio+")"))
	}

	doc.WriteString(compactSectionStyle.Render(pagoSection.String()))
	doc.WriteString("\n")

	// Sección de importes compacta
	importesSection := strings.Builder{}

	// Una sola línea con todos los importes
	importesSection.WriteString(labelStyle.Render("SubTotal:") + inlineValueStyle.Render(ac.FormatMoney(cfdi.SubTotal)))

	if cfdi.Descuento > 0 {
		importesSection.WriteString(labelStyle.Render("Descuento:") + inlineValueStyle.Render(warningStyle.Render(ac.FormatMoney(cfdi.Descuento))))
	}

	importesSection.WriteString(labelStyle.Render("Total:") + moneyStyle.Render(ac.FormatMoney(cfdi.Total)))

	doc.WriteString(compactSectionStyle.Render(importesSection.String()))

	return doc.String()
}

func generateCFDITable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10), // Altura reducida para ocupar menos espacio
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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)

	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

func (m model) ViewFilter(lists list.Model) string {
	//Create tabs for the filters
	doc := strings.Builder{}
	var renderTabs []string

	for i, tab := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}

		border, _, _, _, _ := style.GetBorder()

		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}

		style = style.Border(border, true)
		renderTabs = append(renderTabs, style.Render(tab))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderTabs...)
	doc.WriteString(row)
	doc.WriteString("\n")
	//The list of filters selected
	doc.WriteString(lists.View())
	return doc.String()
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
					MarginTop(0). // Sin margen superior para ahorrar espacio
					Render(m.textarea),
				baseStyle.Width(128).Render(resumenView(m.resumen)),
			),
			lipgloss.JoinVertical(
				lipgloss.Top,
				baseStyle.Render(
					m.ViewFilter(m.filter),
				),
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
					MarginTop(0). // Sin margen superior para ahorrar espacio
					Render(m.textarea),
				baseStyle.Width(128).Render(resumenView(m.resumen)),
			),
			lipgloss.JoinVertical(
				lipgloss.Top,
				baseStyle.
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("69")).
					Render(
						m.ViewFilter(m.filter),
					),
			),
		)
	}
	return s
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

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(i.text))
}

func filterListView(activeTab int) list.Model {
	items := []list.Item{}

	for _, key := range filterTabsContent[activeTab] {
		f := listFilters[key]
		if _, ok := activeFilters[key]; ok {
			items = append(items, item{text: fmt.Sprintf("%s %s-%s", filterActiveStyle.Render("✓"), f.ID, f.Text)})
		} else {
			items = append(items, item{text: fmt.Sprintf("%s %s-%s", filterInactiveStyle.Render("□"), f.ID, f.Text)})
		}
	}

	const defaultListWidth = 40
	const defaultListHeight = 20 // Altura reducida para la lista de filtros

	l := list.New(items, itemDelegate{}, defaultListWidth, defaultListHeight)
	l.Title = filterTitleStyle.Render("Filtros") // Título más corto
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.InfiniteScrolling = true

	// Personalizar estilos de la lista
	styles := list.DefaultStyles()
	styles.Title = styles.Title.
		BorderStyle(lipgloss.NormalBorder()). // Borde normal para ahorrar espacio
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)

	l.Styles = styles

	return l
}
