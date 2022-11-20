package main

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend"
	"log"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	app := frontend.NewApp()
	defer app.Stop()
	app.Start()
}
