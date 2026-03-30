package constants

type UserStatus string

const (
	ACTIVE   = UserStatus("active")
	INACTIVE = UserStatus("inactive")
	DELETED  = UserStatus("deleted")
)
