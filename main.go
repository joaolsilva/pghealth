package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

func main() {
	app := tview.NewApplication()

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
		panic(err)
	}
	databases, err := pg.ListDatabases()
	if err != nil {
		log.Printf("pghealth: %v", err)
		return
	}
	var cellText string
	for r, d := range databases {
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
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		dbRef := table.GetCell(row, column).GetReference()
		if dbRef != nil {
			if database, ok := dbRef.(*postgres.Database); ok {
				// database.Name selected
				_ = database
			}
		}
	}).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-C to exit")

	dbListView := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	pages := tview.NewPages()
	pages.AddAndSwitchToPage("dbList", dbListView, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
