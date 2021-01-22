package upload

import (
	"jmsupload/cmd/common"
	"testing"
	"time"
)

func TestParseSessionID(t *testing.T) {
	replayFile := "/opt/jumpserver/koko/data/replay/2021-01-20/90b11402-39df-4ba9-b14a-2344b8585888.replay.gz"
	sid, err := common.ParseSessionID(replayFile)
	if err != nil {
		t.Fatal(err)
	}
	if sid != "90b11402-39df-4ba9-b14a-2344b8585888" {
		t.Fatalf("ParseSessionID should be 90b11402-39df-4ba9-b14a-2344b8585888, but %s", sid)
	}
	t.Log("sid success: ", sid)
	t.Log("date: ", time.Now().Format("2006-01-02"))
}
