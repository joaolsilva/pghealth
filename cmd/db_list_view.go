package cmd

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

func dbListView() (view tview.Primitive) {
	table := tview.NewTable().
		SetFixed(1, 1).
		SetSelectable(true, false).
		SetSeparator(tview.Borders.Vertical)

	table.SetBorderPadding(0, 0, 1, 1)

	// Table heading
	tableColumns := []string{"Database", "Commit Ratio", "Cache Ratio", "Blocks Read", "Size"}
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
	var alignment int
	for r, d := range databases {
		d := d
		for c := 0; c < len(tableColumns); c++ {
			alignment = tview.AlignRight
			if c == 0 {
				cellText = d.Name
				alignment = tview.AlignLeft
			} else if c == 1 {
				cellText = d.CommitRatio
			} else if c == 2 {
				cellText = d.CacheHitRation
			} else if c == 3 {
				cellText = fmt.Sprintf("%v", d.BlocksRead)
			} else if c == 4 {
				cellText = d.FormattedSize
			}
			table.SetCell(r+1, c,
				tview.NewTableCell(cellText).
					SetTextColor(tcell.ColorWhite).
					SetAlign(alignment).
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
		//log.Printf("Mouse capture %v %v", action, event)
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
