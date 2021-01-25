package upload

import (
	"encoding/base64"
	"net/url"
	"testing"
	"time"

	"jms-upload/cmd/common"
)

func TestParseSessionID(t *testing.T) {
	replayFile := "/Users/eric/Documents/fit2cloud/jumpserver/data/media/replay/3a52b5bc-f155-4a5b-9143-1ae887fe5d5e.replay.gz"
	sid, err := common.ParseSessionID(replayFile)
	if err != nil {
		t.Fatal(err)
	}
	if sid != "3a52b5bc-f155-4a5b-9143-1ae887fe5d5e" {
		t.Fatalf("ParseSessionID should be 3a52b5bc-f155-4a5b-9143-1ae887fe5d5e, but %s", sid)
	}
	t.Log("sid success: ", sid)
	t.Log("date: ", time.Now().Format("2006-01-02"))
}

func TestExecute(t *testing.T) {
	key := "853db481-b604-4839-9169-5a3f2588c416:f8ea716c-b088-44d5-848f-d5fecbdc381c"
	result := base64.StdEncoding.EncodeToString([]byte(key))
	t.Log(result)
	query := url.Values{}
	t.Logf("query: `%s`",query.Encode())
}
