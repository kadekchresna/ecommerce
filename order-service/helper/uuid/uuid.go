package helper_uuid

import "github.com/google/uuid"

type UUIDHelper struct {
	u uuid.UUID
}

func NewUUID(u uuid.UUID) UUIDHelper {
	return UUIDHelper{
		u: u,
	}
}

func (u *UUIDHelper) New() uuid.UUID {
	if u.u != uuid.Nil {
		return u.u
	}

	return uuid.New()
}
