package util

import uuid "github.com/satori/go.uuid"

func IsUUID(sid string) bool {
	_, err := uuid.FromString(sid)
	return err == nil
}
