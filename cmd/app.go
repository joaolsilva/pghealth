package cmd

import (
	"github.com/rivo/tview"
)

type View interface {
	GetPrimitive() tview.Primitive
	Close()
}

type App struct {
	tviewApp  *tview.Application
	viewStack []View
}

func NewApp() *App {
	a := App{}
	a.tviewApp = tview.NewApplication().EnableMouse(true)
	a.viewStack = []View{}
	return &a
}

func (a *App) Push(view View) {
	a.viewStack = append(a.viewStack, view)
	a.tviewApp.SetRoot(view.GetPrimitive(), true)
}

func (a *App) Pop() {
	if len(a.viewStack) > 0 {
		a.viewStack[len(a.viewStack)-1].Close()
	}

	if len(a.viewStack) > 1 {
		a.viewStack = a.viewStack[:len(a.viewStack)-1]
		a.tviewApp.SetRoot(a.viewStack[len(a.viewStack)-1].GetPrimitive(), true)
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
