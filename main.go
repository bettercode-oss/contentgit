package main

import (
	"contentgit/app"
	"contentgit/app/datasource"
	"contentgit/config"
	"contentgit/ports/in/web"
	"log"
)

func main() {
	if err := config.InitConfig("./config"); err != nil {
		panic(err)
	}

	if err := app.NewApp(web.Router{}, datasource.ProductionDbConnector{}, app.NewComponentRegistry()).Run(); err != nil {
		log.Fatal(err)
	}
}
