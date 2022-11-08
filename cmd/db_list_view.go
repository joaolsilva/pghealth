package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

func dbListView() (view tview.Primitive) {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false)

	// Table heading
	tableColumns := []string{"Database", "Size"}
	for i, col := range tableColumns {
		table.SetCell(0, i,
			tview.NewTableCell(col).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

	pg, err := postgres.NewPG()
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	databases, err := pg.ListDatabases()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	var cellText string
	for r, d := range databases {
		d := d
		for c := 0; c < len(tableColumns); c++ {
			if c == 0 {
				cellText = d.Name
			} else {
				cellText = d.FormattedSize
			}
			table.SetCell(r+1, c,
				tview.NewTableCell(cellText).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignCenter).
					SetSelectable(true).
					SetReference(&d))
		}
	}

	table.SetBorder(true).SetTitle("pghealth")

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Pop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		dbRef := table.GetCell(row, column).GetReference()
		if dbRef != nil {
			if database, ok := dbRef.(*postgres.Database); ok {
				app.Push(dbDetailView(database))
			}
		}
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		log.Printf("Mouse capture %v %v", action, event)
		return action, event
	})

	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-C to exit")

	view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return view
}
