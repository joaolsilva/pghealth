package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

type DatabaseListView struct {
	pg   *postgres.PG
	view tview.Primitive
}

func (dbListView *DatabaseListView) GetPrimitive() tview.Primitive {
	return dbListView.view
}

func (dbListView *DatabaseListView) Close() {
	if dbListView.pg != nil {
		dbListView.pg.Close()
	}
}

func NewDatabaseListView() (databaseListView *DatabaseListView) {
	databaseListView = &DatabaseListView{}
	var err error
	databaseListView.pg, err = postgres.NewPG()
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	databases, err := databaseListView.pg.ListDatabases()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}

	table, err := tableForList(" pghealth ", databases)

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Pop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		dbRef := table.GetCell(row, column).GetReference()
		if dbRef != nil {
			if database, ok := dbRef.(*postgres.Database); ok {
				app.Push(NewDatabaseDetailView(database))
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
		//log.Printf("Mouse capture %v %v", action, event)
		return action, event
	})

	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-C to exit")

	databaseListView.view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(table, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return databaseListView
}
