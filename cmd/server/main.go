package main

import (
	"jargon/pkg/api"
	"jargon/pkg/db"
	"jargon/pkg/server"
)

func main() {

	// Initialise database
	database := db.OpenDatabase("/tmp/jargon.db")
	defer database.CloseDatabase()

	database.InitialiseTables()

	j := api.MakeJargon(
		"flavour-tagging",
		"Flavour Tagging",
		"The practice of obtaining the type of particle from which a jet originated",
		[]api.Tag{api.Experimental},
	)
	database.AddJargon(j)

	j = api.MakeJargon(
		"blinding",
		"Data Blinding",
		"The practice of not including actual physics data when constructing an analysis. This is to avoid biasing the analyst.",
		[]api.Tag{api.Stats, api.Experimental},
	)
	database.AddJargon(j)

	server.StartServer(&database, 8080)
}
