package main

import (
	"bitbucket.org/ziggy192/ng_lu/src/frontend"
	"bitbucket.org/ziggy192/ng_lu/src/logger"
	"context"
	"errors"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	ctx := context.Background()
	app := frontend.NewApp()
	defer app.Stop()
	if err := app.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Err(ctx, err)
		panic(err)
	}
}
