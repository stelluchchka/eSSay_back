package main

import (
	"essay/src/internal/app"
)

func main() {
	appInstance := app.NewApp()

	appInstance.Start()
}
