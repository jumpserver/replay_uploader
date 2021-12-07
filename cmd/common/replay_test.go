package common

import "testing"

func TestParseDateFromPath(t *testing.T) {
	path := "/opt/jumpserver/koko/data/replays/2021-11-26/sidreplay_filename"
	targetDate, ok := ParseDateFromPath(path)
	if !ok {
		t.Fatalf("parse date from path failed: %s", path)

	}
	t.Log("target Date: ", targetDate)

}
