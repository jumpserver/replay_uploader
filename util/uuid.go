package util

import "github.com/google/uuid"

func IsUUID(sid string) bool {
	_, err := uuid.Parse(sid)
	return err == nil
}
