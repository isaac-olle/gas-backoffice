package ids

import "github.com/google/uuid"

type Uuid struct {
	uuid.UUID
}

func GenerateUuid() Uuid {
	return Uuid{uuid.New()}
}

func NewUuid(id uuid.UUID) Uuid {
	return Uuid{id}
}

func (id Uuid) Equals(other Uuid) bool {
	return id.String() == other.String()
}

func (id Uuid) IsEmpty() bool {
	return id.UUID == uuid.Nil
}

type Id interface {
	String() string
}
