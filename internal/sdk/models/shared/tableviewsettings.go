// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type TableViewSettings struct {
	ColumnFilterSettings   *ColumnFilterSettings `json:"columnFilterSettings,omitempty"`
	ColumnFormatSettings   *ColumnFormatSettings `json:"columnFormatSettings,omitempty"`
	ColumnOrderSettings    *ColumnOrderSettings  `json:"columnOrderSettings,omitempty"`
	ColumnSortSettings     *ColumnSortSettings   `json:"columnSortSettings,omitempty"`
	EventDetailsPanel      *bool                 `json:"eventDetailsPanel,omitempty"`
	EventTableFields       []string              `json:"eventTableFields,omitempty"`
	RowNumberColumnWidth   *float64              `json:"rowNumberColumnWidth,omitempty"`
	ShowColumnTotals       bool                  `json:"showColumnTotals"`
	ShowColumnTotalsPinned bool                  `json:"showColumnTotalsPinned"`
	ShowRowNumbers         bool                  `json:"showRowNumbers"`
	ShowRowTotals          bool                  `json:"showRowTotals"`
	ShowRowTotalsPinned    bool                  `json:"showRowTotalsPinned"`
	WrapCells              bool                  `json:"wrapCells"`
}

func (o *TableViewSettings) GetColumnFilterSettings() *ColumnFilterSettings {
	if o == nil {
		return nil
	}
	return o.ColumnFilterSettings
}

func (o *TableViewSettings) GetColumnFormatSettings() *ColumnFormatSettings {
	if o == nil {
		return nil
	}
	return o.ColumnFormatSettings
}

func (o *TableViewSettings) GetColumnOrderSettings() *ColumnOrderSettings {
	if o == nil {
		return nil
	}
	return o.ColumnOrderSettings
}

func (o *TableViewSettings) GetColumnSortSettings() *ColumnSortSettings {
	if o == nil {
		return nil
	}
	return o.ColumnSortSettings
}

func (o *TableViewSettings) GetEventDetailsPanel() *bool {
	if o == nil {
		return nil
	}
	return o.EventDetailsPanel
}

func (o *TableViewSettings) GetEventTableFields() []string {
	if o == nil {
		return nil
	}
	return o.EventTableFields
}

func (o *TableViewSettings) GetRowNumberColumnWidth() *float64 {
	if o == nil {
		return nil
	}
	return o.RowNumberColumnWidth
}

func (o *TableViewSettings) GetShowColumnTotals() bool {
	if o == nil {
		return false
	}
	return o.ShowColumnTotals
}

func (o *TableViewSettings) GetShowColumnTotalsPinned() bool {
	if o == nil {
		return false
	}
	return o.ShowColumnTotalsPinned
}

func (o *TableViewSettings) GetShowRowNumbers() bool {
	if o == nil {
		return false
	}
	return o.ShowRowNumbers
}

func (o *TableViewSettings) GetShowRowTotals() bool {
	if o == nil {
		return false
	}
	return o.ShowRowTotals
}

func (o *TableViewSettings) GetShowRowTotalsPinned() bool {
	if o == nil {
		return false
	}
	return o.ShowRowTotalsPinned
}

func (o *TableViewSettings) GetWrapCells() bool {
	if o == nil {
		return false
	}
	return o.WrapCells
}
