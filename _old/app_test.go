package main

import (
	"os"
	"testing"
)

const testcsv = "./csv/test.csv"

var a *App

func TestMain(m *testing.M) {
	a, _ = InitApp("./useragent.protobuf", "./templates", "./giftee.db")
	os.Exit(m.Run())
}

func TestInitApp(t *testing.T) {
	if a == nil {
		t.FailNow()
	}

	t.Logf("%+v\n", a)
}

func TestParseCSV(t *testing.T) {
	err := a.ParseCSV(true, testcsv)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}
