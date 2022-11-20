package main

import (
	"bitbucket.org/ziggy192/ng_lu/src/api"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	app := api.NewApp()
	defer app.Stop()
	app.Start()
}
