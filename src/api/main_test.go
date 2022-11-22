package api

import (
	"bitbucket.org/ziggy192/ng_lu/src/api/config"
	"bitbucket.org/ziggy192/ng_lu/src/api/store"
	"context"
	"os"
	"path"
	"runtime"
	"testing"
)

var dbStores *store.DBStores

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	GoToProjectDir("../..")

	dbStores, err = store.NewDBStores(ctx, config.New().MySQL)
	if err != nil {
		panic(err)
	}

	if err := dbStores.Reset(ctx); err != nil {
		panic(err)
	}

	code := m.Run()

	_ = dbStores.Close()
	os.Exit(code)
}

// GoToProjectDir move to root folder
func GoToProjectDir(baseDir string) {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), baseDir)
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}
