package table

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dannywolfmx/cfdi-xls/complemento"
)

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
		filter:     filterListView(0),
		resumen:    calcularResumen(cfdi),
		Tabs:       filterTabsTitles,
		activeTab:  0,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
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

				selectedFilter := filterTabsContent[m.activeTab][selectedIndex]
				if _, ok := activeFilters[selectedFilter]; ok {
					delete(activeFilters, selectedFilter)
				} else {
					activeFilters[selectedFilter] = listFilters[selectedFilter]
				}

				m.cfdis = filterCFDIS(originalCFDIS)
				m.table = generateCFDITable(m.table.Columns(), transformCFDIToRow(m.cfdis))

				m.filter = filterListView(m.activeTab)
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
		case "left", "h":
			if m.focusState == focusFilter {
				//Move tab to the left
				if m.activeTab > 0 {
					m.activeTab--
					m.filter = filterListView(m.activeTab)
					m.filter.Select(0)
					m.cur = 0
				}
			}

		case "right", "l":
			if m.focusState == focusFilter {
				//Move tab to the right
				if m.activeTab < len(m.Tabs)-1 {
					m.activeTab++
					m.filter = filterListView(m.activeTab)
					m.filter.Select(0)
					m.cur = 0
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
