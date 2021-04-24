package model

import (
	"testing"
)

func TestParseSessionID(t *testing.T) {

	testDat := []map[string]string{
		{
			"name": "/opt/fit2cloud/jumpserver/data/media/replay/3a52b5bc-f155-4a5b-9143-1ae887fe5d5e.replay.gz",
			"id":   "3a52b5bc-f155-4a5b-9143-1ae887fe5d5e",
		},
		{
			"name": "/tmp/xrdp_records/44043bc0-a4a7-11eb-b4a0-fa163e9f4f83.guac",
			"id":   "44043bc0-a4a7-11eb-b4a0-fa163e9f4f83",
		},
	}

	for i := range testDat {
		sid, err := ParseSessionID(testDat[i]["name"])
		if err != nil {
			t.Fatal(err)
		}
		if sid != testDat[i]["id"] {
			t.Fatalf("sid 不相符 %s %s", sid, testDat[i]["id"])
		}

	}
	t.Log("sid success success!")
}
