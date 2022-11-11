package cmd

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"reflect"
)

func tableForList[T any](title string, list []T) (table *tview.Table, err error) {
	table = tview.NewTable()

	table.SetFixed(1, 1).
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical).
		SetBorderPadding(0, 0, 1, 1).
		SetBorder(true).
		SetTitle(title)

	var elem reflect.Value
	if len(list) > 0 {
		elem = reflect.ValueOf(list[0])
	} else {
		elem = reflect.ValueOf(*new(T))
	}
	nFields := elem.NumField()
	col := 0
	for i := 0; i < nFields; i++ {
		fieldName := elem.Type().Field(i).Name
		tableTag := elem.Type().Field(i).Tag.Get("table")
		if tableTag == "-" {
			continue
		} else if tableTag != "" {
			fieldName = tableTag
		}

		table.SetCell(0, col,
			tview.NewTableCell(fieldName).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
		col++
	}

	for r, e := range list {
		e := e
		elem := reflect.ValueOf(e)
		col := 0
		for i := 0; i < nFields; i++ {
			fieldValue := elem.Field(i)
			tableTag := elem.Type().Field(i).Tag.Get("table")
			if tableTag == "-" {
				continue
			}

			alignment := tview.AlignLeft
			if fieldValue.Type().Name() == "int" {
				alignment = tview.AlignRight
			}

			table.SetCell(r+1, col,
				tview.NewTableCell(fmt.Sprintf("%v", fieldValue)).
					SetTextColor(tcell.ColorWhite).
					SetAlign(alignment).
					SetSelectable(true).
					SetReference(&e))
			col++
		}
	}

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Pop()
		}
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || (event.Key() == tcell.KeyRune && event.Rune() == 'q') {
			app.Pop()
			return nil
		}
		return event
	})

	return table, nil
}
