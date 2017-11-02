package csv

import "testing"

const testcsv = "./test.csv"

func TestGetUserList(t *testing.T) {
	var (
		expectedLength = 57
	)
	ul, err := GetUserList(testcsv)
	if err != nil {
		t.Fatal(err)
	}

	gotLenth := len(ul)
	if gotLenth != expectedLength {
		t.Logf("! got %v, expected %v\n", gotLenth, expectedLength)
		t.FailNow()
	}
}
