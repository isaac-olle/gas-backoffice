package auth

type SessionPrimitives struct {
	Id     string
	UserId string
}

func (s *Session) ToPrimitives() SessionPrimitives {
	return SessionPrimitives{
		Id:     s.id.String(),
		UserId: s.userId.String(),
	}
}
