package interfaces

type RegisterInterface interface {
	// Register registers a new user.
	// Return user_id, status_code, status_msg, token
	Register(username string, password string) (int64, int32, string, string)
}
