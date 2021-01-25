package model

import (
	"testing"
)

func TestParseSessionID(t *testing.T) {
	replayFile := "/opt/fit2cloud/jumpserver/data/media/replay/3a52b5bc-f155-4a5b-9143-1ae887fe5d5e.replay.gz"
	sid, err := ParseSessionID(replayFile)
	if err != nil {
		t.Fatal(err)
	}
	if sid != "3a52b5bc-f155-4a5b-9143-1ae887fe5d5e" {
		t.Fatalf("ParseSessionID should be 3a52b5bc-f155-4a5b-9143-1ae887fe5d5e, but %s", sid)
	}
	t.Log("sid success: ", sid)
}
