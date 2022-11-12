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
	activityTable, err := tableForList(string(db.Name)+": Activity", activity)

	bloat, err := dbDetailView.dbConnection.ListBloat()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	bloatTable, err := tableForList(string(db.Name)+": Bloat", bloat)

	vacuumStats, err := dbDetailView.dbConnection.GetVacuumStats()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	vacuumStatsTable, err := tableForList(string(db.Name)+": Vacuum Stats", vacuumStats)

	tableCacheHitRatios, err := dbDetailView.dbConnection.GetTableCacheHitRatios()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	cacheHitRatiosTable, err := tableForList(string(db.Name)+": Cache Hit Ratio", tableCacheHitRatios)

	missingIndexes, err := dbDetailView.dbConnection.GetMissingIndexes()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	missingIndexesTable, err := tableForList(string(db.Name)+": Missing Indexes?", missingIndexes)

	uselessIndexes, err := dbDetailView.dbConnection.GetUselessIndexes()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	uselessIndexesTable, err := tableForList(string(db.Name)+": Useless Indexes?", uselessIndexes)

	tableSizes, err := dbDetailView.dbConnection.GetTableSizes()
	if err != nil {
		log.Printf("pghealth: %v", err)
		panic(err)
	}
	tableSizesTable, err := tableForList(string(db.Name)+": Table Size", tableSizes)

	dbDetailView.currentPage = 0
	pages := tview.NewPages()
	pages.AddPage("page-0", activityTable, true, true)
	pages.AddPage("page-1", cacheHitRatiosTable, true, false)
	pages.AddPage("page-2", missingIndexesTable, true, false)
	pages.AddPage("page-3", uselessIndexesTable, true, false)
	pages.AddPage("page-4", tableSizesTable, true, false)
	pages.AddPage("page-5", bloatTable, true, false)
	pages.AddPage("page-6", vacuumStatsTable, true, false)
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
		SetText(" Esc - back\tCtrl-N - next panel\tCtrl-P - previous panel")

	dbDetailView.view = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(helpInfo, 1, 1, false)

	return dbDetailView
}
