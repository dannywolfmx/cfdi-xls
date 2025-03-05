package table

import (
	"github.com/dannywolfmx/cfdi-xls/complemento"
)

// FilterStrategy define una interfaz para estrategias de filtrado
type FilterStrategy interface {
	Apply(cfdis []complemento.CFDI) []complemento.CFDI
	IsActive() bool
}

// MetodoPagoFilter implementa filtrado por método de pago (PUE, PPD)
type MetodoPagoFilter struct {
	filters map[string]cfdiFilterOption
}

func NewMetodoPagoFilter(activeFilters map[string]cfdiFilterOption) *MetodoPagoFilter {
	return &MetodoPagoFilter{
		filters: activeFilters,
	}
}

func (f *MetodoPagoFilter) IsActive() bool {
	_, pueActive := f.filters[filterPUE.ID]
	_, ppdActive := f.filters[filterPPD.ID]
	return pueActive || ppdActive
}

func (f *MetodoPagoFilter) Apply(cfdis []complemento.CFDI) []complemento.CFDI {
	if !f.IsActive() {
		return cfdis
	}

	result := make([]complemento.CFDI, 0)
	for _, c := range cfdis {
		if _, ok := f.filters[c.MetodoPago]; ok {
			result = append(result, c)
		}
	}
	return result
}

// FormaPagoFilter implementa filtrado por forma de pago
type FormaPagoFilter struct {
	filters map[string]cfdiFilterOption
}

func NewFormaPagoFilter(activeFilters map[string]cfdiFilterOption) *FormaPagoFilter {
	return &FormaPagoFilter{
		filters: activeFilters,
	}
}

func (f *FormaPagoFilter) IsActive() bool {
	for _, id := range listFilterFormaPago {
		if _, active := f.filters[id]; active {
			return true
		}
	}
	return false
}

func (f *FormaPagoFilter) Apply(cfdis []complemento.CFDI) []complemento.CFDI {
	if !f.IsActive() {
		return cfdis
	}

	result := make([]complemento.CFDI, 0)
	for _, c := range cfdis {
		if _, ok := f.filters[c.FormaPago]; ok {
			result = append(result, c)
		}
	}
	return result
}

// UsoCFDIFilter implementa filtrado por uso de CFDI
type UsoCFDIFilter struct {
	filters map[string]cfdiFilterOption
}

func NewUsoCFDIFilter(activeFilters map[string]cfdiFilterOption) *UsoCFDIFilter {
	return &UsoCFDIFilter{
		filters: activeFilters,
	}
}

func (f *UsoCFDIFilter) IsActive() bool {
	for _, id := range listFilterUsoCFDI {
		if _, active := f.filters[id]; active {
			return true
		}
	}
	return false
}

func (f *UsoCFDIFilter) Apply(cfdis []complemento.CFDI) []complemento.CFDI {
	if !f.IsActive() {
		return cfdis
	}

	result := make([]complemento.CFDI, 0)
	for _, c := range cfdis {
		if _, ok := f.filters[c.Receptor.UsoCFDI]; ok {
			result = append(result, c)
		}
	}
	return result
}

// TipoComprobanteFilter implementa filtrado por tipo de comprobante
type TipoComprobanteFilter struct {
	filters map[string]cfdiFilterOption
}

func NewTipoComprobanteFilter(activeFilters map[string]cfdiFilterOption) *TipoComprobanteFilter {
	return &TipoComprobanteFilter{
		filters: activeFilters,
	}
}

func (f *TipoComprobanteFilter) IsActive() bool {
	for _, id := range listFilterTipoComprobante {
		if _, active := f.filters[id]; active {
			return true
		}
	}
	return false
}

func (f *TipoComprobanteFilter) Apply(cfdis []complemento.CFDI) []complemento.CFDI {
	if !f.IsActive() {
		return cfdis
	}

	result := make([]complemento.CFDI, 0)
	for _, c := range cfdis {
		if _, ok := f.filters[c.TipoDeComprobante]; ok {
			result = append(result, c)
		}
	}
	return result
}

// FilterChain gestiona la cadena de filtros aplicados a los CFDIs
type FilterChain struct {
	filters []FilterStrategy
}

func NewFilterChain() *FilterChain {
	return &FilterChain{
		filters: []FilterStrategy{},
	}
}

func (c *FilterChain) AddFilter(filter FilterStrategy) {
	c.filters = append(c.filters, filter)
}

func (c *FilterChain) Apply(cfdis []complemento.CFDI) []complemento.CFDI {
	result := cfdis
	for _, filter := range c.filters {
		result = filter.Apply(result)
	}
	return result
}

// FilterFactory crea una cadena de filtros basada en el estado actual
func FilterFactory(activeFilters map[string]cfdiFilterOption) *FilterChain {
	// Comprobar si debemos ignorar los filtros
	if _, ok := activeFilters[filterIgnoreFilter.ID]; ok {
		return NewFilterChain() // Cadena vacía, no aplicará filtros
	}

	chain := NewFilterChain()
	chain.AddFilter(NewMetodoPagoFilter(activeFilters))
	chain.AddFilter(NewFormaPagoFilter(activeFilters))
	chain.AddFilter(NewUsoCFDIFilter(activeFilters))
	chain.AddFilter(NewTipoComprobanteFilter(activeFilters))

	return chain
}

// GenericFilterCFDIS es la nueva función genérica que reemplaza a filterCFDIS
func GenericFilterCFDIS(cfdis []complemento.CFDI) []complemento.CFDI {
	// Comprobar si hay filtros activos o si está activado el filtro de ignorar
	if len(activeFilters) == 0 || activeFilters[filterIgnoreFilter.ID] != (cfdiFilterOption{}) {
		return cfdis
	}

	// Crear la cadena de filtros y aplicarla
	filterChain := FilterFactory(activeFilters)
	return filterChain.Apply(cfdis)
}
