package csv

import "testing"

func TestGetPairs(t *testing.T) {
	var (
		expect = make(map[string]string)
		got    map[string]string
	)

	expect["santa"] = "giftee"

	got, err := GetPairs("./SantaList.csv")
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range got {
		if expect[k] != v {
			t.Logf("! got: %v => %v, expected: %v\n", k, v, expect[k])
			t.FailNow()
		}
	}
}
