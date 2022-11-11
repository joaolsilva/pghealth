package cmd

import "log"

var app *App

func Execute() (err error) {
	app = NewApp()
	app.Push(NewDatabaseListView())

	if err := app.Run(); err != nil {
		log.Printf("%v", err)
		return err
	}

	return nil
}
