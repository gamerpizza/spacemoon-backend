package model

import (
	"github.com/google/uuid"
)

type UserID string

func ParseUUID(issClaim string) UserID {
	return UserID(uuid.MustParse(issClaim).String())
}

func UUIDToUserID(uuid uuid.UUID) UserID {
	return UserID(uuid.String())
}
