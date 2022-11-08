package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
)

func dbDetailView(db *postgres.Database) (view tview.Primitive) {
	dbDetail := tview.NewTextView().
		SetTitle(db.Name).
		SetBorder(true)

	dbDetail.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
		AddItem(dbDetail, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return view
}
