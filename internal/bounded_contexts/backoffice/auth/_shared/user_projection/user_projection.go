package user_projection

type UserProjection struct {
	UserId          string
	Email           string
	Role            string
	AdminId         string
	SubscriptionIds []string
}
