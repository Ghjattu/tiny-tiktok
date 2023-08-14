package interfaces

type FollowInterface interface {
	// CreateNewFollowRel creates a new follow relationship.
	// Return status code, status message.
	CreateNewFollowRel(followerID, followingID int64) (int32, string)

	// DeleteFollowRel delete a follow relationship by follower id and following id.
	// Return status code, status message.
	DeleteFollowRel(followerID, followingID int64) (int32, string)
}
