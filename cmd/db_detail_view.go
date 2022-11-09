package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

func dbDetailView(db *postgres.Database) (view tview.Primitive) {

	bloatTable := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical)

	bloatTable.SetBorderPadding(0, 0, 1, 1)

	// Table heading
	tableColumns := []string{"Type", "Schema Name", "Object Name", "Bloat", "Waste"}
	for i, col := range tableColumns {
		bloatTable.SetCell(0, i,
			tview.NewTableCell(col).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignCenter).
				SetSelectable(false))
	}

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
	var cellText string
	var alignment int
	for r, b := range bloat {
		b := b
		for c := 0; c < len(tableColumns); c++ {
			alignment = tview.AlignRight
			if c == 0 {
				cellText = string(b.Type)
				alignment = tview.AlignLeft
			} else if c == 1 {
				cellText = b.SchemaName
				alignment = tview.AlignLeft
			} else if c == 2 {
				cellText = b.ObjectName
				alignment = tview.AlignLeft
			} else if c == 3 {
				cellText = b.Bloat
			} else if c == 4 {
				cellText = b.Waste
			}
			bloatTable.SetCell(r+1, c,
				tview.NewTableCell(cellText).
					SetTextColor(tcell.ColorWhite).
					SetAlign(alignment).
					SetSelectable(true).
					SetReference(&b))
		}
	}

	bloatTable.SetBorder(true).SetTitle("Bloat")

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
