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

// View resumen
func resumenView(r resumen) string {
	var s string

	s += "SubTotal: " + ac.FormatMoney(r.SubTotal) + "\n"
	s += "Descuento: " + ac.FormatMoney(r.Descuento) + "\n"
	s += "Total: " + ac.FormatMoney(r.Total) + "\n"
	s += "Cantidad de facturas: " + strconv.Itoa(r.CantidadFacturas) + "\n"

	return s
}

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
