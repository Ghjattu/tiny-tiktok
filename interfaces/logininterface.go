package interfaces

type LoginInterface interface {
	// Login logs in a user.
	// Return user_id, status_code, status_msg, token
	Login(username string, password string) (int64, int32, string, string)
}
