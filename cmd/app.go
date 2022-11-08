package cmd

import (
	"github.com/rivo/tview"
)

type App struct {
	tviewApp  *tview.Application
	viewStack []tview.Primitive
}

func NewApp() *App {
	a := App{}
	a.tviewApp = tview.NewApplication().EnableMouse(true)
	a.viewStack = []tview.Primitive{}
	return &a
}

func (a *App) Push(view tview.Primitive) {
	a.viewStack = append(a.viewStack, view)
	a.tviewApp.SetRoot(view, true)
}

func (a *App) Pop() {
	if len(a.viewStack) > 1 {
		a.viewStack = a.viewStack[:len(a.viewStack)-1]
		a.tviewApp.SetRoot(a.viewStack[len(a.viewStack)-1], true)
	} else {
		a.Stop()
	}
}

func (a *App) Run() (err error) {
	return a.tviewApp.Run()
}

func (a *App) Stop() {
	a.tviewApp.Stop()
}
