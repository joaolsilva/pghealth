package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

func dbDetailView(db *postgres.Database) (view tview.Primitive) {
	dbConnection, err := postgres.NewDBConnection(db.Name)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	bloat, err := dbConnection.ListBloat()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}

	bloatTable, err := tableForList("Bloat", bloat)

	bloatTable.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Pop()
		}
	})

	bloatTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape || (event.Key() == tcell.KeyRune && event.Rune() == 'q') {
			app.Pop()
			return nil
		}
		return event
	})

	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-C to exit")

	view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(bloatTable, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return view
}
