package cmd

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

var currentPage int

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

	vacuumStats, err := dbConnection.GetVacuumStats()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	vacuumStatsTable, err := tableForList("Vacuum Stats", vacuumStats)

	tableCacheHitRatio, err := dbConnection.GetTableCacheHitRatio()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	cacheHitRatioTable, err := tableForList("Cache Hit Ratio", tableCacheHitRatio)

	currentPage = 0
	pages := tview.NewPages()
	pages.AddPage("page-0", cacheHitRatioTable, true, true)
	pages.AddPage("page-1", bloatTable, true, false)
	pages.AddPage("page-2", vacuumStatsTable, true, false)
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			if currentPage < (pages.GetPageCount() - 1) {
				currentPage++
			} else {
				currentPage = 0
			}

			pages.SwitchToPage(fmt.Sprintf("page-%v", currentPage))
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			if currentPage > 0 {
				currentPage--
			} else {
				currentPage = pages.GetPageCount() - 1
			}

			pages.SwitchToPage(fmt.Sprintf("page-%v", currentPage))
			return nil
		}
		return event
	})
	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-C to exit")

	view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return view
}
