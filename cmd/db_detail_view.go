package cmd

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/joaolsilva/pghealth/pkg/postgres"
	"github.com/rivo/tview"
	"log"
)

type DatabaseDetailView struct {
	currentPage  int
	dbConnection *postgres.DBConnection
	view         tview.Primitive
}

func (dbDetailView *DatabaseDetailView) GetPrimitive() tview.Primitive {
	return dbDetailView.view
}

func (dbDetailView *DatabaseDetailView) Close() {
	if dbDetailView.dbConnection != nil {
		dbDetailView.dbConnection.Close()
	}
}

func NewDatabaseDetailView(db *postgres.Database) (dbDetailView *DatabaseDetailView) {
	dbDetailView = &DatabaseDetailView{}
	var err error
	dbDetailView.dbConnection, err = postgres.NewDBConnection(db.Name)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}

	activity, err := dbDetailView.dbConnection.GetActivity()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	activityTable, err := tableForList("Activity", activity)

	bloat, err := dbDetailView.dbConnection.ListBloat()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	bloatTable, err := tableForList("Bloat", bloat)

	vacuumStats, err := dbDetailView.dbConnection.GetVacuumStats()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	vacuumStatsTable, err := tableForList("Vacuum Stats", vacuumStats)

	tableCacheHitRatio, err := dbDetailView.dbConnection.GetTableCacheHitRatio()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	cacheHitRatioTable, err := tableForList("Cache Hit Ratio", tableCacheHitRatio)

	missingIndexes, err := dbDetailView.dbConnection.GetMissingIndexes()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	missingIndexesTable, err := tableForList("Missing Indexes?", missingIndexes)

	dbDetailView.currentPage = 0
	pages := tview.NewPages()
	pages.AddPage("page-0", activityTable, true, true)
	pages.AddPage("page-1", cacheHitRatioTable, true, false)
	pages.AddPage("page-2", missingIndexesTable, true, false)
	pages.AddPage("page-3", bloatTable, true, false)
	pages.AddPage("page-4", vacuumStatsTable, true, false)
	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlN {
			if dbDetailView.currentPage < (pages.GetPageCount() - 1) {
				dbDetailView.currentPage++
			} else {
				dbDetailView.currentPage = 0
			}

			pages.SwitchToPage(fmt.Sprintf("page-%v", dbDetailView.currentPage))
			return nil
		} else if event.Key() == tcell.KeyCtrlP {
			if dbDetailView.currentPage > 0 {
				dbDetailView.currentPage--
			} else {
				dbDetailView.currentPage = pages.GetPageCount() - 1
			}

			pages.SwitchToPage(fmt.Sprintf("page-%v", dbDetailView.currentPage))
			return nil
		}
		return event
	})
	helpInfo := tview.NewTextView().
		SetText(" Press Ctrl-C to exit")

	dbDetailView.view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return dbDetailView
}
