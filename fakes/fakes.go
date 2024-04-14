package fakes

const (
	Admin             = "Admin"
	AuthorizedUser    = "AuthorizedUser"
	NotAuthorizedUser = "NotAuthorizedUser"
)

func GetRole(token string) string {
	switch token {
	case Admin:
		return Admin
	case AuthorizedUser:
		return AuthorizedUser
	default:
		return NotAuthorizedUser
	}
}
